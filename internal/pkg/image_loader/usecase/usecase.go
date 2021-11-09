package usecase

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"
	internalError "yula/internal/error"
	imageloader "yula/internal/pkg/image_loader"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type ImageLoaderUsecase struct {
	imageLoaderRepo imageloader.ImageLoaderRepository
}

func NewImageLoaderUsecase(imageLoaderRepo imageloader.ImageLoaderRepository) imageloader.ImageLoaderUsecase {
	return &ImageLoaderUsecase{
		imageLoaderRepo: imageLoaderRepo,
	}
}

func (ilu *ImageLoaderUsecase) Upload(headerFile *multipart.FileHeader, dir string) (string, error) {
	availableFormats := []string{"png", "jpeg"}

	ct := headerFile.Header.Get("Content-Type")

	extension := ct[strings.LastIndex(ct, "/")+1:]
	if !contains(availableFormats, extension) {
		return "", internalError.UnknownExtension
	}

	currentTime := time.Now()
	avatarId := fmt.Sprintf("%d_%d_%d", currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond())
	filename := avatarId + "." + extension

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

func (ilu *ImageLoaderUsecase) RemoveAdvertImages(imageUrls []string) error {
	for _, url := range imageUrls {
		if url == imageloader.DefaultAdvertImage {
			continue
		}
		err := ilu.imageLoaderRepo.Delete(url)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ilu *ImageLoaderUsecase) UploadAdvertImages(headerFiles []*multipart.FileHeader) ([]string, error) {
	var urls []string
	for _, file := range headerFiles {
		url, err := ilu.Upload(file, imageloader.AdvertImageDirectory)
		if err != nil {
			return urls, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
