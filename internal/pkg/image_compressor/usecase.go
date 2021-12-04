package imagecompressor

type ImageCompressorUsecase interface {
	HandleImages(folders []string) error
}
