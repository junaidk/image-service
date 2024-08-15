package image

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	imageapi "github.com/junaidk/image-service"
)

func TestGetMetadata(t *testing.T) {
	file, err := os.Open("../../sample/images/Canon_40D.jpg")
	if err != nil {
		t.Fatalf("error openinging image %v", err)
	}

	imgMeta, err := GetMetadata(file, ".jpg")
	if err != nil {
		t.Fatalf("error getting metadata %v", err)
	}

	expectedImage := &imageapi.Image{
		MetaData: imageapi.MetaData{
			Type: "jpg",
			Size: imageapi.Size{
				Width:  100, // Expected width
				Height: 68,  // Expected height
			},
			CameraModel: "Canon EOS 40D",
			LensModel:   "",
			Location: imageapi.Location{
				Latitude:  0,
				Longitude: 0,
			},
		},
	}

	if !cmp.Equal(imgMeta, expectedImage) {
		t.Fatalf("output is not equal to expected\n%s", cmp.Diff(imgMeta, expectedImage))
	}
}
