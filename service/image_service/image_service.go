package image_service

import (
	"iaas-api-server/proto/image"
	"golang.org/x/net/context"
	"unicode"
)


type ImageService struct {}

func (is *ImageService) GetImage(context.Context, *image.GetImageReq) (*image.Image, error) {
	//your logic
}

func (is *ImageService) ListImages(context.Context, *image.ListImagesReq) (*image.ListImagesRes, error) {
	//your logic
}
