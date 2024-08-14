package image

import (
	"io"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
	imageapi "github.com/junaidk/image-service"
)

func GetMetadata(f io.ReadSeeker, fileExt string) (*imageapi.Image, error) {

	var data exif2.Exif
	var err error

	switch fileExt {
	case ".png":
		data, err = imagemeta.DecodePng(f)
	default:
		data, err = imagemeta.Decode(f)
	}

	if err != nil {
		return nil, err
	}

	imgData := imageapi.Image{
		MetaData: imageapi.MetaData{
			Type: data.ImageType.Extension(),
			Size: imageapi.Size{
				Width:  data.ImageWidth,
				Height: data.ImageHeight,
			},
			CameraModel: data.CameraMake.String(),
			LensModel:   data.LensModel,
			Location: imageapi.Location{
				Latitude:  data.GPS.Latitude(),
				Longitude: data.GPS.Longitude(),
			},
		},
	}

	return &imgData, nil
}
