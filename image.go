package imageapi

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID       uuid.UUID
	Name     string
	Hash     string
	CreateAt time.Time
	MetaData MetaData
}

type MetaData struct {
	Type        string
	Size        Size
	CameraModel string
	LensModel   string
	Location    Location
}
type Size struct {
	Width  uint16
	Height uint16
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type ImageService interface {
	CreateImage(context.Context, *Image) error
	GetImage(context.Context, uuid.UUID) (*Image, error)
}
