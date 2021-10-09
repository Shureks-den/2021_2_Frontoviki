package usecase

import (
	"log"
	"mime/multipart"
	internalError "yula/internal/error"
	imageloader "yula/internal/pkg/image_loader"

	"github.com/google/uuid"
)

type ImageLoaderUsecase struct {
	imageLoaderRepo imageloader.ImageLoaderRepository
}

func NewImageLoaderUsecase(imageLoaderRepo imageloader.ImageLoaderRepository) imageloader.ImageLoaderUsecase {
	return &ImageLoaderUsecase{
		imageLoaderRepo: imageLoaderRepo,
	}
}

func (ilu *ImageLoaderUsecase) UploadAvatar(headerFile *multipart.FileHeader) (string, error) {
	extension := ""

	ct := headerFile.Header.Get("Content-Type")
	switch ct {
	case "image/png":
		extension = "png"
	case "image/jpeg":
		extension = "jpeg"
	default:
		return "", internalError.UnknownExtension
	}

	avatarId := uuid.NewString()
	filename := avatarId + "." + extension
	log.Println(filename)

	err := ilu.imageLoaderRepo.Insert(headerFile, imageloader.AvatarsDirectory, filename)
	if err != nil {
		return "", err
	}

	return imageloader.AvatarsDirectory + "/" + filename, nil
}

func (ilu *ImageLoaderUsecase) RemoveAvatar(filePath string) error {
	return ilu.imageLoaderRepo.Delete(filePath)
}
