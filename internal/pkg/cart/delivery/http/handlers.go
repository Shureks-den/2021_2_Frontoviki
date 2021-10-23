package http

import (
	"encoding/json"
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/cart"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CartHandler struct {
	cartUsecase   cart.CartUsecase
	userUsecase   user.UserUsecase
	advertUsecase advt.AdvtUsecase
	logger        logging.Logger
}

func NewCartHandler(cartUsecase cart.CartUsecase, userUsecase user.UserUsecase, advertUsecase advt.AdvtUsecase, logger logging.Logger) *CartHandler {
	return &CartHandler{
		cartUsecase:   cartUsecase,
		userUsecase:   userUsecase,
		advertUsecase: advertUsecase,
		logger:        logger,
	}
}

func (ch *CartHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/cart").Subrouter()
	s.Use(sm.CheckAuthorized)

	s.HandleFunc("/one", ch.UpdateOneAdvertHandler).Methods(http.MethodPost, http.MethodOptions)
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
	ch.logger = ch.logger.GetLoggerWithFields((r.Context().Value("logger fields")).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	var cartInputed models.CartHandler
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&cartInputed)
	if err != nil {
		ch.logger.Warnf("invalid body: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(cartInputed)
	if err != nil {
		ch.logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	advert, err := ch.advertUsecase.GetAdvert(cartInputed.AdvertId)
	if err != nil {
		ch.logger.Warnf("unable to get the advert: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = ch.cartUsecase.UpdateCart(userId, &cartInputed, advert.Amount)
	if err != nil {
		ch.logger.Warnf("unable to update the cart: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "successfully updated", nil))
}
