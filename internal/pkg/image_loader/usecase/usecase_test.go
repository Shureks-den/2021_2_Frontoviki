package usecase

import (
	"mime/multipart"
	"net/textproto"
	"testing"
	imageloader "yula/internal/pkg/image_loader"
	ILMock "yula/internal/pkg/image_loader/mocks"

	myerr "yula/internal/error"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadSuccess(t *testing.T) {
	ilr := ILMock.ImageLoaderRepository{}
	ilu := NewImageLoaderUsecase(&ilr)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "image/png")
	file := multipart.FileHeader{
		Filename: "aboba.png",
		Header:   partHeaders,
		Size:     28915,
	}

	ilr.On("Insert", &file, ".", mock.AnythingOfType("string")).Return(nil)

	_, err := ilu.Upload(&file, ".")
	assert.Nil(t, err)
}

func TestUploadBadExtension(t *testing.T) {
	ilr := ILMock.ImageLoaderRepository{}
	ilu := NewImageLoaderUsecase(&ilr)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "aboba")
	file := multipart.FileHeader{
		Filename: "aboba.png",
		Header:   partHeaders,
		Size:     28915,
	}

	ilr.On("Insert", &file, ".", mock.AnythingOfType("string")).Return(myerr.UnknownExtension)

	_, err := ilu.Upload(&file, ".")
	assert.Equal(t, err, myerr.UnknownExtension)
}

func TestUploadAvatar(t *testing.T) {
	ilr := ILMock.ImageLoaderRepository{}
	ilu := NewImageLoaderUsecase(&ilr)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "image/png")
	file := multipart.FileHeader{
		Filename: "aboba.png",
		Header:   partHeaders,
		Size:     28915,
	}

	ilr.On("Insert", &file, imageloader.AvatarsDirectory, mock.AnythingOfType("string")).Return(nil)

	_, err := ilu.UploadAvatar(&file)
	assert.Nil(t, err)
}

func TestRemoveAvatar(t *testing.T) {
	ilr := ILMock.ImageLoaderRepository{}
	ilu := NewImageLoaderUsecase(&ilr)
	ilr.On("Delete", "./a.png").Return(nil)

	err := ilu.RemoveAvatar("./a.png")
	assert.Nil(t, err)
}

func TestRemoveAdImages(t *testing.T) {
	ilr := ILMock.ImageLoaderRepository{}
	ilu := NewImageLoaderUsecase(&ilr)
	imageUrls := []string{"./a.png"}
	ilr.On("Delete", "./a.png").Return(nil)

	err := ilu.RemoveAdvertImages(imageUrls)
	assert.Nil(t, err)
}

func TestUploadAdImages(t *testing.T) {
	ilr := ILMock.ImageLoaderRepository{}
	ilu := NewImageLoaderUsecase(&ilr)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "image/png")
	files := []*multipart.FileHeader{
		{
			Filename: "a.png",
			Header:   partHeaders,
			Size:     28915,
		},
	}
	ilr.On("Insert", files[0], imageloader.AdvertImageDirectory, mock.AnythingOfType("string")).Return(nil)

	_, err := ilu.UploadAdvertImages(files)
	assert.Nil(t, err)
}
