package main

import (
	"yula/internal/config"

	imgcmprsRep "yula/internal/pkg/image_compressor/repository"
	imgcmprsUse "yula/internal/pkg/image_compressor/usecase"
	"yula/internal/pkg/logging"
)

func main() {
	logger := logging.GetLogger()

	if err := config.LoadConfig(); err != nil {
		logger.Errorf("cannot load config: %s", err.Error())
		return
	}

	folders := config.Cfg.GetStaticDirs()

	icr := imgcmprsRep.NewImageCompressorRepository()
	icu := imgcmprsUse.NewImageCompressorUsecase(icr)

	if err := icu.HandleImages(folders); err != nil {
		logger.Errorf("serve compressor error: %v", err)
	}
}
