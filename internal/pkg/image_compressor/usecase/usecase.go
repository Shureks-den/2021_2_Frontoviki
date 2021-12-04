package usecase

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
	"yula/internal/models"
	imagecompressor "yula/internal/pkg/image_compressor"
	"yula/internal/pkg/logging"

	"github.com/go-co-op/gocron"
)

var (
	compressedFormat     = "webp"
	notCompressedFormats = []string{"jpg", "jpeg", "png"}
)

var logger logging.Logger = logging.GetLogger()

type ImageCompressorUsecase struct {
	imageCompressorRepository imagecompressor.ImageCompressorRepository
	sheduler                  *gocron.Scheduler
}

func NewImageCompressorUsecase(icr imagecompressor.ImageCompressorRepository) imagecompressor.ImageCompressorUsecase {
	return &ImageCompressorUsecase{
		imageCompressorRepository: icr,
		sheduler:                  nil,
	}
}

func (icu *ImageCompressorUsecase) StartByCron(config *models.ImageCompressorConfig) error {
	icu.sheduler = gocron.NewScheduler(config.Start.UTC().Location())
	_, err := icu.sheduler.Every(config.EveryHours).Minute().Do(icu.HandleImages, config.PathToDir)
	icu.sheduler.StartAsync()
	return err
}

func (icu *ImageCompressorUsecase) HandleImages(folders []string) error {
	for _, folder := range folders {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			logger.Errorf("cannot read path: %s\nerror: %v", folder, err)
			return err
		}

		for _, file := range files {
			_, isUpdated, err := icu.CompressImage(file, folder)
			if err != nil {
				logger.Errorf("cannot compress image: %s in folder: %s\nerror: %v", file.Name(), folder, err)
				return err
			}

			if isUpdated {
				err = icu.imageCompressorRepository.Delete(folder, file.Name())
				if err != nil {
					logger.Errorf("cannot delete image: %s in folder: %s\nerror: %v", file.Name(), folder, err)
					return err
				}
			}
		}
	}

	return nil
}

func searchInSlice(slice []string, value string) bool {
	for _, elem := range slice {
		if value == elem {
			return true
		}
	}
	return false
}

func (icu *ImageCompressorUsecase) CompressImage(file fs.FileInfo, folder string) (string, bool, error) {
	fileParts := strings.Split(file.Name(), ".")
	filename, extension := fileParts[0], fileParts[1]
	newFilename := ""
	if searchInSlice(notCompressedFormats, extension) {
		err := icu.imageCompressorRepository.Compress(folder, filename, extension)
		if err != nil {
			return "", true, err
		}

		newFilename = fmt.Sprintf("%s.%s", filename, compressedFormat)
		return newFilename, true, err
	}
	return "", false, nil
}

func (icu *ImageCompressorUsecase) StopJob() {
	if icu.sheduler != nil && icu.sheduler.IsRunning() {
		icu.sheduler.Stop()
		icu.sheduler.Clear()
		fmt.Println("stop job")
	}
}
