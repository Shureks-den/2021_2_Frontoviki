package imageloader

import "mime/multipart"

const (
	AvatarsDirectory     string = "static/avatars"
	DefaultAvatar        string = AvatarsDirectory + "/default_avatar.png"
	AdvertImageDirectory string = "static/advertimages"
	DefaultAdvertImage   string = AdvertImageDirectory + "/default_image.png"
)

//go:generate mockery -name=ImageLoaderUsecase

type ImageLoaderUsecase interface {
	Upload(headerFile *multipart.FileHeader, dir string) (string, error)
	UploadAvatar(headerFile *multipart.FileHeader) (string, error)
	RemoveAvatar(filePath string) error

	UploadAdvertImages(headerFiles []*multipart.FileHeader) ([]string, error)
	RemoveAdvertImages(imageUrls []string) error
}
