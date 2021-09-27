package http

import (
	"encoding/json"
	"log"
	"net/http"
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
	r.HandleFunc("/", ah.AdvtListHandler).Methods(http.MethodGet, http.MethodOptions)
}

func (ah *AdvtHandler) AdvtListHandler(w http.ResponseWriter, r *http.Request) {
	advts, err := ah.advtUsecase.GetListAdvt(0, 100, true)
	log.Println("1")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		log.Println("2")
		response := models.HttpError{Code: http.StatusInternalServerError, Message: err.Error()}
		js, _ := json.Marshal(response)

		w.Write(js)
		log.Println("c")
		return
	}
	log.Println("d")
	w.WriteHeader(http.StatusOK)

	response := models.HttpBodyInterface{Code: http.StatusOK, Message: "list of ads found successfully",
		Body: models.HttpBodyAdvts{Advts: advts}}
	js, err := json.Marshal(response)
	if err != nil {
		log.Println("e")
		response := models.HttpError{Code: http.StatusInternalServerError, Message: err.Error()}
		js, _ = json.Marshal(response)
	}

	log.Println(response.Body)
	w.Write(js)
}
