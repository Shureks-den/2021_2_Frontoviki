package imagecompressor

type ImageCompressorRepository interface {
	Compress(filepath, filename, extension string) error
	Delete(dirpath, filename string) error
}
