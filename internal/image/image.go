package image

import (
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
	imageapi "github.com/junaidk/image-service"
)

func GetMetadata(f io.ReadSeeker, fileExt string) (*imageapi.Image, error) {

	var data exif2.Exif
	var err error

	data, err = imagemeta.Decode(f)

	if err != nil {
		return nil, err
	}

	fmt.Println(data)
	imgData := imageapi.Image{
		MetaData: imageapi.MetaData{
			Type: data.ImageType.Extension(),
			Size: imageapi.Size{
				Width:  data.ImageWidth,
				Height: data.ImageHeight,
			},
			CameraModel: data.Make + "-" + data.Model,
			LensModel:   data.LensModel,
			Location: imageapi.Location{
				Latitude:  data.GPS.Latitude(),
				Longitude: data.GPS.Longitude(),
			},
		},
	}

	return &imgData, nil
}
