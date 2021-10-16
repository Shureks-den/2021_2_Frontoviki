package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
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

func (ah *AdvertHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/adverts").Subrouter()

	s.HandleFunc("", ah.AdvertListHandler).Methods(http.MethodGet, http.MethodOptions)

	s.Handle("", sm.CheckAuthorized(http.HandlerFunc(ah.CreateAdvertHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/{id:[0-9]+}", ah.AdvertDetailHandler).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.AdvertUpdateHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.DeleteAdvertHandler))).Methods(http.MethodDelete, http.MethodOptions)

	s.Handle("/{id:[0-9]+}/close", sm.CheckAuthorized(http.HandlerFunc(ah.CloseAdvertHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/{id:[0-9]+}/upload", sm.CheckAuthorized(http.HandlerFunc(ah.UploadImageHandler))).Methods(http.MethodPost, http.MethodOptions)
}

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
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyAdvts{Advert: advts}
	w.Write(models.ToBytes(http.StatusOK, "adverts found successfully", body))
}

func (ah *AdvertHandler) CreateAdvertHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	var advert models.Advert
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(advert)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	err = ah.advtUsecase.CreateAdvert(userId, &advert)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertShort{AdvertShort: *advert.ToShort()}
	w.Write(models.ToBytes(http.StatusCreated, "advert created successfully", body))
}

func (ah *AdvertHandler) AdvertDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	advert, err := ah.advtUsecase.GetAdvert(advertId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	salesman, err := ah.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertDetail{Advert: *advert, Salesman: *salesman}
	w.Write(models.ToBytes(http.StatusOK, "advert found successfully", body))
}

func (ah *AdvertHandler) AdvertUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var newAdvert models.Advert
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&newAdvert)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	if newAdvert.PublisherId != userId {
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusConflict, "no rights to access", nil))
		return
	}

	_, err = govalidator.ValidateStruct(newAdvert)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	err = ah.advtUsecase.UpdateAdvert(advertId, &newAdvert)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvert{Advert: newAdvert}
	w.Write(models.ToBytes(http.StatusCreated, "advert updated successfully", body))
}

func (ah *AdvertHandler) DeleteAdvertHandler(w http.ResponseWriter, r *http.Request) {
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

	err = ah.advtUsecase.DeleteAdvert(advertId, userId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "advert deleted successfully", nil))
}

func (ah *AdvertHandler) CloseAdvertHandler(w http.ResponseWriter, r *http.Request) {
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

	err = ah.advtUsecase.CloseAdvert(advertId, userId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "advert closed successfully", nil))
}

func (ah *AdvertHandler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.InternalError)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	if len(r.MultipartForm.File["images"]) == 0 {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	files := r.MultipartForm.File["images"]
	advert, err := ah.advtUsecase.UploadImages(files, advertId, userId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	salesman, err := ah.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertDetail{Advert: *advert, Salesman: *salesman}
	w.Write(models.ToBytes(http.StatusOK, "images uploaded successfully", body))
}
