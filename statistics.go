package imageapi

import "context"

type Statistics struct {
	TopImageFormat       ImageFormatItem   `json:"top_image_format"`
	TopCameraModel       []CameraModelItem `json:"top_camera_model"`
	ImageUploadFrequency []FrequencyItem   `json:"image_upload_frequency"`
}

type ImageFormatItem struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type CameraModelItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type FrequencyItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type StatisticsService interface {
	GetStatistics(context.Context) (*Statistics, error)
}
