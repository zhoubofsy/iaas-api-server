package vpcsvc

import (
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/vpc"
)

type SetVpcInfoRPCTask struct {
	Req *vpc.SetVpcInfoReq
	Res *vpc.VpcRes
	Err *common.Error
}

func (rpctask *SetVpcInfoRPCTask) Run(context.Context) {
	defer rpctask.setResult()

	if err := rpctask.checkParam(); nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("check param failed.")
		rpctask.Err = common.EPARAM
		return
	}

	providers, err := common.GetOpenstackClient(rpctask.Req.Apikey, rpctask.Req.TenantId, rpctask.Req.PlatformUserid)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("call common, get openstack client error")
		rpctask.Err = common.EGETOPSTACKCLIENT
		return
	}

	rpctask.Err = rpctask.execute(providers)
}

func (rpctask *SetVpcInfoRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return &common.Error{
			Code: common.ENEWNETWORK.Code,
			Msg:  err.Error(),
		}
	}

	//注意当前逻辑未处理VPC名称的更改，因VPC名称是获取路由器信息的参数，更改后不能获取路由器信息。
	//setVpcName := rpctask.Req.GetVpcName()
	setVpcDesc := rpctask.Req.GetVpcDesc()
	setOpts := networks.UpdateOpts{
		//Name:        &setVpcName,
		Description: &setVpcDesc,
	}

	networkInfo, err := networks.Update(client, rpctask.Req.VpcId, setOpts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("networks update failed")
		return &common.Error{
			Code: common.ENETWORKSSET.Code,
			Msg:  err.Error(),
		}
	}

	//获取子网信息
	vpcSubnets := make([]*vpc.VpcRes_Vpc_Subnet, 0)
	for _, v := range networkInfo.Subnets {
		subnetInfo, err := subnets.Get(client, v).Extract()
		if nil != err {
			log.WithFields(log.Fields{
				"err": err,
				"req": rpctask.Req.String(),
			}).Error("subnets get failed")
			return &common.Error{
				Code: common.ESUBNETGET.Code,
				Msg:  err.Error(),
			}
		}

		vpcSubnet := vpc.VpcRes_Vpc_Subnet{
			Subnet:            subnetInfo.CIDR,
			SubnetId:          subnetInfo.ID,
			SubnetCreatedTime: "",
		}
		vpcSubnets = append(vpcSubnets, &vpcSubnet)
	}

	//获取路由器信息
	routerName := "router-" + networkInfo.Name
	routerPages, err := routers.List(client, routers.ListOpts{Name: routerName}).AllPages()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers list failed")
		return &common.Error{
			Code: common.EROUTERLIST.Code,
			Msg:  err.Error(),
		}
	}

	routersInfo, err := routers.ExtractRouters(routerPages)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers extract failed")
		return &common.Error{
			Code: common.EROUTEREXTRACT.Code,
			Msg:  err.Error(),
		}
	}

	//根据名称查到的路由信息应该唯一
	if 1 != len(routersInfo) {
		log.WithFields(log.Fields{
			"err": "",
			"req": rpctask.Req.String(),
		}).Error("routersInfo is not unique")
		return &common.Error{
			Code: common.EROUTERINFO.Code,
			Msg:  "",
		}
	}

	portsPages, err := ports.List(client, ports.ListOpts{DeviceID: routersInfo[0].ID}).AllPages()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("ports list failed")
		return &common.Error{
			Code: common.EPORTSLIST.Code,
			Msg:  err.Error(),
		}
	}

	portsInfo, err := ports.ExtractPorts(portsPages)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("ports extract failed")
		return &common.Error{
			Code: common.EPORTSEXTRACT.Code,
			Msg:  err.Error(),
		}
	}

	vpcRouter := vpc.VpcRes_Vpc_Router{
		RouterId:          routersInfo[0].ID,
		RouterName:        routersInfo[0].Name,
		RouterCreatedTime: "",
		Intfs:             make([]*vpc.VpcRes_Vpc_Router_Intf, 0),
	}

	for _, v := range portsInfo {
		vpcRouter.Intfs = append(vpcRouter.Intfs, &vpc.VpcRes_Vpc_Router_Intf{
			IntfId:          v.ID,
			IntfName:        v.Name,
			IntfIp:          v.FixedIPs[0].IPAddress,
			SubnetId:        v.FixedIPs[0].SubnetID,
			IntfCreatedTime: "",
		})
	}

	rpctask.Res.Vpc = &vpc.VpcRes_Vpc{
		VpcId:          networkInfo.ID,
		VpcName:        networkInfo.Name,
		VpcDesc:        networkInfo.Description,
		Region:         "RegionOne",
		Subnet:         vpcSubnets,
		VpcStatus:      networkInfo.Status,
		VpcCreatedTime: networkInfo.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		Router:         &vpcRouter,
	}

	return common.EOK
}

func (rpctask *SetVpcInfoRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetVpcId() ||
		"" == rpctask.Req.GetVpcName() ||
		"" == rpctask.Req.GetVpcDesc() {
		return errors.New("input param is wrong")
	}
	return nil
}

func (rpctask *SetVpcInfoRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
