package instancesvc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"errors"
	"iaas-api-server/common"
	"iaas-api-server/common/config"
	"iaas-api-server/proto/instance"

	log "github.com/sirupsen/logrus"
	gophercloud "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	servers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	flavors "github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	volumes "github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	startstop "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	suspendresume "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/suspendresume"
	"golang.org/x/net/context"
)

// InstanceService 计算实例服务
type InstanceService struct {
}

// CreateInstance 创建云主机
//   - 认证、获取 token
//   - 查询租户信息
//   - 创建数据盘
//   - 挂载系统盘、数据盘 (gosdk 暂不支持, 用原生 api 实现)
//   - 修改 hostname、root密码
func (is *InstanceService) CreateInstance(ctx context.Context, req *instance.CreateInstanceReq) (*instance.InstanceRes, error) {
	timer := common.NewTimer()
	log.Info("rpc CreateInstance req: ", req)
	res := &instance.InstanceRes{}

	token, err := common.AuthAndGetToken(req.Apikey, req.TenantId, req.PlatformUserid)
	if err != nil {
		log.Error("auth failed: ", err)
		res.Code = common.EUNAUTHORED.Code
		res.Msg = common.EUNAUTHORED.Msg
		return res, err
	}

	tenantInfo, err := common.QueryTenantInfoByTenantIdAndApikey(req.TenantId, req.Apikey)
	if err != nil {
		log.Error("query tenant info failed: ", err, ", tenant_id: ", req.TenantId, ", apikey: ", req.Apikey)
		res.Code = common.ENINSQUERYTENANT.Code
		res.Msg = common.ENINSQUERYTENANT.Msg
		return res, err
	}

	// 创建数据盘
	var volumeIDs []string
	for i := 0; i < len(req.DataDisks); i++ {
		disk := req.DataDisks[i]
		s, err := createVolume(disk.VolumeType, disk.SizeInG, tenantInfo.OpenstackProjectid,
			req.AvailabilityZone, token)
		if err != nil {
			res.Code = common.ENINSCREATEVOLUME.Code
			res.Msg = common.ENINSCREATEVOLUME.Msg
			return res, err
		}

		volumeIDs = append(volumeIDs, s)
	}

	// 构建脚本修改 hostname 与 root 密码
	script := "#!/bin/sh\n"
	if req.GuestOsHostname != "" {
		hostname := fmt.Sprintf("hostname '%s'\necho '%s' > /etc/hostname\n", req.GuestOsHostname, req.GuestOsHostname)
		script = script + hostname
	}
	if req.RootPass != "" {
		rootpass := "echo 'root:" + req.RootPass + "'|chpasswd\n"
		script = script + rootpass
	}
	if len(script) > 10  {
		fmt.Println("script: ", script)
		script = base64.StdEncoding.EncodeToString([]byte(script))
	} else {
		script = ""
	}

	js :=
`{
    "server": {
        "name": "{{.ServName}}",
        "flavorRef": "{{.FlavorRef}}",
        "availability_zone": "{{.AvailZone}}",
        "adminPass": "{{.RootPass}}",
        "networks": [{
            "uuid": "{{.Network}}"
        }],
        "block_device_mapping_v2": {{.Disks}},
        "security_groups": {{.SecurityGroup}},
        "user_data": "{{.Script}}"
    }
}`

	sysdisk :=
`{
    "uuid": "%s",
    "source_type": "image",
    "destination_type": "volume",
    "delete_on_termination": true,
    "boot_index": "0",
    "volume_size": %d
}`
	datadisk :=
`{
	"uuid": "%s",
	"source_type": "volume",
	"destination_type": "volume",
	"delete_on_termination": true,
	"boot_index": "-1",
	"volume_size": %d
}`

	disks := "[" + fmt.Sprintf(sysdisk, req.ImageRef, req.SystemDisk.SizeInG)
	for i := 0; i < len(volumeIDs); i++ {
		disks = disks + "," + fmt.Sprintf(datadisk, volumeIDs[i], req.DataDisks[i].SizeInG)
	}
	disks = disks + "]"

	secGroup := "["
	for i := 0; i < len(req.SecurityGroupName); i++ {
		if i > 0 {
			secGroup = secGroup + fmt.Sprintf(",{\"name\":\"%s\"}", req.SecurityGroupName[i])
		} else {
			secGroup = secGroup + fmt.Sprintf("{\"name\":\"%s\"}", req.SecurityGroupName[i])
		}
	}
	secGroup = secGroup + "]"

	mp := map[string]string {
		"ServName": req.InstanceName,
		"FlavorRef": req.FlavorRef,
		"AvailZone": req.AvailabilityZone,
		"RootPass": req.RootPass,
		"Network": req.NetworkUuid,
		"Disks": disks,
		"SecurityGroup": secGroup,
		"Script": script,
	}

	jsbody, _ := common.CreateJsonByTmpl(js, mp)
	log.Info("create instance json: ", string(jsbody))

	computeEndpoint, _ := config.GetString("compute_endpoint")
	body, err := common.CallRawAPI(computeEndpoint + "/servers", "POST", jsbody, token)
	if err != nil {
		res.Code = common.ENINSCREATE.Code
		res.Msg = common.ENINSCREATE.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	log.Info("create instance res: ", string(body))
	resmap := map[string]interface{}{}
	_ = json.Unmarshal(body, &resmap)
	servmap, ok := resmap["server"].(map[string]interface{})
	if ok {
		resmap = servmap
	}

	id, ok := resmap["id"].(string)
	if !ok || id == "" {
		log.Error("create instance failed, no id found in the response..")
		res.Code = common.ENINSCREATE.Code
		res.Msg = common.ENINSCREATE.Msg
		return res, err
	}

	/*
	serv, err := queryInstance(id, token)
	if err != nil {
		res.Code = common.ENINSQUERY.Code
		res.Msg = common.ENINSQUERY.Msg
		return res, err
	}
	 */

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	//flavor := serv["flavor"].(map[string]interface{})
	res.Instance = &instance.InstanceRes_Instance{
		InstanceId:          id,
		/*
		InstanceStatus:      serv["status"].(string),
		Region:              req.Region,
		AvailabilityZone:    req.AvailabilityZone,
		Flavor:              &instance.Flavor {
			FlavorId: flavor["id"].(string),
			//FlavorName: flavor["original_name"].(string),
			//FlavorVcpus: strconv.Itoa(flavor["vcpus"].(int)),
			//FlavorRam: strconv.Itoa(flavor["ram"].(int)),
			//FlavorDisk: strconv.Itoa(flavor["disk"].(int)),
		},
		ImageRef:            serv["image"].(string),//serv["image"].(map[string]interface{})["id"].(string),
		NetworkUuid:         req.NetworkUuid,
		InstanceName:        serv["name"].(string),
		GuestOsHostname:     req.GuestOsHostname,
		CreatedTime:         serv["created"].(string),
		UpdatedTime:         serv["updated"].(string),
		 */
	}

	// ip address
	/*
	addr := serv["addresses"].(map[string]interface{})
	for _, val := range addr {
		addrs := val.([]interface{})
		if len(addrs) > 0 {
			res.Instance.InstanceAddr = addrs[0].(map[string]interface{})["addr"].(string)
		}
		break
	}

	res.Instance.SystemDisk = req.SystemDisk
	res.Instance.DataDisks = req.DataDisks

	sec := serv["security_groups"].([]interface{})
	for i := 0; i < len(sec); i++ {
		res.Instance.SecurityGroupName = append(res.Instance.SecurityGroupName, sec[i].(map[string]interface{})["name"].(string))
	}
	 */

	log.Info("rpc CreateInstance res: ", res, ". time elapse: ", timer.Elapse())
	return res, nil
}

// GetInstance 获取云主机信息
func (is *InstanceService) GetInstance(ctx context.Context, req *instance.GetInstanceReq) (*instance.InstanceRes, error) {
	timer := common.NewTimer()
	log.Info("rpc GetInstance req: ", req)
	res := &instance.InstanceRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		res.Code = common.ENEWCPU.Code
		res.Msg = common.ENEWCPU.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	ss, err := servers.Get(client, req.InstanceId).Extract()
	if err != nil {
		res.Code = common.ENINSGET.Code
		res.Msg = common.ENINSGET.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	ssjson, err := json.Marshal(ss)
	log.Info("get instance res: ", string(ssjson))

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	res.Instance = &instance.InstanceRes_Instance{
		InstanceId:          ss.ID,
		InstanceStatus:      ss.Status,
		InstanceName:        ss.Name,
		CreatedTime:         ss.Created.Local().Format("2006-01-02 03:04:05"),
		UpdatedTime:         ss.Updated.Local().Format("2006-01-02 03:04:05"),
	}

	// ip address
	for _, val := range ss.Addresses {
		addrs, ok := val.([]interface{})
		if ok && len(addrs) > 0 {
			addr, ok := addrs[0].(map[string]interface{})["addr"].(string)
			if ok {
				res.Instance.InstanceAddr = addr
				break
			}
		}
	}

	// image id
	imageId, ok := ss.Image["id"].(string)
	if ok {
		res.Instance.ImageRef = imageId
	}

	// security group
	for i := 0; i < len(ss.SecurityGroups); i++ {
		secGroup, ok := ss.SecurityGroups[i]["name"].(string)
		if ok {
			res.Instance.SecurityGroupName = append(res.Instance.SecurityGroupName, secGroup)
		}
	}

	// volume
	volumeClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Error("openstack NewBlockStorageV3 failed: ", err)
	} else {
		for i := 0; i < len(ss.AttachedVolumes); i++ {
			volume, err := volumes.Get(volumeClient, ss.AttachedVolumes[i].ID).Extract()
			if err != nil {
				log.Error("query volume failed: ", err, ", id: ", ss.AttachedVolumes[i].ID)
				continue
			}

			// TODO: 暂时将 device name 为 /dev/vda 的卷作为系统卷，后续可能要根据实际情况调整
			if len(volume.Attachments) > 0 && volume.Attachments[0].Device == "/dev/vda" {
				res.Instance.SystemDisk = &instance.CloudDiskInfo{
					VolumeType:  volume.VolumeType,
					SizeInG:     int32(volume.Size),
					Device:      volume.Attachments[0].Device,
					VolumeId:    volume.ID,
				}
			} else {
				dataDisk := &instance.CloudDiskInfo{
					VolumeType:  volume.VolumeType,
					SizeInG:     int32(volume.Size),
					VolumeId:    volume.ID,
				}
				if len(volume.Attachments) > 0 {
					dataDisk.Device = volume.Attachments[0].Device
				}
				res.Instance.DataDisks = append(res.Instance.DataDisks, dataDisk)
			}
		}
	}

	// flavor
	flavorId, ok := ss.Flavor["id"].(string)
	if ok {
		x, err := flavors.Get(client, flavorId).Extract()
		if err != nil {
			log.Error("query flavor info failed: ", err, ", id: ", flavorId)
		} else {
			res.Instance.Flavor = &instance.Flavor{
				FlavorId:     x.ID,
				FlavorName:   x.Name,
				FlavorVcpus:  strconv.Itoa(x.VCPUs),
				FlavorRam:    strconv.Itoa(x.RAM),
				FlavorDisk:   strconv.Itoa(x.Disk),
			}
		}
	}

	log.Info("rpc GetInstance res: ", res, ", time elapse: ", timer.Elapse())
	return res, nil
}

// UpdateInstanceFlavor 修改云主机规格
//   1. 需要修改 nova-ncpu, nova-api 上的 nova.conf:
//      allow_resize_to_same_host=True
//   2. instance 为 ACTIVE 或 SHUTOFF 状态，才能执行 resize 操作
//   3. resize 后，要等 instance 变成 VERIFY_RESIZE 状态，才能执行 confirm_resize 操作
func (is *InstanceService) UpdateInstanceFlavor(ctx context.Context, req *instance.UpdateInstanceFlavorReq) (*instance.InstanceRes, error) {
	timer := common.NewTimer()
	log.Info("rpc UpdateInstanceFlavor req: ", req)
	res := &instance.InstanceRes{}

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		res.Code = common.ENEWCPU.Code
		res.Msg = common.ENEWCPU.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	// instance status 为 ACTIVE 或 SHUTOFF 时，才能执行 resize 操作
	ss, err := servers.Get(client, req.InstanceId).Extract()
	if err != nil || (ss.Status != "ACTIVE" && ss.Status != "SHUTOFF") {
		res.Code = common.ENINSSTATUS.Code
		res.Msg = common.ENINSSTATUS.Msg
		log.Error(res.Msg, ": ", ss.Status)
		return res, err
	}

	opts := servers.ResizeOpts{
		FlavorRef:        req.FlavorRef,
	}

	err = servers.Resize(client, req.InstanceId, opts).ExtractErr()
	if err != nil {
		res.Code = common.ENINSUPFLAVOR.Code
		res.Msg = common.ENINSUPFLAVOR.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	// status 为 VERIFY_RESIZE 才能执行 confirm_resize 操作
	resizeTimeout, err := config.GetInt("instance_resize_timeout")
	if err != nil {
		resizeTimeout = 20
	}
	resizeTimer := common.NewTimer()

	for ;; {
		ss, err := servers.Get(client, req.InstanceId).Extract()
		if err != nil {
			res.Code = common.ENINSGET.Code
			res.Msg = common.ENINSGET.Msg
			log.Error(res.Msg, ": ", err)
			return res, err
		}

		if ss.Status == "VERIFY_RESIZE" {
			break
		}

		if resizeTimer.Elapse().Seconds() > float64(resizeTimeout) {
			log.Warn("resize timeout: ", resizeTimeout)
			res.Code = common.ENINSUPFLAVOR.Code
			res.Msg = common.ENINSUPFLAVOR.Msg
			log.Error(res.Msg)
			return res, nil
		}

		time.Sleep(time.Duration(1) * time.Second)
	}

	err = servers.ConfirmResize(client, req.InstanceId).ExtractErr()
	if err != nil {
		res.Code = common.ENINSCONFIRMRESIZE.Code
		res.Msg = common.ENINSCONFIRMRESIZE.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	getReq := &instance.GetInstanceReq{
		Apikey:          req.Apikey,
		TenantId:        req.TenantId,
		PlatformUserid:  req.PlatformUserid,
		InstanceId:      req.InstanceId,
	}

	res, err = is.GetInstance(ctx, getReq)
	if err != nil {
		log.Warn("UpdateInstanceFlavor query instance failed: ", err)
		res.Code = common.EOK.Code
		res.Msg = common.EOK.Msg
		res.Instance = &instance.InstanceRes_Instance{
			InstanceId:  req.InstanceId,
		}
	} else {
		log.Info("rpc UpdateInstanceFlavor res: ", res, ", time elapse: ", timer.Elapse())
	}

	return res, nil
}

// DeleteInstance 删除云主机
func (is *InstanceService) DeleteInstance(ctx context.Context, req *instance.DeleteInstanceReq) (*instance.DeleteInstanceRes, error) {
	timer := common.NewTimer()
	log.Info("rpc DeleteInstance req: ", req)
	res := &instance.DeleteInstanceRes{}
	res.InstanceId = req.InstanceId;

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		res.Code = common.ENEWCPU.Code
		res.Msg = common.ENEWCPU.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	err = servers.Delete(client, req.InstanceId).ExtractErr()
	if err != nil {
		res.Code = common.ENINSDEL.Code
		res.Msg = common.ENINSDEL.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg
	res.DeletedTime = common.Now()

	log.Info("rpc DeleteInstance res: ", res, ", time elapse: ", timer.Elapse())
	return res, nil
}

// OperateInstance 启动-停止-挂起-重启云主机
func (is *InstanceService) OperateInstance(ctx context.Context, req *instance.OperateInstanceReq) (*instance.OperateInstanceRes, error) {
	timer := common.NewTimer()
	log.Info("rpc OperateInstance req: ", req)
	res := &instance.OperateInstanceRes{}
	res.InstanceId = req.InstanceId

	provider, err := common.GetOpenstackClient(req.Apikey, req.TenantId, req.PlatformUserid)
	if provider == nil {
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		res.Code = common.ENEWCPU.Code
		res.Msg = common.ENEWCPU.Msg
		log.Error(res.Msg, ": ", err)
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
		res.Code = common.ENINSOPUNKNOWN.Code
		res.Msg = common.ENINSOPUNKNOWN.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	if err != nil {
		res.Code = common.ENINSOP.Code
		res.Msg = common.ENINSOP.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg
	res.OperatedTime = common.Now()
	res.OperateType = req.OperateType

	log.Info("rpc OperateInstance res: ", res, ", time elapse: ", timer.Elapse())
	return res, nil
}

func queryInstance(id string, token string) (map[string]interface{}, error) {
	url, _ := config.GetString("compute_endpoint")
	url = url + "/servers/" + id
	timer := common.NewTimer()

	for ;; {
		res, err := common.CallRawAPI(url, "GET", []byte{}, token)
		if err != nil {
			log.Error("query instance failed: ", err, ", id: ", id)
			return nil, errors.New("query instance failed")
		}

		mp := map[string]interface{}{}
		err = json.Unmarshal(res, &mp)
		if err != nil {
			log.Error("json conv to map failed: ", err)
			return nil, err
		}

		server, ok := mp["server"].(map[string]interface{})
		if !ok {
			log.Error("query instance failed: ", string(res))
			return nil, errors.New("query instance error")
		}

		status, ok := server["status"].(string)
		//(server.Name != "" && len(server.Addresses) > 0 && len(server.Flavor) > 0 && server.Status != "")
		if ok && status != "" && len(server["addresses"].(map[string]interface{})) > 0 {
			return server, nil
		}

		if timer.Elapse().Seconds() > 30 {
			break
		}

		time.Sleep(time.Duration(3) * time.Second)
	}

	return nil, errors.New("query instance timedout")
}

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

	volumeEndpoint, err := config.GetString("volume_endpoint")
	if err != nil {
		log.Error("volume_endpoint not found in config file")
		return "", err
	}

	url := volumeEndpoint + "/" + projectID + "/volumes"
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
	vol, ok := rmp["volume"].(map[string]interface{})
	if !ok {
		return "", common.ENINSCREATEVOLUME
	}

	id, ok := vol["id"].(string)
	if !ok || id == "" {
		return "", common.ENINSCREATEVOLUME
	}

	// query volume status
	url = volumeEndpoint + "/" + projectID + "/volumes/" + id
	timer := common.NewTimer()
	for ;; {
		res, err = common.CallRawAPI(url, "GET", []byte{}, token)
		if err != nil {
			log.Error("query volume info error: ", err, ", volume id: ", id)
			return "", common.ENINSCREATEVOLUME
		}

		mp := make(map[string]interface{})
		err = json.Unmarshal(res, &mp)
		if err != nil {
			log.Error("json conv to map failed: ", err)
			return "", err
		}

		vmp, ok := mp["volume"].(map[string]interface{})
		if ok {
			mp = vmp
		}

		status, ok := mp["status"].(string)
		if ok && status == "available" {
			log.Info("query volume res: ", string(res))
			return id, nil
		}

		timeout, err := config.GetInt("create_volume_timeout")
		if err != nil {
			timeout = 20
		}

		if timer.Elapse().Seconds() > float64(timeout){
			return "", common.ENINSCREATEVOLUME
		}

		time.Sleep(time.Duration(1) * time.Second)
	}
}
