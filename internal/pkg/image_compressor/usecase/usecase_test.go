package usecase

import (
	"testing"
	"yula/internal/pkg/image_compressor/mocks"

	"github.com/stretchr/testify/assert"
)

func TestHandleImagesError(t *testing.T) {
	icr := mocks.ImageCompressorRepository{}

	var folder string = "folder"
	var str string = "filename"

	icr.On("Compress", folder, str, "ext").Return(nil)
	icr.On("Delete", folder, str).Return(nil)

	icu := NewImageCompressorUsecase(&icr)
	err := icu.HandleImages([]string{folder})

	assert.Error(t, err)
}

func TestSearchInSliceOk(t *testing.T) {
	slice := []string{"a", "b", "c"}
	elem := "a"

	out := searchInSlice(slice, elem)
	assert.Equal(t, true, out)
}

func TestSearchInSliceError(t *testing.T) {
	slice := []string{"a", "b", "c"}
	elem := "d"

	out := searchInSlice(slice, elem)
	assert.Equal(t, false, out)
}
