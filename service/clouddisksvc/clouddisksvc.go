package clouddisksvc

import (
	//"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	cinder_op "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	cinder "github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	nova_op "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	log "github.com/sirupsen/logrus"
	"iaas-api-server/common/config"
	"iaas/configmap"

	//"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/clouddisk"
	"time"
)

type CloudDiskService struct {
}

//创建云硬盘
func (clouddisktask *CloudDiskService) CreateCloudDisk(req *clouddisk.CreateCloudDiskReq) (*clouddisk.CloudDiskRes, error) {

	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: req.Region,
	})

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack NewBlockStorageV3 failed"
		log.Error("openstack NewBlockStorageV3 failed: ", err)
		return res, err
	}

	ret, err := cinder.Create(client, cinder.CreateOpts{
		Size:             int(req.CloudDiskConf.SizeInG), //类型不一致，使用强转
		Name:             req.VolumeName,
		Description:      req.VolumeDesc,
		AvailabilityZone: req.AvailabilityZone,
		VolumeType:       req.CloudDiskConf.VolumeType,
	}).Extract()

	if err != nil {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack cinder create failed"
		log.Error("openstack cinder create failed: ", err)
		return res, err
	}

	res.CloudDisk.VolumeId = ret.ID
	res.CloudDisk.AvailabilityZone = ret.AvailabilityZone
	res.CloudDisk.VolumeName = ret.Name
	res.CloudDisk.CloudDiskConf.VolumeType = ret.VolumeType
	res.CloudDisk.CloudDiskConf.SizeInG = int32(ret.Size) //类型不一致，使用强转
	res.CloudDisk.VolumeDesc = ret.Description
	res.CloudDisk.Region = req.Region
	res.CloudDisk.VolumeStatus = ret.Status
	res.CloudDisk.CreatedTime = ret.CreatedAt.String()
	res.CloudDisk.UpdatedTime = ret.UpdatedAt.String()
	res.CloudDisk.AttachInstanceId = ret.Attachments[0].ServerID
	res.CloudDisk.AttachInstanceDevice = ret.Attachments[0].Device
	res.CloudDisk.AttachedTime = ret.Attachments[0].AttachedAt.String()

	return res, err
}

//删除云硬盘
func (clouddisktask *CloudDiskService) DeleteCloudDisk(req *clouddisk.DeleteCloudDiskReq) (*clouddisk.DeleteCloudDiskRes, error) {
	res := &clouddisk.DeleteCloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack NewBlockStorageV3 failed"
		log.Error("openstack NewBlockStorageV3 failed: ", err)
		return res, err
	}

	cascade, _ := config.GetString("Cascade")

	err = cinder.Delete(client, req.VolumeId, cinder.DeleteOpts{
		Cascade: cascade,
	}).ExtractErr()

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack cinder delete failed"
		log.Error("openstack cinder delete failed: ", err)
		return res, err
	}

	res.VolumeId = req.VolumeId
	res.DeletedTime = time.Now().String()

	return res, err
}

//获取云硬盘信息
func (clouddisktask *CloudDiskService) GetCloudDisk(req *clouddisk.GetCloudDiskReq) (*clouddisk.CloudDiskRes, error) {

	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack NewBlockStorageV3 failed"
		log.Error("openstack NewBlockStorageV3 failed: ", err)
		return res, err
	}

	ret, err := cinder.Get(client, req.VolumeId).Extract()

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack get clouddisk info failed"
		log.Error("openstack get clouddisk info failed: ", err)
		return res, err
	}

	res.CloudDisk.VolumeId = ret.ID
	res.CloudDisk.VolumeName = ret.Name
	res.CloudDisk.VolumeDesc = ret.Description
	res.CloudDisk.VolumeStatus = ret.Status
	res.CloudDisk.CreatedTime = ret.CreatedAt.String()
	res.CloudDisk.AvailabilityZone = ret.AvailabilityZone
	res.CloudDisk.CloudDiskConf.VolumeType = ret.VolumeType
	res.CloudDisk.CloudDiskConf.SizeInG = int32(ret.Size) //使用强转
	res.CloudDisk.UpdatedTime = ret.UpdatedAt.String()
	res.CloudDisk.AttachInstanceId = ret.Attachments[0].ServerID
	res.CloudDisk.AttachInstanceDevice = ret.Attachments[0].Device
	res.CloudDisk.AttachedTime = ret.Attachments[0].AttachedAt.String()

	return res, err
}

//云硬盘扩容
func (clouddisktask *CloudDiskService) ReqizeCloudDisk(req *clouddisk.ReqizeCloudDiskReq) (*clouddisk.CloudDiskRes, error) {

	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack NewBlockStorageV3 failed"
		log.Error("openstack NewBlockStorageV3 failed: ", err)
		return res, err
	}

	err = cinder_op.ExtendSize(client, req.VolumeId, cinder_op.ExtendSizeOpts{
		NewSize: int(req.CloudDiskConf.SizeInG), //类型不匹配，使用强转
	}).ExtractErr()

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack cinder volume extract failed"
		log.Error("openstack cinder volume extract failed: ", err)
		return res, err
	}

	res.CloudDisk.CloudDiskConf.SizeInG = req.CloudDiskConf.SizeInG

	return res, err
}

//更新云硬盘信息
func (clouddisktask *CloudDiskService) ModifyCloudDiskInfo(req *clouddisk.ModifyCloudDiskInfoReq) (*clouddisk.CloudDiskRes, error) {

	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack NewBlockStorageV3 failed"
		log.Error("openstack NewBlockStorageV3 failed: ", err)
		return res, err
	}

	ret, err := cinder.Update(client, req.VolumeId, cinder.UpdateOpts{
		Name:        &req.VolumeName,
		Description: &req.VolumeDesc,
	}).Extract()

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack cinder volume info update failed"
		log.Error("openstack cinder volume info update failed: ", err)
		return res, err
	}

	res.CloudDisk.VolumeDesc = ret.Description
	res.CloudDisk.VolumeName = ret.Name

	return res, err
}

//云主机挂载、卸载
func (clouddisktask *CloudDiskService) OperateCloudDisk(req *clouddisk.OperateCloudDiskReq) (*clouddisk.CloudDiskRes, error) {

	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = xxx //todo 错误码待定义
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	if req.OpsType == "Attach" {
		ret, err := nova_op.Create(client, req.InstanceId, nova_op.CreateOpts{
			//Device:   req.,
			VolumeID: req.VolumeId,
		}).Extract()

		if err != nil {
			res.Code = xxx //todo 错误码待定义
			res.Msg = "openstack nova attach failed"
			log.Error("openstack nova attach failed: ", err)
			return res, err
		}

		res.CloudDisk.VolumeId = ret.VolumeID
		res.CloudDisk.AttachInstanceId = ret.ServerID
		res.CloudDisk.AttachInstanceDevice = ret.Device
		res.CloudDisk.AttachedTime = time.Now().String()

	} else {
		err = nova_op.Delete(client, req.InstanceId, req.VolumeId).ExtractErr()

		if err != nil {
			res.Code = xxx //todo 错误码待定义
			res.Msg = "openstack nova detach failed"
			log.Error("openstack nova detach failed: ", err)
			return res, err
		}

		res.CloudDisk.VolumeId = req.VolumeId
	}
	return res, err
}
