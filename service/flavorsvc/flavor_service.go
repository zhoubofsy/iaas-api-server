package flavorsvc

import (
	"strconv"
	"iaas-api-server/common"
	"iaas-api-server/proto/flavor"

	log "github.com/sirupsen/logrus"
	gophercloud "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	flavors "github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"golang.org/x/net/context"
)

// FlavorService 计算实例服务
type FlavorService struct {
}

// ListFlavors 获取规格列表
func (*FlavorService) ListFlavors(ctx context.Context, req *flavor.ListFlavorsReq) (*flavor.ListFlavorsRes, error) {
	log.Info("rpc ListFlavors req: ", req)
	res := &flavor.ListFlavorsRes{}

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

	opts := flavors.ListOpts{
		//Marker:       req.PageNumber,
		//Limit:        int(req.PageSize),
		AccessType:   flavors.PublicAccess,
	}

	pages, err := flavors.ListDetail(client, opts).AllPages()
	if err != nil {
		res.Code = 20101
		res.Msg = "openstack list flavors failed"
		log.Error("openstack list flavors failed: ", err)
		return res, err
	}

	allflavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		res.Code = 20102
		res.Msg = "openstack extract flavors failed"
		log.Error("openstack extract flavors failed: ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	if req.PageNumber < 0 {
		log.Warn("recv negative pagenumber: ", req.PageNumber)
		req.PageNumber = 0
	}

	for i := int(req.PageNumber); i < len(allflavors); i++ {
		if req.PageSize > 0 && len(res.Flavors) >= int(req.PageSize) {
			break
		}

		x := allflavors[i]
		res.Flavors = append(res.Flavors, &flavor.Flavor{
			FlavorId:     x.ID,
			FlavorName:   x.Name,
			FlavorVcpus:  strconv.Itoa(x.VCPUs),
			FlavorRam:    strconv.Itoa(x.RAM),
			FlavorDisk:   strconv.Itoa(x.Disk),
		})
	}

	log.Info("rpc ListFlavors res: ", res)
	return res, nil
}

// GetFlavor 获取规格信息
func (*FlavorService) GetFlavor(ctx context.Context, req *flavor.GetFlavorReq) (*flavor.GetFlavorRes, error) {
	log.Info("rpc GetFlavor req: ", req)
	res := &flavor.GetFlavorRes{}

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

	x, err := flavors.Get(client, req.FlavorId).Extract()
	if err != nil {
		res.Code = 20103
		res.Msg = "openstack get flavor failed"
		log.Error("openstack get flavor failed: ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg
	res.Flavor = &flavor.Flavor{
		FlavorId:     x.ID,
		FlavorName:   x.Name,
		FlavorVcpus:  strconv.Itoa(x.VCPUs),
		FlavorRam:    strconv.Itoa(x.RAM),
		FlavorDisk:   strconv.Itoa(x.Disk),
	}

	log.Info("rpc GetFlavor res: ", res)
	return res, err
}
