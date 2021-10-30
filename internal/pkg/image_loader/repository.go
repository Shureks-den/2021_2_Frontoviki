package imageloader

import "mime/multipart"

//go:generate mockery -name=ImageLoaderRepository

type ImageLoaderRepository interface {
	Insert(fileHeader *multipart.FileHeader, dir string, name string) error
	Delete(filePath string) error
}
