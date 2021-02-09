/*================================================================
*
*  文件名称：image_op.go
*  创 建 者: tiantingting
*  创建日期：2021年02月05日
*
================================================================*/
package imagesvc

import (
	sdk "github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	log "github.com/sirupsen/logrus"
	"iaas-api-server/common"
	"iaas-api-server/proto/image"
)

type Op interface {
	// Predo
	Predo() error
	// Do
	Do() error
	// Done
	Done(error) (interface{}, error)
}

type BasicOp struct {
	provider *sdk.ProviderClient
}

func getImageOpenStackClient( provider *sdk.ProviderClient ) (*sdk.ServiceClient,error) {
	sc, serviceErr := openstack.NewImageServiceV2(provider, sdk.EndpointOpts{})
	if serviceErr != nil {
		log.WithFields(log.Fields{
			"err": serviceErr,
		}).Error("get identity failed.")
		return nil,common.ETTGETIDENTITYCLIENT
	}
	return sc,common.EOK
}

type GetImageInfoOp struct {
	BasicOp
	Req *image.GetImageReq
	Res *image.GetImageRes
}
func (o *GetImageInfoOp) Predo() error {
	// check params
	if o.Req == nil {
		o.Res.Code=common.EPARAM.Code
		o.Res.Msg=common.EPARAM.Msg
		return common.EPARAM
	}
	o.Res=&image.GetImageRes{}
	var err error
	o.provider,err=common.GetOpenstackClient(o.Req.Apikey,o.Req.TenantId,o.Req.PlatformUserid)
	if err !=nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get identity auth failed.")
		o.Res.Msg=common.ETTGETIDENTITYCLIENT.Msg
		o.Res.Code=common.ETTGETIDENTITYCLIENT.Code
		return common.ETTGETIDENTITYCLIENT
	}
	return common.EOK
}
func (o *GetImageInfoOp) Do() error {
	sc,err:=getImageOpenStackClient(o.provider)
	if err!=common.EOK {
		return common.ETTGETIDENTITYCLIENT
	}
	imageResult,err:=images.Get(sc,o.Req.ImageId).Extract()
	if err!=nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get image info  failed.")
		o.Res.Code=common.EIGGETIMAGE.Code
		o.Res.Msg=common.EIGLISTIMAGES.Msg
		return common.EIGGETIMAGE
	}
	o.Res.Image = &image.Image{}
	o.Res.Image.ImageId=imageResult.ID
	o.Res.Image.ImageContainerformat=imageResult.ContainerFormat
	o.Res.Image.ImageDiskformat=imageResult.DiskFormat
	o.Res.Image.ImageName=imageResult.Name
	o.Res.Code=common.EOK.Code
	o.Res.Msg=common.EOK.Msg
	return common.EOK
}

func (o *GetImageInfoOp) Done(e error) (interface{}, error) {
	//Translate error code
	if e == common.EOK {
		return o.Res, nil
	}
	return o.Res, e
}

type ListImageInfoOp struct {
	BasicOp
	Req  *image.ListImagesReq
	Res *image.ListImagesRes
}

//Predo
func (o *ListImageInfoOp) Predo() error {
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res=new (image.ListImagesRes)
	var err error
	o.provider,err=common.GetOpenstackClient(o.Req.Apikey,o.Req.TenantId,o.Req.PlatformUserid)
	if err !=nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get identity auth failed.")
		return common.ETTGETIDENTITYCLIENT
	}
	return common.EOK
}
func (o *ListImageInfoOp) Do() error {
	startImageId:=o.Req.StartImageId
	pageSize:=o.Req.PageSize
	sc,err:=getImageOpenStackClient(o.provider)
	if err!=common.EOK {
		o.Res.Code=common.ETTGETIDENTITYCLIENT.Code
		o.Res.Msg=common.ETTGETIDENTITYCLIENT.Msg
		return common.ETTGETIDENTITYCLIENT
	}
	//TODO 根据租户ID获取projectID
	//projectID:=o.Req.TenantId
	var listOpts = images.ListOpts{
		//Owner:  projectID,
		Marker:startImageId,
		Limit:  int(pageSize),
	}
	if listOpts.Limit>1000 {
		listOpts.Limit=1000
	}
	lastAllImages, err3, done1 := getListImages(err, sc, listOpts)
	if done1 {
		o.Res.Code=common.EIGLISTIMAGES.Code
		o.Res.Msg=common.EIGLISTIMAGES.Msg
		return err3
	}
	getResult(lastAllImages,o)
	return common.EOK
}

func getResult(lastAllImages []images.Image, o *ListImageInfoOp) {
	o.Res.Images = []*image.Image{}
	for _, img := range lastAllImages {
		imageRes := &image.Image{ImageId: img.ID, ImageName: img.Name, ImageDiskformat: img.DiskFormat, ImageContainerformat: img.ContainerFormat}
		o.Res.Images = append(o.Res.Images, imageRes)
	}
	o.Res.Code = common.EOK.Code
	o.Res.Msg = common.EOK.Msg
}

func getListImages(err error, sc *sdk.ServiceClient, listOpts images.ListOpts) ([]images.Image, error, bool) {
	allPages, err := images.List(sc, listOpts).AllPages()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("create project failed.")
		return nil, common.EIGLISTIMAGES, true
	}
	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("create project failed.")
		return nil, common.EIGGETIMAGE, true
	}
	return allImages, nil, false
}

func (o *ListImageInfoOp) Done(e error) (interface{}, error) {
	//Translate error code
	if e == common.EOK {
		return o.Res, nil
	}
	return o.Res, e
}







