package usecase

import (
	"testing"
	"yula/services/category/mocks"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	mocksRepository := mocks.CategoryRepository{}
	categoryUsecase := NewCategoryUsecase(&mocksRepository)
	mocksRepository.On("SelectCategories").Return(nil, nil)
	categories, err := categoryUsecase.GetCategories()

	assert.Nil(t, categories)
	assert.Nil(t, err)
}
