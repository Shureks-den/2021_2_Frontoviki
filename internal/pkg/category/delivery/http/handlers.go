package http

import (
	"context"
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"

	proto "yula/proto/generated/category"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CategoryHandler struct {
	categoryUsecase proto.CategoryClient
}

func NewCategoryHandler(categoryUsecase proto.CategoryClient) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
	}
}

var (
	logger logging.Logger = logging.GetLogger()
)

func (ch *CategoryHandler) Routing(r *mux.Router) {
	s := r.PathPrefix("/category").Subrouter()
	s.HandleFunc("", middleware.SetSCRFToken(http.HandlerFunc(ch.CategoriesListHandler))).Methods(http.MethodGet, http.MethodOptions)
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
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	protocategories, err := ch.categoryUsecase.GetCategories(context.Background(), &proto.Nothing{Dummy: true})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	var categories []*models.Category
	for _, category := range protocategories.Categories {
		categories = append(categories, &models.Category{
			Name: category.Name,
		})
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyCategories{Categories: categories}
	w.Write(models.ToBytes(http.StatusOK, "categories got successfully", body))
}
