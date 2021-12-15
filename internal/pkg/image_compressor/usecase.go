package imagecompressor

import "io/fs"

type ImageCompressorUsecase interface {
	HandleImages(folders []string) error
	CompressImage(file fs.FileInfo, folder string) (string, bool, error)
}
