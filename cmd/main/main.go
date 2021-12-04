package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"yula/internal/config"

	_ "github.com/jackc/pgx/stdlib"

	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"

	"yula/internal/pkg/logging"
	userHttp "yula/internal/pkg/user/delivery/http"
	userRep "yula/internal/pkg/user/repository"
	userUse "yula/internal/pkg/user/usecase"

	"yula/internal/pkg/middleware"
	sessHttp "yula/internal/pkg/session/delivery/http"

	advtHttp "yula/internal/pkg/advt/delivery/http"
	advtRep "yula/internal/pkg/advt/repository"
	advtUse "yula/internal/pkg/advt/usecase"

	cartHttp "yula/internal/pkg/cart/delivery/http"
	cartRep "yula/internal/pkg/cart/repository"
	cartUse "yula/internal/pkg/cart/usecase"

	srchHttp "yula/internal/pkg/search/delivery/http"
	srchRep "yula/internal/pkg/search/repository"
	srchUse "yula/internal/pkg/search/usecase"

	categoryHttp "yula/internal/pkg/category/delivery/http"

	chatHttp "yula/internal/pkg/chat/delivery/http"
	metrics "yula/internal/pkg/metrics"
	metricsHttp "yula/internal/pkg/metrics/delivery"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	authProto "yula/proto/generated/auth"
	categoryProto "yula/proto/generated/category"
	chatProto "yula/proto/generated/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// _ "yula/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	govalidator.SetFieldsRequiredByDefault(true)
}

func getPostgres(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	return db
}

func CreateGRPCClient(endPoint string, opt grpc.DialOption) *grpc.ClientConn {
	grpcAuthClient, err := grpc.Dial(endPoint, opt)
	if err != nil {
		log.Fatal("cant open grpc conn")
	}

	return grpcAuthClient
}

func CreateSecureGRPCClient(endPoint string, pemServerCA []byte) *grpc.ClientConn {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		log.Fatal("can not append certs from pem")
	}

	// Create the credentials and return it
	configG := &tls.Config{
		RootCAs: certPool,
	}

	return CreateGRPCClient(endPoint, grpc.WithTransportCredentials(credentials.NewTLS(configG)))
}

// @title Volchock's API
// @version 1.0
// @description Advert placement service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	logger := logging.GetLogger()

	if err := config.LoadConfig(); err != nil {
		logger.Errorf("error with load config: %s", err.Error())
		return
	}

	sqlDB := getPostgres(config.Cfg.GetPostgresUrl())
	defer sqlDB.Close()

	r := mux.NewRouter()
	r.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)
	api := r.PathPrefix("").Subrouter()

	// ставим мидлварину с метриками
	m := metrics.NewMetrics(r)
	mmw := metricsHttp.NewMetricsMiddleware(m)
	r.Use(mmw.ScanMetrics)

	api.Use(middleware.CorsMiddleware)
	api.Use(middleware.ContentTypeMiddleware)
	api.Use(middleware.LoggerMiddleware)
	//api.Use(middleware.CSRFMiddleWare())

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ar := advtRep.NewAdvtRepository(sqlDB)
	ur := userRep.NewUserRepository(sqlDB)
	rr := userRep.NewRatingRepository(sqlDB)
	cr := cartRep.NewCartRepository(sqlDB)
	serr := srchRep.NewSearchRepository(sqlDB)

	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)
	au := advtUse.NewAdvtUsecase(ar, ilu)
	uu := userUse.NewUserUsecase(ur, rr, ilu)
	cu := cartUse.NewCartUsecase(cr)
	seru := srchUse.NewSearchUsecase(serr, ar)

	ah := advtHttp.NewAdvertHandler(au, uu)
	ch := cartHttp.NewCartHandler(cu, uu, au)
	serh := srchHttp.NewSearchHandler(seru)

	pemServerCA, err := ioutil.ReadFile(config.Cfg.GetSelfSignedCrt())
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcChatClient := CreateSecureGRPCClient(config.Cfg.GetChatEndPoint(), pemServerCA)
	defer grpcChatClient.Close()

	grpcAuthClient := CreateGRPCClient(config.Cfg.GetAuthEndPoint(), grpc.WithInsecure())
	defer grpcAuthClient.Close()

	grpcCategoryClient := CreateGRPCClient(config.Cfg.GetCategoryEndPoint(), grpc.WithInsecure())
	defer grpcCategoryClient.Close()

	uh := userHttp.NewUserHandler(uu, authProto.NewAuthClient(grpcAuthClient))
	sh := sessHttp.NewSessionHandler(authProto.NewAuthClient(grpcAuthClient), uu)
	cath := categoryHttp.NewCategoryHandler(categoryProto.NewCategoryClient(grpcCategoryClient))
	chth := chatHttp.NewChatHandler(chatProto.NewChatClient(grpcChatClient))

	sm := middleware.NewSessionMiddleware(authProto.NewAuthClient(grpcAuthClient))

	ah.Routing(api, sm)
	uh.Routing(api, sm)
	sh.Routing(api)
	ch.Routing(api, sm)
	serh.Routing(api)
	cath.Routing(api)
	middleware.Routing(api)
	chth.Routing(api, sm)

	port := config.Cfg.GetMainPort()
	fmt.Printf("start serving ::%s\n", port)

	var error error
	secure := config.Cfg.IsSecure()
	if secure {
		error = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), config.Cfg.GetHTTPSCrt(), config.Cfg.GetHTTPSKey(), r)
	} else {
		error = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	}
	logger.Errorf("http serve error %v", error)
}
