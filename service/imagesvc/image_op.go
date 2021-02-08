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

type Authorization interface {
	Auth() bool
}

type OpenstackAPIAuthorization struct {
	Apikey         string
	TenantId       string
	PlatformUserid string
}

func (o *OpenstackAPIAuthorization) Auth() bool {
	//common.APIAuth(o.Apikey, o.TenantId, o.PlantformUserid)
	return true
}

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
	provider *sdk.ProviderClient
}
func (o *GetImageInfoOp) Predo() error {
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res=&image.GetImageRes{}
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
		return common.EIGGETIMAGE
	}
	o.Res.Image = &image.Image{}
	o.Res.Image.ImageId=imageResult.ID
	o.Res.Image.ImageContainerformat=imageResult.ContainerFormat
	o.Res.Image.ImageDiskformat=imageResult.DiskFormat
	o.Res.Image.ImageName=imageResult.Name
	o.Res.Code=200
	o.Res.Msg="获取镜像信息成功"
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
	pageNum:=o.Req.PageNumber
	pageSize:=o.Req.PageSize
	sc,err:=getImageOpenStackClient(o.provider)
	if err!=common.EOK {
		return common.ETTGETIDENTITYCLIENT
	}
	//TODO 根据租户ID获取projectID
	//projectID:=o.Req.TenantId
	//当页数为1时，直接取数据
	if pageNum==1 {
		var listOpts = images.ListOpts{
			//Owner: projectID,
			Marker: "",
			Limit: int(pageSize),
		}
		lastAllImages, err, done := getListImages(err, sc, listOpts)
		if done {
			return err
		}
		getResult(lastAllImages, o)
		return common.EOK
	}
	//当页数大于1时，需要分两步获取数据
	//获取maker
	var listOptsGetMaker = images.ListOpts{
		//Owner:  projectID,
		Marker: "",
		Limit:  int(pageSize * (pageNum - 1)),
	}
	allImages, err2, done := getListImages(err, sc, listOptsGetMaker)
	if done {
		return err2
	}
	size := len(allImages)
	id := allImages[size-1].ID

	var listOpts = images.ListOpts{
		//Owner:  projectID,
		Marker: id,
		Limit:  int(pageSize * (pageNum - 1)),
	}
	lastAllImages, err3, done1 := getListImages(err, sc, listOpts)
	if done1 {
		return err3
	}
	getResult(lastAllImages,o)
	return common.EOK
}

func getResult(lastAllImages []images.Image, o *ListImageInfoOp) {
	var imageSlice = []*image.Image{}
	for _, img := range lastAllImages {
		imageRes := &image.Image{ImageId: img.ID, ImageName: img.Name, ImageDiskformat: img.DiskFormat, ImageContainerformat: img.ContainerFormat}
		imageSlice = append(imageSlice, imageRes)
	}
	o.Res.Images = []*image.Image{}
	o.Res.Images = imageSlice
	o.Res.Code = 200
	o.Res.Msg = "获取镜像列表成功"
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







