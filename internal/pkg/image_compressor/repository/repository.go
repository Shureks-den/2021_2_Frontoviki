package repository

import (
	"fmt"
	"os"
	internalError "yula/internal/error"
	imagecompressor "yula/internal/pkg/image_compressor"

	"github.com/h2non/bimg"
)

const CompressedFormat = "webp"

type ImageCompressorRepository struct {
}

func NewImageCompressorRepository() imagecompressor.ImageCompressorRepository {
	return &ImageCompressorRepository{}
}

func (icr *ImageCompressorRepository) Compress(pathToDir, filename, extension string) error {
	path := fmt.Sprintf("%s%s.%s", pathToDir, filename, extension)
	buffer, err := bimg.Read(path)
	if err != nil {
		return internalError.ImageNotExist
	}

	newImage, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return internalError.UnableToConvert
	}

	if bimg.NewImage(newImage).Type() != CompressedFormat {
		return internalError.NotConverted
	}

	path = fmt.Sprintf("%s%s.%s", pathToDir, filename, CompressedFormat)
	err = bimg.Write(path, newImage)
	return err
}

func (icr *ImageCompressorRepository) Delete(dirpath, filename string) error {
	path := fmt.Sprintf("%s%s", dirpath, filename)
	err := os.Remove(path)
	if err != nil {
		return internalError.UnableToRemove
	}
	return nil
}
