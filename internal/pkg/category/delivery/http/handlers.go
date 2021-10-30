package http

import (
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/category"
	"yula/internal/pkg/logging"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryUsecase category.CategoryUsecase
	logger          logging.Logger
}

func NewCategoryHandler(categoryUsecase category.CategoryUsecase, logger logging.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
		logger:          logger,
	}
}

func (ch *CategoryHandler) Routing(r *mux.Router) {
	s := r.PathPrefix("/category").Subrouter()
	s.HandleFunc("", ch.CategoriesListHandler).Methods(http.MethodGet, http.MethodOptions)
}

// CategoriesListHandler godoc
// @Summary Get list of all categories
// @Description Get list of all categories
// @Tags category
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface{body=[]models.Category}
// @failure default {object} models.HttpError
// @Router /category [get]
func (ch CategoryHandler) CategoriesListHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := ch.categoryUsecase.GetCategories()
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyCategories{Categories: categories}
	w.Write(models.ToBytes(http.StatusOK, "categories got successfully", body))
}
