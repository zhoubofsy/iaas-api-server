package imagesvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/proto/image"
)

type ImageService struct{}

func (is *ImageService) GetImage(context.Context, *image.GetImageReq) (*image.Image, error) {
	//your logic
	return nil, nil
}

func (is *ImageService) ListImages(context.Context, *image.ListImagesReq) (*image.ListImagesRes, error) {
	//your logic
	return nil, nil
}
