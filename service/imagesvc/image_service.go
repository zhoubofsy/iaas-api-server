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
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(GetImageInfoOp)

	op.Req = r
	res, err := is.Process(auth, op)
	return res.(*image.GetImageRes), err
}

func (is *ImageService) ListImages(cxt context.Context,r *image.ListImagesReq) (*image.ListImagesRes, error) {
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(ListImageInfoOp)

	op.Req = r
	res, err := is.Process(auth, op)
	return res.(*image.ListImagesRes), err
}

func (o *ImageService) Process(auth Authorization, op Op) (interface{}, error) {
	err := op.Predo()
	if err == common.EOK {
		if auth.Auth() == false {
			err = common.EUNAUTHORED
		} else {
			err = op.Do()
		}
	}
	return op.Done(err)
}

