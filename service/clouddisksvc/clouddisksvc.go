package clouddisksvc

import (
	//"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	cinder_op "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	cinder "github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	nova_op "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/common/config"
	"iaas-api-server/proto/clouddisk"
)

type CloudDiskService struct {
}

//创建云硬盘
func (clouddisktask *CloudDiskService) CreateCloudDisk(ctx context.Context, req *clouddisk.CreateCloudDiskReq) (*clouddisk.CloudDiskRes, error) {
	timer_elasp := common.NewTimer()
	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: req.Region,
	})
	if nil != err {
		res.Code = common.ENEWBLOCK.Code
		res.Msg = common.ENEWBLOCK.Msg
		log.Error(res.Msg, ": ", err)
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
		res.Code = common.ENEWVOLUME.Code
		res.Msg = common.ENEWVOLUME.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.CloudDisk = &clouddisk.CloudDiskRes_CloudDisk{
		VolumeId:         ret.ID,
		AvailabilityZone: ret.AvailabilityZone,
		VolumeName:       ret.Name,
		CloudDiskConf: &clouddisk.CloudDiskConf{
			VolumeType: ret.VolumeType,
			SizeInG:    int32(ret.Size), //类型不一致，使用强转
		},
		VolumeDesc:   ret.Description,
		Region:       req.Region,
		VolumeStatus: ret.Status,
		CreatedTime:  ret.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		UpdatedTime:  ret.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
	}

	log.Info("rpc CreateVolume: ", res, ". time elapse: ", timer_elasp.Elapse())
	return res, err
}

//删除云硬盘
func (clouddisktask *CloudDiskService) DeleteCloudDisk(ctx context.Context, req *clouddisk.DeleteCloudDiskReq) (*clouddisk.DeleteCloudDiskRes, error) {
	timer_elasp := common.NewTimer()
	res := &clouddisk.DeleteCloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})
	if nil != err {
		res.Code = common.ENEWBLOCK.Code
		res.Msg = common.ENEWBLOCK.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	cascade, _ := config.GetBool("Cascade")

	err = cinder.Delete(client, req.VolumeId, cinder.DeleteOpts{
		Cascade: cascade,
	}).ExtractErr()

	if nil != err {
		res.Code = common.EDELETEVOLUME.Code
		res.Msg = common.EDELETEVOLUME.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.VolumeId = req.VolumeId
	res.DeletedTime = common.Now()
	log.Info("rpc CreateVolume: ", res, ". time elapse: ", timer_elasp.Elapse())
	return res, err
}

//获取云硬盘信息
func (clouddisktask *CloudDiskService) GetCloudDisk(ctx context.Context, req *clouddisk.GetCloudDiskReq) (*clouddisk.CloudDiskRes, error) {
	timer_elasp := common.NewTimer()
	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = common.ENEWBLOCK.Code
		res.Msg = common.ENEWBLOCK.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	ret, err := cinder.Get(client, req.VolumeId).Extract()

	if nil != err {
		res.Code = common.ESHOWVOLUME.Code
		res.Msg = common.ESHOWVOLUME.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.CloudDisk = &clouddisk.CloudDiskRes_CloudDisk{
		VolumeId:         ret.ID,
		AvailabilityZone: ret.AvailabilityZone,
		VolumeName:       ret.Name,
		CloudDiskConf: &clouddisk.CloudDiskConf{
			VolumeType: ret.VolumeType,
			SizeInG:    int32(ret.Size), //类型不一致，使用强转
		},
		VolumeDesc:   ret.Description,
		VolumeStatus: ret.Status,
		CreatedTime:  ret.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		UpdatedTime:  ret.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
	}

	if len(ret.Attachments) != 0 {
		res.CloudDisk.AttachInstanceId = ret.Attachments[0].ServerID
		res.CloudDisk.AttachInstanceDevice = ret.Attachments[0].Device
		res.CloudDisk.AttachedTime = ret.Attachments[0].AttachedAt.Local().Format("2006-01-02 15:04:05")
	}
	log.Info("rpc CreateVolume: ", res, ". time elapse: ", timer_elasp.Elapse())
	return res, err
}

//云硬盘扩容
func (clouddisktask *CloudDiskService) ReqizeCloudDisk(ctx context.Context, req *clouddisk.ReqizeCloudDiskReq) (*clouddisk.CloudDiskRes, error) {
	timer_elasp := common.NewTimer()
	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = common.ENEWBLOCK.Code
		res.Msg = common.ENEWBLOCK.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	err = cinder_op.ExtendSize(client, req.VolumeId, cinder_op.ExtendSizeOpts{
		NewSize: int(req.CloudDiskConf.SizeInG), //类型不匹配，使用强转
	}).ExtractErr()

	if nil != err {
		res.Code = common.EEXTENDVOLUMESIZE.Code
		res.Msg = common.EEXTENDVOLUMESIZE.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.CloudDisk = &clouddisk.CloudDiskRes_CloudDisk{
		VolumeId: req.VolumeId,
		CloudDiskConf: &clouddisk.CloudDiskConf{
			VolumeType: req.CloudDiskConf.VolumeType,
			SizeInG:    int32(req.CloudDiskConf.SizeInG), //类型不一致，使用强转
		},
	}
	log.Info("rpc CreateVolume: ", res, ". time elapse: ", timer_elasp.Elapse())
	return res, err
}

//更新云硬盘信息
func (clouddisktask *CloudDiskService) ModifyCloudDiskInfo(ctx context.Context, req *clouddisk.ModifyCloudDiskInfoReq) (*clouddisk.CloudDiskRes, error) {
	timer_elasp := common.NewTimer()
	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = common.ENEWBLOCK.Code
		res.Msg = common.ENEWBLOCK.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	ret, err := cinder.Update(client, req.VolumeId, cinder.UpdateOpts{
		Name:        &req.VolumeName,
		Description: &req.VolumeDesc,
	}).Extract()

	if nil != err {
		res.Code = common.EVOLUMEUPDATE.Code
		res.Msg = common.EVOLUMEUPDATE.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.CloudDisk = &clouddisk.CloudDiskRes_CloudDisk{
		VolumeId:     ret.ID,
		VolumeName:   ret.Name,
		VolumeDesc:   ret.Description,
		UpdatedTime:  common.Now(),
		VolumeStatus: ret.Status,
	}
	log.Info("rpc CreateVolume: ", res, ". time elapse: ", timer_elasp.Elapse())
	return res, err
}

//云主机挂载、卸载
func (clouddisktask *CloudDiskService) OperateCloudDisk(ctx context.Context, req *clouddisk.OperateCloudDiskReq) (*clouddisk.CloudDiskRes, error) {
	timer_elasp := common.NewTimer()
	res := &clouddisk.CloudDiskRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if nil != err {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = common.ENEWCPU.Code
		res.Msg = common.ENEWCPU.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	cinder_client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})

	if nil != err {
		res.Code = common.ENEWBLOCK.Code
		res.Msg = common.ENEWBLOCK.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	if req.OpsType == "Attach" {
		ret, err := nova_op.Create(client, req.InstanceId, nova_op.CreateOpts{
			//Device:   req.,
			VolumeID: req.VolumeId,
		}).Extract()

		if err != nil {
			res.Code = common.EVOLUMEATTACH.Code
			res.Msg = common.EVOLUMEATTACH.Msg
			log.Error(res.Msg, ": ", err)
			return res, err
		}

		cinder_ret, err := cinder.Get(cinder_client, req.VolumeId).Extract()

		if nil != err {
			res.Code = common.ESHOWVOLUME.Code
			res.Msg = common.ESHOWVOLUME.Msg
			log.Error(res.Msg, ": ", err)
			return res, err
		}

		res.CloudDisk = &clouddisk.CloudDiskRes_CloudDisk{
			VolumeId:             ret.VolumeID,
			AttachInstanceId:     ret.ServerID,
			AttachInstanceDevice: ret.Device,
			AttachedTime:         common.Now(),
			UpdatedTime:          common.Now(),
			VolumeStatus:         cinder_ret.Status,
		}

	} else {
		err = nova_op.Delete(client, req.InstanceId, req.VolumeId).ExtractErr()

		if err != nil {
			res.Code = common.EVOLUMEDETACH.Code
			res.Msg = common.EVOLUMEDETACH.Msg
			log.Error(res.Msg, ": ", err)
			return res, err
		}

		cinder_ret, err := cinder.Get(cinder_client, req.VolumeId).Extract()

		if nil != err {
			res.Code = common.ESHOWVOLUME.Code
			res.Msg = common.ESHOWVOLUME.Msg
			log.Error(res.Msg, ": ", err)
			return res, err
		}

		res.CloudDisk = &clouddisk.CloudDiskRes_CloudDisk{
			VolumeId:         req.VolumeId,
			UpdatedTime:      common.Now(),
			AttachInstanceId: req.InstanceId,
			VolumeStatus:     cinder_ret.Status,
		}

	}
	log.Info("rpc CreateVolume: ", res, ". time elapse: ", timer_elasp.Elapse())
	return res, err
}
