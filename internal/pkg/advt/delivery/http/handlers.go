package http

import (
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"

	"github.com/gorilla/mux"
)

type AdvtHandler struct {
	advtUsecase advt.AdvtUsecase
}

func NewAdvtHandler(advtUsecase advt.AdvtUsecase) *AdvtHandler {
	return &AdvtHandler{
		advtUsecase: advtUsecase,
	}
}

func (ah *AdvtHandler) Routing(r *mux.Router) {
	s := r.PathPrefix("/adverts").Subrouter()
	s.HandleFunc("", ah.AdvtListHandler).Methods(http.MethodGet, http.MethodOptions)
}

func (ah *AdvtHandler) AdvtListHandler(w http.ResponseWriter, r *http.Request) {
	advts, err := ah.advtUsecase.GetListAdvt(0, 100, true)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyAdvts{Advts: advts}
	w.Write(models.ToBytes(http.StatusOK, "adverts found successfully", body))
}
