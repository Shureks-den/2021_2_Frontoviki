package models

import "time"

const (
	avatarsDirectory     string = "static/avatars/"
	advertImageDirectory string = "static/advertimages/"
)

type ImageCompressorConfig struct {
	PathToDir  []string
	Start      time.Time
	EveryHours int
}

func NewImageCompressorConfig(t time.Time, repeatInHour int) *ImageCompressorConfig {
	return &ImageCompressorConfig{
		[]string{avatarsDirectory, advertImageDirectory},
		t,
		repeatInHour,
	}
}
