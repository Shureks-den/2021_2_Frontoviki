package delivery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
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

	s.HandleFunc("/{id:[0-9]+}", middleware.SetSCRFToken(sm.SoftCheckAuthorized(ah.AdvertDetailHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.AdvertUpdateHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.DeleteAdvertHandler))).Methods(http.MethodDelete, http.MethodOptions)
	s.Handle("/{id:[0-9]+}/close", sm.CheckAuthorized(http.HandlerFunc(ah.CloseAdvertHandler))).Methods(http.MethodPost, http.MethodOptions)

	s.Handle("/{id:[0-9]+}/images", sm.CheckAuthorized(http.HandlerFunc(ah.UploadImageHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/{id:[0-9]+}/images", sm.CheckAuthorized(http.HandlerFunc(ah.RemoveImageHandler))).Methods(http.MethodDelete, http.MethodOptions)

	s.HandleFunc("/salesman/{id:[0-9]+}", middleware.SetSCRFToken(sm.SoftCheckAuthorized(ah.SalesmanPageHandler))).Methods(http.MethodGet, http.MethodOptions)

	s.Handle("/favorite", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ah.FavoriteListHandler)))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/favorite/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.AddFavoriteHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/favorite/{id:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ah.RemoveFavoriteHandler))).Methods(http.MethodDelete, http.MethodOptions)

	s.Handle("/price_history", sm.CheckAuthorized(http.HandlerFunc(ah.UpdatePriceHistory))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/price_history/{id:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ah.GetPriceHistory)))).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/promotion", ah.HandlePromotion).Methods(http.MethodPost, http.MethodOptions)

	s.Handle("/recomendations/{id:[0-9]+}", middleware.SetSCRFToken(sm.SoftCheckAuthorized(ah.RecomendationsHandler))).Methods(http.MethodGet, http.MethodOptions)
}

func (ah *AdvertHandler) HandlePromotion(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	err := r.ParseForm()
	fmt.Println("request.Form::")
	for key, value := range r.Form {
		fmt.Printf("Key:%s, Value:%s\n", key, value)
	}
	fmt.Println("\nrequest.PostForm::")
	for key, value := range r.PostForm {
		fmt.Printf("Key:%s, Value:%s\n", key, value)
	}
	if err != nil {
		logger.Warnf("invalid form: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("invalid write: %s", err.Error())
		}
		return
	}

	label := r.PostFormValue("label")
	fmt.Println(label)
	vals := strings.Split(label, "__")
	userId, _ := strconv.ParseInt(vals[0], 10, 64)
	adId, _ := strconv.ParseInt(vals[1], 10, 64)
	lvl, _ := strconv.ParseInt(vals[2], 10, 64)
	promo := &models.Promotion{
		AdvertId:   adId,
		PromoLevel: lvl,
	}

	err = ah.advtUsecase.UpdatePromotion(userId, promo)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "promotion updated successfully", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	query := u.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	advts, err := ah.advtUsecase.GetListAdvt(page.PageNum, page.Count, true)
	if err != nil {
		logger.Warnf("get list advt error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyAdverts{Advert: advts}
	_, err = w.Write(models.ToBytes(http.StatusOK, "adverts found successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, &advert)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
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

		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.CreateAdvert(userId, &advert)
	if err != nil {
		logger.Infof("can not create adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertShort{AdvertShort: *advert.ToShort()}
	_, err = w.Write(models.ToBytes(http.StatusCreated, "advert created successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
	var userId int64 = 0
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}
	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse id adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	advert, err := ah.advtUsecase.GetAdvert(advertId, userId, true)
	if err != nil {
		logger.Warnf("can not get adv by advId: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	advert.Views, err = ah.advtUsecase.GetAdvertViews(advertId)
	if err != nil {
		logger.Warnf("can not get views: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	salesman, err := ah.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		logger.Warnf("can not parse id adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	rateStat, err := ah.userUsecase.GetRating(userId, salesman.Id)
	if err != nil {
		logger.Debugf("can not get user's statistic with id %d: %s", salesman.Id, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	history, err := ah.advtUsecase.GetPriceHistory(advertId)
	if err != nil {
		logger.Debugf("can not get price history: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	favCount, err := ah.advtUsecase.GetFavoriteCount(advertId)
	if err != nil {
		logger.Debugf("can not get favorite count: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertDetail{Advert: *advert, Salesman: *salesman, Rating: *rateStat,
		PriceHistory: history, FavoriteCount: favCount}
	_, err = w.Write(models.ToBytes(http.StatusOK, "advert found successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	var newAdvert models.Advert
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, &newAdvert)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	if newAdvert.PublisherId != userId {
		logger.Info("no rights to access")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(models.ToBytes(http.StatusConflict, "no rights to access", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
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

		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.UpdateAdvert(advertId, &newAdvert)
	if err != nil {
		logger.Warnf("bad update adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvert{Advert: newAdvert}
	_, err = w.Write(models.ToBytes(http.StatusCreated, "advert updated successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.DeleteAdvert(advertId, userId)
	if err != nil {
		logger.Warnf("can not delete adv with id %d: %s", advertId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "advert deleted successfully", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.CloseAdvert(advertId, userId)
	if err != nil {
		logger.Warnf("can not close adv with id %d: %s", advertId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "advert closed successfully", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
// @Router /adverts/{id}/image [post]
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
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	defer r.Body.Close()
	err = r.ParseMultipartForm(8 << 20) // 8Мб
	if err != nil {
		logger.Warnf("can not parse adv id: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.InternalError)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	if len(r.MultipartForm.File["images"]) == 0 {
		logger.Warnf("empty image form: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	files := r.MultipartForm.File["images"]
	advert, err := ah.advtUsecase.UploadImages(files, advertId, userId)
	if err != nil {
		logger.Warnf("user %d can not upload image in adv %d: %s", userId, advertId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	salesman, err := ah.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		logger.Warnf("can not get user by id %d: %s", advert.PublisherId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdvertDetail{Advert: *advert, Salesman: *salesman}
	_, err = w.Write(models.ToBytes(http.StatusOK, "images uploaded successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Debug("image uploaded successfully")
}

// UploadImageHandler godoc
// @Summary Upload images for advert
// @Description Upload images for advert
// @Tags advert
// @Accept multipart/form-data
// @Produce application/json
// @Param id path integer true "Advert id"
// @Param images body models.AdvertImages true "Pathes to images"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /adverts/{id}/image [delete]
func (ah *AdvertHandler) RemoveImageHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64 = 0
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	images := &models.AdvertImages{}
	err = json.Unmarshal(buf, images)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.RemoveImages(images.ImagesPath, advertId, userId)
	if err != nil {
		logger.Warnf("can not delete images of %d advert: %s", advertId, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "images removed successfully", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
	var userId int64 = 0
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	salesmanId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse string: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		logger.Warnf("can not parse path: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	query := u.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	salesman, err := ah.userUsecase.GetById(salesmanId)
	if err != nil {
		logger.Warnf("can not parse path: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	adverts, err := ah.advtUsecase.GetAdvertListByPublicherId(salesmanId, true, page)
	if err != nil {
		logger.Warnf("can not get adverts: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	shortAdverts := ah.advtUsecase.AdvertsToShort(adverts)

	rateStat, err := ah.userUsecase.GetRating(userId, salesman.Id)
	if err != nil {
		logger.Debugf("can not get user's statistic with id %d: %s", salesman.Id, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodySalesmanPage{Salesman: *salesman, Adverts: shortAdverts, Rating: *rateStat}
	_, err = w.Write(models.ToBytes(http.StatusOK, "salesman profile provided", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	query := u.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	adverts, err := ah.advtUsecase.GetAdvertListByPublicherId(userId, false, page)
	if err != nil {
		logger.Warnf("unable to got adverts: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdverts{Advert: adverts}
	_, err = w.Write(models.ToBytes(http.StatusOK, "archive got", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
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
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	categoryName = bluemonday.UGCPolicy().Sanitize(categoryName)

	adverts, err := ah.advtUsecase.GetAdvertListByCategory(categoryName, page)
	if err != nil {
		logger.Warnf("can not get adverts: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdverts{Advert: adverts}
	_, err = w.Write(models.ToBytes(http.StatusOK, "adverts got successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// FavoriteListHandler godoc
// @Summary Get list of favorite adverts
// @Description Get list of favorite adverts
// @Tags favorite
// @Accept application/json
// @Produce application/json
// @Param page query string false "Page num"
// @Param count query string false "Count adverts per page"
// @Success 200 {object} models.HttpBodyInterface{body=HttpBodyAdverts}
// @failure default {object} models.HttpError
// @Router /adverts/favorite [get]
func (ah *AdvertHandler) FavoriteListHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	query := r.URL.Query()
	page, err := models.NewPage(query.Get("page"), query.Get("count"))
	if err != nil {
		logger.Warnf("can not create page: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	adverts, err := ah.advtUsecase.GetFavoriteList(userId, page)
	if err != nil {
		logger.Warnf("can not get favorite list: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdverts{Advert: adverts}
	_, err = w.Write(models.ToBytes(http.StatusOK, "favorite adverts got successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// AddFavoriteHandler godoc
// @Summary Add to favorites
// @Description Add to favorites
// @Tags favorite
// @Accept application/json
// @Produce application/json
// @Param id query int true "Advert id"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /adverts/favorite/{id} [post]
func (ah *AdvertHandler) AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse string: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.AddFavorite(userId, advertId)
	if err != nil {
		logger.Warnf("can not add to favorite: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "added to favorite", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// RemoveFavoriteHandler godoc
// @Summary Remove to favorites
// @Description Remove to favorites
// @Tags favorite
// @Accept application/json
// @Produce application/json
// @Param id query int true "Advert id"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /adverts/favorite/{id} [delete]
func (ah *AdvertHandler) RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse string: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.RemoveFavorite(userId, advertId)
	if err != nil {
		logger.Warnf("can not remove from favorite: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "removed from favorite", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

func (ah *AdvertHandler) UpdatePriceHistory(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	defer r.Body.Close()
	adPrice := &models.AdvertPrice{}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, adPrice)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	_, err = govalidator.ValidateStruct(adPrice)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ah.advtUsecase.UpdateAdvertPrice(userId, adPrice)
	if err != nil {
		logger.Warnf("can not decode images: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, fmt.Sprintf("price of %d advert updated", adPrice.AdvertId), nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

func (ah *AdvertHandler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse string: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	priceHistory, err := ah.advtUsecase.GetPriceHistory(advertId)
	if err != nil {
		logger.Warnf("can not get price history of advert %d: %s", advertId, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyPriceHistory{History: priceHistory}
	_, err = w.Write(models.ToBytes(http.StatusOK, "favorite adverts got successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

func (ah *AdvertHandler) RecomendationsHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64 = 0
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}
	vars := mux.Vars(r)
	advertId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logger.Warnf("can not parse id adv: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	recs, err := ah.advtUsecase.GetRecomendations(advertId, 10, userId)
	if err != nil {
		logger.Warnf("can not get recomendations: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyAdverts{Advert: recs}
	_, err = w.Write(models.ToBytes(http.StatusOK, "adverts found successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}
