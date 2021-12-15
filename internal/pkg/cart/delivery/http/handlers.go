package delivery

import (
	"io/ioutil"
	"net/http"
	"strconv"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/cart"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

var (
	logger logging.Logger = logging.GetLogger()
)

type CartHandler struct {
	cartUsecase   cart.CartUsecase
	userUsecase   user.UserUsecase
	advertUsecase advt.AdvtUsecase
}

func NewCartHandler(cartUsecase cart.CartUsecase, userUsecase user.UserUsecase, advertUsecase advt.AdvtUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase:   cartUsecase,
		userUsecase:   userUsecase,
		advertUsecase: advertUsecase,
	}
}

func (ch *CartHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/cart").Subrouter()
	s.Use(sm.CheckAuthorized)

	s.HandleFunc("/one", ch.UpdateOneAdvertHandler).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("", ch.UpdateAllCartHandler).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("", middleware.SetSCRFToken(http.HandlerFunc(ch.GetCartHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/clear", ch.ClearCartHandler).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
}

// UpdateOneAdvertHandler godoc
// @Summary Update single advert in cart
// @Description Update single advert in cart
// @Tags cart
// @Accept application/json
// @Produce application/json
// @Param add_to_cart body models.CartHandler true "Add to Cart model"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /cart/one [post]
func (ch *CartHandler) UpdateOneAdvertHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	var cartInputed models.CartHandler
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err := w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Errorf(err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, &cartInputed)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err := w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Errorf(err.Error())
		}
		return
	}

	_, err = govalidator.ValidateStruct(cartInputed)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Errorf(err.Error())
		}
		return
	}

	advert, err := ch.advertUsecase.GetAdvert(cartInputed.AdvertId, userId, false)
	if err != nil {
		logger.Warnf("unable to get the advert: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err := w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write response to body: %s", err.Error())
		}
		return
	}

	_, err = ch.cartUsecase.UpdateCart(userId, &cartInputed, advert.Amount)
	if err != nil {
		logger.Warnf("unable to update the cart: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "successfully updated", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// UpdateAllCartHandler godoc
// @Summary Update all cart
// @Description Update all cart
// @Tags cart
// @Accept application/json
// @Produce application/json
// @Param add_to_cart body []models.CartHandler true "Add to Cart models"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyCartAll}
// @failure default {object} models.HttpError
// @Router /cart [post]
func (ch *CartHandler) UpdateAllCartHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	cartInputed := models.CHs{}
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

	err = easyjson.Unmarshal(buf, &cartInputed)
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

	for _, elCart := range cartInputed {
		_, err = govalidator.ValidateStruct(elCart)
		if err != nil {
			logger.Warnf("invalid data: %s", err.Error())
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
			if err != nil {
				logger.Warnf("cannot write answer to body %s", err.Error())
			}
			return
		}
	}

	adverts := make([]*models.Advert, 0, len(cartInputed))
	for _, element := range cartInputed {
		advert, err := ch.advertUsecase.GetAdvert(element.AdvertId, userId, false)
		if err != nil {
			logger.Warnf("unable to get the advert: %s", err.Error())
			w.WriteHeader(http.StatusOK)
			metaCode, metaMessage := internalError.ToMetaStatus(err)
			_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
			if err != nil {
				logger.Warnf("cannot write answer to body %s", err.Error())
			}
			return
		}

		adverts = append(adverts, advert)
	}

	cart, advs, messages, err := ch.cartUsecase.UpdateAllCart(userId, cartInputed, adverts)
	if err != nil {
		logger.Warnf("unable to update the cart: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyCartAll{Cart: cart, Adverts: advs, Hints: messages}
	_, err = w.Write(models.ToBytes(http.StatusOK, "successfully updated", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// GetCartHandler godoc
// @Summary Get user's cart
// @Description Get user's cart
// @Tags cart
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyCart}
// @failure default {object} models.HttpError
// @Router /cart [get]
func (ch *CartHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	cart, err := ch.cartUsecase.GetCart(userId)
	if err != nil {
		logger.Warnf("unable to get the cart: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	archived_advert := make([]int64, 0)
	adverts := make([]*models.Advert, 0)
	for _, e := range cart {
		advert, err := ch.advertUsecase.GetAdvert(e.AdvertId, userId, false)
		if err == internalError.EmptyQuery {
			archived_advert = append(archived_advert, e.AdvertId)
		} else if err != nil {
			logger.Warnf("unable to get the advert: %s", err.Error())
			w.WriteHeader(http.StatusOK)
			metaCode, metaMessage := internalError.ToMetaStatus(err)
			_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
			if err != nil {
				logger.Warnf("cannot write answer to body %s", err.Error())
			}
			return
		}

		adverts = append(adverts, advert)
	}

	if len(archived_advert) != 0 {
		s := ""
		for _, i := range archived_advert {
			s += strconv.FormatInt(i, 10) + " "
		}

		logger.Warnf("cart contain archived advert")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(models.ToBytes(http.StatusConflict, "cart contain archived advert: "+s, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyCart{Cart: cart, Adverts: adverts}
	_, err = w.Write(models.ToBytes(http.StatusOK, "successfully updated", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// GetCartHandler godoc
// @Summary Get user's cart
// @Description Get user's cart
// @Tags cart
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /cart/clear [post]
func (ch *CartHandler) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	err := ch.cartUsecase.ClearAllCart(userId)
	if err != nil {
		logger.Warnf("unable to clear the cart: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "cart cleared", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}

// CheckoutHandler godoc
// @Summary Checkout
// @Description Checkout
// @Tags cart
// @Accept application/json
// @Produce application/json
// @Param id path integer true "Advert id"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyOrder}
// @failure default {object} models.HttpError
// @Router /cart/{id}/checkout [post]
func (ch *CartHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
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

	order, err := ch.cartUsecase.GetOrderFromCart(userId, advertId)
	if err != nil {
		logger.Warnf("error with getting order: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	advert, err := ch.advertUsecase.GetAdvert(advertId, userId, false)
	if err != nil {
		logger.Warnf("error with getting advert: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	salesman, err := ch.userUsecase.GetById(advert.PublisherId)
	if err != nil {
		logger.Warnf("error with getting salesman: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = ch.cartUsecase.MakeOrder(order, advert)
	if err != nil {
		logger.Warnf("can not make order: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyOrder{Salesman: *salesman, Order: *order}
	_, err = w.Write(models.ToBytes(http.StatusOK, "order made successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}
