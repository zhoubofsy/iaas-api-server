/*================================================================
*
*  文件名称：image_service.go
*  创 建 者: tiantingting
*  创建日期：2021年02月05日
*
================================================================*/
package imagesvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/image"
)

type ImageService struct{
	image.UnimplementedImageServiceServer
}

func (is *ImageService) GetImage(cxt context.Context,r *image.GetImageReq) (*image.GetImageRes, error) {
	op := new(GetImageInfoOp)
	op.Req = r
	res, err := is.Process(op)
	return res.(*image.GetImageRes), err
}

func (is *ImageService) ListImages(cxt context.Context,r *image.ListImagesReq) (*image.ListImagesRes, error) {
	op := new(ListImageInfoOp)
	op.Req = r
	res, err := is.Process(op)
	return res.(*image.ListImagesRes), err
}

func (o *ImageService) Process( op Op) (interface{}, error) {
	err := op.Predo()
	if err == common.EOK {
		err = op.Do()
	}
	return op.Done(err)
}

