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

func (ilu *ImageLoaderUsecase) Upload(headerFile *multipart.FileHeader, dir string) (string, error) {
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

	err := ilu.imageLoaderRepo.Insert(headerFile, dir, filename)
	if err != nil {
		return "", err
	}

	return dir + "/" + filename, nil
}

func (ilu *ImageLoaderUsecase) UploadAvatar(headerFile *multipart.FileHeader) (string, error) {
	url, err := ilu.Upload(headerFile, imageloader.AvatarsDirectory)
	return url, err
}

func (ilu *ImageLoaderUsecase) RemoveAvatar(filePath string) error {
	return ilu.imageLoaderRepo.Delete(filePath)
}

func (ilu *ImageLoaderUsecase) UploadAdvertImages(headerFiles []*multipart.FileHeader) ([]string, error) {
	var urls []string
	for _, file := range headerFiles {
		url, err := ilu.Upload(file, imageloader.AdvertImageDirectory)
		if err != nil {
			return []string{}, err
		}

		urls = append(urls, url)
	}
	return urls, nil
}

func (ilu *ImageLoaderUsecase) RemoveAdvertImages(imageUrls []string) error {
	for _, url := range imageUrls {
		err := ilu.imageLoaderRepo.Delete(url)
		if err != nil {
			return err
		}
	}
	return nil
}
