package instancesvc

import (
	"encoding/json"
	"strconv"
	"time"
	"iaas-api-server/common"
	"iaas-api-server/common/config"
	"iaas-api-server/proto/instance"

	log "github.com/sirupsen/logrus"
	gophercloud "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	servers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	startstop "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	suspendresume "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/suspendresume"
	"golang.org/x/net/context"
)

// InstanceService 计算实例服务
type InstanceService struct {
}

/*
func (is *InstanceService) CreateInstance(ctx context.Context, req *instance.CreateInstanceReq) (*instance.InstanceRes, error) {
	log.Info("rpc CreateInstance req: ", req)
	res := &instance.InstanceRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
        Region: req.Region,
	})

	// TODO: 后续在 common/error.go 中定义错误码
	if err != nil {
		res.Code = 20000
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	opts := servers.CreateOpts{
		Name:             req.InstanceName,
		ImageRef:         req.ImageRef,
		FlavorRef:        req.FlavorRef,
		SecurityGroups:   req.SecurityGroupName,
		AvailabilityZone: req.AvailabilityZone,
		Networks: []servers.Network{
			servers.Network{UUID: req.NetworkUuid},
		},
		AdminPass: req.RootPass,
	}

	// TODO: 后续在 common/error.go 中定义错误码
	ss, err := servers.Create(client, opts).Extract()
	if err != nil {
		res.Code = 20001
		res.Msg = "openstack create instance failed"
		log.Error("openstack create instance failed: ", err)
		return res, err
	}

	ssjson, err := json.Marshal(ss)
	log.Info("create instance res: ", string(ssjson))

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	res.Instance = &instance.InstanceRes_Instance{
		InstanceId:          ss.ID,
		InstanceStatus:      ss.Status,
		InstanceAddr:        ss.AccessIPv4,  // TODO: 填充 AccessIPv4?
		Region:              req.Region,
		AvailabilityZone:    req.AvailabilityZone,
		Flavor:              &instance.Flavor {  // TODO: 后续根据 FlavorMap 填充
			FlavorId: "xxx",
			FlavorName: "xxx",
			FlavorVcpus: "xxx",
			FlavorRam: "xxx",
			FlavorDisk: "xxx",
		}, 
		ImageRef:            req.ImageRef,
		// TODO: SystemDisk:
		NetworkUuid:         req.NetworkUuid,
		// TODO: SecurityGroupName:
		InstanceName:        ss.Name,
		//GuestOsHostname:  ss.HostID,  // TODO: HostID?
		CreatedTime:         ss.Created.Local().Format("2006-01-02 03:04:05"),
		UpdatedTime:         ss.Updated.Local().Format("2006-01-02 03:04:05"),
	}

	log.Info("rpc CreateInstance res: ", res)
	return res, nil
}
*/

// 创建 volume, 返回 volume id
//   volume endpoint: http://192.168.66.131/volume/v3
func createVolume(volumeType string, sizeInG int32, projectID string, availZone string, token string) (string, error) {
	jstmpl :=
`{
    "volume": {
        "size": {{.VolumeSize}},
        "multiattach": false,
        "volume_type": "{{.VolumeType}}",
        "availability_zone": "{{.AvailZone}}"
    }
}`

	mp := map[string]string {
		"VolumeSize": strconv.Itoa(int(sizeInG)),
		"VolumeType": volumeType,
		"AvailZone": availZone,
	}

	jsbody, _ := common.CreateJsonByTmpl(jstmpl, mp)

	url, err := config.GetString("volume_endpoint")
	if err != nil {
		log.Error("volume endpoint not found")
		return "", err
	}

	url = url + "/" + projectID + "/volumes"

	res, err := common.CallRawAPI(url, "POST", jsbody, token)
	if err != nil {
		log.Error("create volume failed: ", err)
		return "", err
	}

	rmp := make(map[string]interface{})
	err = json.Unmarshal(res, &rmp)
	if err != nil {
		log.Error("json conv to map failed: ", err)
		return "", err
	}

	log.Info("create volume res: ", string(res))
	vol, _ := rmp["volume"]
	id, ok := vol.(map[string]interface{})["id"]
	if ok == true && id != "" {
		return id.(string), nil
	}

	return "", nil
}

// CreateInstance 创建云主机
func (is *InstanceService) CreateInstance(ctx context.Context, req *instance.CreateInstanceReq) (*instance.InstanceRes, error) {
	log.Info("rpc CreateInstance req: ", req)
	res := &instance.InstanceRes{}

	token, err := common.AuthAndGetToken(req.Apikey, req.TenantId, req.PlatformUserid)
	if err != nil {
		log.Error("get token failed: ", err)
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		return res, err
	}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: req.Region,
	})

	// TODO: 后续在 common/error.go 中定义错误码
	if err != nil {
		res.Code = 20000
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	opts := servers.CreateOpts{
		Name:             req.InstanceName,
		ImageRef:         req.ImageRef,
		FlavorRef:        req.FlavorRef,
		SecurityGroups:   req.SecurityGroupName,
		AvailabilityZone: req.AvailabilityZone,
		Networks: []servers.Network{
			servers.Network{UUID: req.NetworkUuid},
		},
		AdminPass: req.RootPass,
	}

	// TODO: 后续在 common/error.go 中定义错误码
	ss, err := servers.Create(client, opts).Extract()
	if err != nil {
		res.Code = 20001
		res.Msg = "openstack create instance failed"
		log.Error("openstack create instance failed: ", err)
		return res, err
	}

	ssjson, err := json.Marshal(ss)
	log.Info("create instance res: ", string(ssjson))

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	res.Instance = &instance.InstanceRes_Instance{
		InstanceId:          ss.ID,
		InstanceStatus:      ss.Status,
		InstanceAddr:        ss.AccessIPv4,  // TODO: 填充 AccessIPv4?
		Region:              req.Region,
		AvailabilityZone:    req.AvailabilityZone,
		Flavor:              &instance.Flavor {  // TODO: 后续根据 FlavorMap 填充
			FlavorId: "xxx",
			FlavorName: "xxx",
			FlavorVcpus: "xxx",
			FlavorRam: "xxx",
			FlavorDisk: "xxx",
		},
		ImageRef:            req.ImageRef,
		// TODO: SystemDisk:
		NetworkUuid:         req.NetworkUuid,
		// TODO: SecurityGroupName:
		InstanceName:        ss.Name,
		//GuestOsHostname:  ss.HostID,  // TODO: HostID?
		CreatedTime:         ss.Created.Local().Format("2006-01-02 03:04:05"),
		UpdatedTime:         ss.Updated.Local().Format("2006-01-02 03:04:05"),
	}

	log.Info("rpc CreateInstance res: ", res)
	return res, nil
}

// GetInstance 获取云主机信息
func (is *InstanceService) GetInstance(ctx context.Context, req *instance.GetInstanceReq) (*instance.InstanceRes, error) {
	log.Info("rpc GetInstance req: ", req)
	res := &instance.InstanceRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	// TODO: 后续在 common/error.go 中定义错误码
	if err != nil {
		res.Code = 20000
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	ss, err := servers.Get(client, req.InstanceId).Extract()
	if err != nil {
		res.Code = 20002
		res.Msg = "openstack get instance failed"
		log.Error("openstack get instance failed: ", err)
		return res, err
	}

	ssjson, err := json.Marshal(ss)
	log.Info("get instance res: ", string(ssjson))

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	res.Instance = &instance.InstanceRes_Instance{
		InstanceId:          ss.ID,
		InstanceStatus:      ss.Status,
		InstanceAddr:        ss.AccessIPv4,  // TODO: 填充 AccessIPv4?
		//Region:              req.Region,
		//AvailabilityZone:    req.AvailabilityZone,
		Flavor:              &instance.Flavor {  // TODO: 后续根据 FlavorMap 填充
			FlavorId: "xxx",
			FlavorName: "xxx",
			FlavorVcpus: "xxx",
			FlavorRam: "xxx",
			FlavorDisk: "xxx",
		},
		//ImageRef:            req.ImageRef,
		// TODO: SystemDisk:  ss.AttachedVolumes
		//NetworkUuid:         req.NetworkUuid,
		// TODO: SecurityGroupName: ss.SecurityGroups
		InstanceName:        ss.Name,
		GuestOsHostname:     ss.HostID,
		CreatedTime:         ss.Created.Local().Format("2006-01-02 03:04:05"),
		UpdatedTime:         ss.Updated.Local().Format("2006-01-02 03:04:05"),
	}

	log.Info("rpc GetInstance res: ", res)
	return res, nil
}

// UpdateInstanceFlavor 修改云主机规格
func (is *InstanceService) UpdateInstanceFlavor(ctx context.Context, req *instance.UpdateInstanceFlavorReq) (*instance.InstanceRes, error) {
	log.Info("rpc UpdateInstanceFlavor req: ", req)
	res := &instance.InstanceRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	// TODO: 后续在 common/error.go 中定义错误码
	if err != nil {
		res.Code = 20000
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	opts := servers.ResizeOpts{
		FlavorRef:        req.FlavorRef,
	}

	// TODO: 后续在 common/error.go 中定义错误码
	err = servers.Resize(client, req.InstanceId, opts).ExtractErr()
	if err != nil {
		res.Code = 20002
		res.Msg = "openstack update instance flavor failed"
		log.Error("openstack update instance flavor failed: ", err)
		return res, err
	}

	err = servers.ConfirmResize(client, req.InstanceId).ExtractErr()
	if err != nil {
		res.Code = 20003
		res.Msg = "openstack confirm resize instance failed"
		log.Error("openstack confirm resize instance failed: ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	res.Instance = &instance.InstanceRes_Instance{
		InstanceId:          req.InstanceId,
	}

	log.Info("rpc UpdateInstanceFlavor res: ", res)
	return res, nil
}

// DeleteInstance 删除云主机
func (is *InstanceService) DeleteInstance(ctx context.Context, req *instance.DeleteInstanceReq) (*instance.DeleteInstanceRes, error) {
	log.Info("rpc DeleteInstance req: ", req)
	res := &instance.DeleteInstanceRes{}
	res.InstanceId = req.InstanceId;

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	// TODO: 后续在 common/error.go 中定义错误码
	if err != nil {
		res.Code = 20000
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	err = servers.Delete(client, req.InstanceId).ExtractErr()
	if err != nil {
		res.Code = 20004
		res.Msg = "openstack delete instance failed"
		log.Error("openstack delete instance failed: ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg
	res.DeletedTime = time.Now().Format("2006-01-02 03:04:05")

	log.Info("rpc DeleteInstance res: ", res)
	return res, nil
}

// OperateInstance 启动-停止-挂起-重启云主机
func (is *InstanceService) OperateInstance(ctx context.Context, req *instance.OperateInstanceReq) (*instance.OperateInstanceRes, error) {
	log.Info("rpc OperateInstance req: ", req)
	res := &instance.OperateInstanceRes{}
	res.InstanceId = req.InstanceId

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		log.Error("GetOpenstackClient failed: ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	// TODO: 后续在 common/error.go 中定义错误码
	if err != nil {
		res.Code = 20000
		res.Msg = "openstack NewComputeV2 failed"
		log.Error("openstack NewComputeV2 failed: ", err)
		return res, err
	}

	if req.OperateType == "start" {
		err = startstop.Start(client, req.InstanceId).ExtractErr()
	} else if req.OperateType == "stop" {
		err = startstop.Stop(client, req.InstanceId).ExtractErr()
	} else if req.OperateType == "suspend" {
		err = suspendresume.Suspend(client, req.InstanceId).ExtractErr()
	} else if req.OperateType == "resume" {
		err = suspendresume.Resume(client, req.InstanceId).ExtractErr()
	} else if req.OperateType == "softreboot" {
		opts := servers.RebootOpts{
			Type: servers.SoftReboot,
		}
		err = servers.Reboot(client, req.InstanceId, opts).ExtractErr()
	} else if req.OperateType == "hardreboot"{
		opts := servers.RebootOpts{
			Type: servers.HardReboot,
		}
		err = servers.Reboot(client, req.InstanceId, opts).ExtractErr()
	} else {
		res.Code = 20006
		res.Msg = "unsupported operate type"
		log.Error("unsupported operate type: ", req.OperateType)
		return res, err
	}

	if err != nil {
		res.Code = 20005
		res.Msg = "openstack operate instance failed"
		log.Error("openstack operate instance failed: ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg
	res.OperatedTime = time.Now().Format("2006-01-02 03:04:05")
	res.OperateType = req.OperateType

	log.Info("rpc OperateInstance res: ", res)
	return res, nil
}
