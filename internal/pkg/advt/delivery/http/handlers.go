package delivery

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strconv"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

type AdvertHandler struct {
	advtUsecase advt.AdvtUsecase
	userUsecase user.UserUsecase
}

func NewAdvertHandler(advtUsecase advt.AdvtUsecase, userUsecase user.UserUsecase) *AdvertHandler {
	return &AdvertHandler{
		advtUsecase: advtUsecase,
		userUsecase: userUsecase,
	}
}

var (
	logger logging.Logger = logging.GetLogger()
)

func (ah *AdvertHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/adverts").Subrouter()

	s.HandleFunc("", middleware.SetSCRFToken(http.HandlerFunc(ah.AdvertListHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("", sm.CheckAuthorized(http.HandlerFunc(ah.CreateAdvertHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/archive", middleware.SetSCRFToken(http.Handler(sm.CheckAuthorized(http.HandlerFunc(ah.ArchiveHandler))))).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/category/{category}", middleware.SetSCRFToken(http.HandlerFunc(ah.AdvertListByCategoryHandler))).Methods(http.MethodGet, http.MethodOptions)

	s.HandleFunc("/{id:[0-9]+}", middleware.SetSCRFToken(http.HandlerFunc(ah.AdvertDetailHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.AdvertUpdateHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.DeleteAdvertHandler))).Methods(http.MethodDelete, http.MethodOptions)
	s.Handle("/{id:[0-9]+}/close", sm.CheckAuthorized(http.HandlerFunc(ah.CloseAdvertHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/{id:[0-9]+}/upload", sm.CheckAuthorized(http.HandlerFunc(ah.UploadImageHandler))).Methods(http.MethodPost, http.MethodOptions)

	s.HandleFunc("/salesman/{id:[0-9]+}", middleware.SetSCRFToken(http.HandlerFunc(ah.SalesmanPageHandler))).Methods(http.MethodGet, http.MethodOptions)
}

// AdvertListHandler godoc
// @Summary Get list of all adverts
// @Description Get list of all adverts
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param page query string false "Page num"
// @Param count query string false "Count adverts per page"
// @Success 200 {object} models.HttpBodyInterface{body=[]models.Advert}
// @failure default {object} models.HttpError
// @Router /adverts [get]
func (ah *AdvertHandler) AdvertListHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	query := u.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	advts, err := ah.advtUsecase.GetListAdvt(page.PageNum, page.Count, true)
	logger = logger.GetLoggerWithFields((r.Context().Value("logger fields")).(logrus.Fields))
	if err != nil {
		logger.Warnf("get list advt error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyAdverts{Advert: advts}
	w.Write(models.ToBytes(http.StatusOK, "adverts found successfully", body))
	logger.Info("adverts found successfully")
}

// CreateAdvertHandler godoc
// @Summary Create advert
// @Description Create advert
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param new_advert body models.Advert true "Advert"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyAdvertShort}
// @failure default {object} models.HttpError
// @Router /adverts [post]
func (ah *AdvertHandler) CreateAdvertHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	var advert models.Advert
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		logger.Warnf("invalid body: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	sanitize := bluemonday.UGCPolicy()
	advert.Name = sanitize.Sanitize(advert.Name)
	advert.Description = sanitize.Sanitize(advert.Description)
	advert.Location = sanitize.Sanitize(advert.Location)
	advert.Category = sanitize.Sanitize(advert.Category)

	_, err = govalidator.ValidateStruct(advert)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	err = ah.advtUsecase.CreateAdvert(userId, &advert)
	if err != nil {
		logger.Infof("can not create adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertShort{AdvertShort: *advert.ToShort()}
	w.Write(models.ToBytes(http.StatusCreated, "advert created successfully", body))
	logger.Debug("advert created successfully")
}

// AdvertDetailHandler godoc
// @Summary Get detail advert
// @Description Get detail advert
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Advert id"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyAdvert}
// @failure default {object} models.HttpError
// @Router /adverts/{id} [get]
func (ah *AdvertHandler) AdvertDetailHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse id adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	advert, err := ah.advtUsecase.GetAdvert(advertId)
	if err != nil {
		logger.Warnf("can not get adv by advId: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	salesman, err := ah.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		logger.Warnf("can not parse id adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertDetail{Advert: *advert, Salesman: *salesman}
	w.Write(models.ToBytes(http.StatusOK, "advert found successfully", body))
	logger.Debug("advert found successfully")
}

// AdvertUpdateHandler godoc
// @Summary Update advert
// @Description Update advert
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Advert id"
// @Param advert body models.Advert true "New advert"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyAdvert}
// @failure default {object} models.HttpError
// @Router /adverts/{id} [post]
func (ah *AdvertHandler) AdvertUpdateHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse adv id: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	var newAdvert models.Advert
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&newAdvert)
	if err != nil {
		logger.Warnf("can not decode adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	if newAdvert.PublisherId != userId {
		logger.Info("no rights to access")
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusConflict, "no rights to access", nil))
		return
	}

	sanitize := bluemonday.UGCPolicy()
	newAdvert.Name = sanitize.Sanitize(newAdvert.Name)
	newAdvert.Description = sanitize.Sanitize(newAdvert.Description)
	newAdvert.Location = sanitize.Sanitize(newAdvert.Location)
	newAdvert.Category = sanitize.Sanitize(newAdvert.Category)

	_, err = govalidator.ValidateStruct(newAdvert)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	err = ah.advtUsecase.UpdateAdvert(advertId, &newAdvert)
	if err != nil {
		logger.Warnf("bad update adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvert{Advert: newAdvert}
	w.Write(models.ToBytes(http.StatusCreated, "advert updated successfully", body))
	logger.Debug("advert updated successfully")
}

// DeleteAdvertHandler godoc
// @Summary Delete advert
// @Description Delete advert
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Advert id"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /adverts/{id} [delete]
func (ah *AdvertHandler) DeleteAdvertHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse adv id: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = ah.advtUsecase.DeleteAdvert(advertId, userId)
	if err != nil {
		logger.Warnf("can not delete adv with id %d: %s", advertId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "advert deleted successfully", nil))
	logger.Debug("advert deleted successfully")
}

// CloseAdvertHandler godoc
// @Summary Close advert
// @Description Close advert
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Advert id"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /adverts/{id}/close [post]
func (ah *AdvertHandler) CloseAdvertHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse adv id: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = ah.advtUsecase.CloseAdvert(advertId, userId)
	if err != nil {
		logger.Warnf("can not close adv with id %d: %s", advertId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "advert closed successfully", nil))
	logger.Debug("advert closed successfully")
}

// UploadImageHandler godoc
// @Summary Upload images for advert
// @Description Upload images for advert
// @Tags advert
// @Accept multipart/form-data
// @Produce application/json
// @Param id path integer true "Advert id"
// @Param images formData file true "Uploaded images"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyAdvertDetail{advert=models.Advert,salesman=models.Profile}}
// @failure default {object} models.HttpError
// @Router /adverts/{id}/upload [post]
func (ah *AdvertHandler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	defer r.Body.Close()
	err = r.ParseMultipartForm(8 << 20) // 8Мб
	if err != nil {
		logger.Warnf("can not parse adv id: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.InternalError)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	if len(r.MultipartForm.File["images"]) == 0 {
		logger.Warnf("empty image form: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	files := r.MultipartForm.File["images"]
	advert, err := ah.advtUsecase.UploadImages(files, advertId, userId)
	if err != nil {
		logger.Warnf("user %d can not upload image in adv %d: %s", userId, advertId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	salesman, err := ah.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		logger.Warnf("can not get user by id %d: %s", advert.PublisherId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertDetail{Advert: *advert, Salesman: *salesman}
	w.Write(models.ToBytes(http.StatusOK, "images uploaded successfully", body))
	logger.Debug("image uploaded successfully")
}

// SalesmanPageHandler godoc
// @Summary Get salesman page and his adverts
// @Description Get salesman page and his adverts
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Salesman id"
// @Param page query string false "Page num"
// @Param count query string false "Count adverts per page"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodySalesmanPage{salesman=models.Profile,adverts=[]models.AdvertShort}}
// @failure default {object} models.HttpError
// @Router /adverts/salesman/{id} [get]
func (ah *AdvertHandler) SalesmanPageHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	vars := mux.Vars(r)
	salesmanId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse string: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		logger.Warnf("can not parse path: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	query := u.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	salesman, err := ah.userUsecase.GetById(salesmanId)
	if err != nil {
		logger.Warnf("can not parse path: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	adverts, err := ah.advtUsecase.GetAdvertListByPublicherId(salesmanId, true, page)
	if err != nil {
		logger.Warnf("can not get adverts: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	shortAdverts := ah.advtUsecase.AdvertsToShort(adverts)

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodySalesmanPage{Salesman: *salesman, Adverts: shortAdverts}
	w.Write(models.ToBytes(http.StatusOK, "salesman profile provided", body))
}

// SalesmanPageHandler godoc
// @Summary Get salesman page and his adverts
// @Description Get salesman page and his adverts
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyAdverts}
// @failure default {object} models.HttpError
// @Router /adverts/archive [get]
func (ah *AdvertHandler) ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		logger.Warnf("can not parse path: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	query := u.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	adverts, err := ah.advtUsecase.GetAdvertListByPublicherId(userId, false, page)
	if err != nil {
		logger.Warnf("unable to got adverts: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdverts{Advert: adverts}
	w.Write(models.ToBytes(http.StatusOK, "archive got", body))
}

// AdvertListByCategoryHandler godoc
// @Summary Get adverts by category
// @Description Get adverts by category
// @Tags advert
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Salesman id"
// @Param page query string false "Page num"
// @Param count query string false "Count adverts per page"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyAdverts}
// @failure default {object} models.HttpError
// @Router /adverts/category/{category} [get]
func (ah *AdvertHandler) AdvertListByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	categoryName := path.Base(r.URL.Path)
	query := r.URL.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	categoryName = bluemonday.UGCPolicy().Sanitize(categoryName)

	adverts, err := ah.advtUsecase.GetAdvertListByCategory(categoryName, page)
	if err != nil {
		logger.Warnf("can not get adverts: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdverts{Advert: adverts}
	w.Write(models.ToBytes(http.StatusOK, "adverts got successfully", body))
}
