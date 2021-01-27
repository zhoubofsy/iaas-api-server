package image_service

import (
	"iaas-api-server/proto/image"
	"golang.org/x/net/context"
	"unicode"
)


type GetImageSvc struct {}

func (gi *GetImageSvc) GetImage(context.Context, *image.GetImageReq) (*image.Image, error) {
	//your logic
}
