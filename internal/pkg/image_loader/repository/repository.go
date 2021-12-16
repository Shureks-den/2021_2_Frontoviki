package repository

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	internalError "yula/internal/error"
	imageloader "yula/internal/pkg/image_loader"
)

type ImageLoaderRepository struct {
}

func NewImageLoaderRepository() imageloader.ImageLoaderRepository {
	return &ImageLoaderRepository{}
}

func (ilr *ImageLoaderRepository) Insert(fileHeader *multipart.FileHeader, dir string, name string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return internalError.UnableToReadFile
	}
	defer file.Close()

	newFile, err := os.Create(dir + "/" + name)
	if err != nil {
		return internalError.InternalError
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return internalError.InternalError
	}

	return nil
}

func (ilr *ImageLoaderRepository) Delete(filePath string) error {
	if filePath == "" {
		return nil
	}

	err := os.Remove(fmt.Sprintf("%s.%s", filePath, imageloader.CompressedFormat))
	if err != nil {
		parts := strings.Split(filePath, "__")
		if len(parts) != 2 {
			return internalError.UnableToRemove
		}
		err = os.Remove(fmt.Sprintf("%s.%s", filePath, parts[1]))
		if err != nil {
			return internalError.UnableToRemove
		}
	}
	return nil
}
