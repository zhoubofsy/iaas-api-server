package flavorsvc

import (
	"iaas-api-server/common/config"
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
	timer := common.NewTimer()
	log.Info("rpc ListFlavors req: ", req)
	res := &flavor.ListFlavorsRes{}

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

	opts := flavors.ListOpts{
		Marker:       req.StartFlavorId,
		Limit:        int(req.PageSize),
		AccessType:   flavors.PublicAccess,
	}

	limit, err := config.GetInt("list_flavors_limit")
	if err != nil {
		limit = 1000
	}
	if opts.Limit > limit {
		opts.Limit = limit
	}

	pages, err := flavors.ListDetail(client, opts).AllPages()
	if err != nil {
		res.Code = common.ENFLVLIST.Code
		res.Msg = common.ENFLVLIST.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	allflavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		res.Code = common.ENFLVEXTRACT.Code
		res.Msg = common.ENFLVEXTRACT.Msg
		log.Error(res.Msg, ": ", err)
		return res, err
	}

	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg

	for i := 0; i < len(allflavors); i++ {
		x := allflavors[i]
		res.Flavors = append(res.Flavors, &flavor.Flavor{
			FlavorId:     x.ID,
			FlavorName:   x.Name,
			FlavorVcpus:  strconv.Itoa(x.VCPUs),
			FlavorRam:    strconv.Itoa(x.RAM),
			FlavorDisk:   strconv.Itoa(x.Disk),
		})

		if req.PageSize > 0 && len(res.Flavors) >= int(req.PageSize) {
			res.NextFlavorId = x.ID
			break
		}
	}

	log.Info("rpc ListFlavors res: ", res, ", time elapse: ", timer.Elapse())
	return res, nil
}

// GetFlavor 获取规格信息
func (*FlavorService) GetFlavor(ctx context.Context, req *flavor.GetFlavorReq) (*flavor.GetFlavorRes, error) {
	timer := common.NewTimer()
	log.Info("rpc GetFlavor req: ", req)
	res := &flavor.GetFlavorRes{}

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

	x, err := flavors.Get(client, req.FlavorId).Extract()
	if err != nil {
		res.Code = common.ENFLVGET.Code
		res.Msg = common.ENFLVGET.Msg
		log.Error(res.Msg, ": ", err)
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

	log.Info("rpc GetFlavor res: ", res, ", time elapse: ", timer.Elapse())
	return res, err
}
