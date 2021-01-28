package imagesvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/proto/image"
	"unicode"
)

type ImageService struct{}

func (is *ImageService) GetImage(context.Context, *image.GetImageReq) (*image.Image, error) {
	//your logic
}

func (is *ImageService) ListImages(context.Context, *image.ListImagesReq) (*image.ListImagesRes, error) {
	//your logic
}
