package imageloader

import "mime/multipart"

type ImageLoaderRepository interface {
	Insert(fileHeader *multipart.FileHeader, dir string, name string) error
	Delete(filePath string) error
}
