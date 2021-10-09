package imageloader

import "mime/multipart"

const (
	AvatarsDirectory string = "static/avatars"
	DefaultAvatar    string = AvatarsDirectory + "/default_avatar.png"
)

type ImageLoaderUsecase interface {
	UploadAvatar(headerFile *multipart.FileHeader) (string, error)
	RemoveAvatar(filePath string) error
}
