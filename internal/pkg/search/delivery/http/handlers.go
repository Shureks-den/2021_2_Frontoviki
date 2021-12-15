package delivery

import (
	"log"
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/search"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

type SearchHandler struct {
	searchUsecase search.SearchUsecase
}

func NewSearchHandler(searchUsecase search.SearchUsecase) *SearchHandler {
	return &SearchHandler{
		searchUsecase: searchUsecase,
	}
}

var (
	logger logging.Logger = logging.GetLogger()
)

func (sh *SearchHandler) Routing(r *mux.Router) {
	r.Handle("/search", middleware.SetSCRFToken(http.HandlerFunc(sh.SearchHandler))).Methods(http.MethodGet, http.MethodOptions)
}

// AdvertListHandler godoc
// @Summary Get list of all adverts
// @Description Get list of all adverts
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param query query string true "Query text"
// @Param time_duration query int false "Time duration"
// @Param category query string false "Category"
// @Param latitude query string false "Latitude"
// @Param longitude query string false "Longitude"
// @Param radius query string false "Radius"
// @Param sorting_name query string false "Sort by name"
// @Param sorting_date query string false "Sort by date"
// @Param page query string false "Page num"
// @Param count query string false "Count adverts per page"
// @Success 200 {object} models.HttpBodyInterface{body=[]models.Advert}
// @failure default {object} models.HttpError
// @Router /search [get]
func (sh *SearchHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	query := r.URL.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			log.Printf("error writing response %v", err.Error())
		}
		return
	}

	filter, err := models.NewSearchFilter(&query)
	if err != nil {
		logger.Warnf("can not create search filter: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			log.Printf("error writing response %v", err.Error())
		}
		return
	}

	_, err = govalidator.ValidateStruct(filter)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			log.Printf("error writing response %v", err.Error())
		}
		return
	}

	sanitize := bluemonday.UGCPolicy()
	filter.Query = sanitize.Sanitize(filter.Query)
	filter.Category = sanitize.Sanitize(filter.Category)

	adverts, err := sh.searchUsecase.SearchWithFilter(filter, page)
	if err != nil {
		logger.Warnf("can not use search: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			log.Printf("error writing response %v", err.Error())
		}
		return
	}

	body := models.HttpBodyAdverts{Advert: adverts}
	_, err = w.Write(models.ToBytes(http.StatusOK, "adverts found successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Info("adverts found successfully")
}
