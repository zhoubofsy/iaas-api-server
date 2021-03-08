package vpcsvc

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
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
	"net"
)

type CreateVpcRPCTask struct {
	Req *vpc.CreateVpcReq
	Res *vpc.VpcRes
	Err *common.Error
}

func (rpctask *CreateVpcRPCTask) Run(context.Context) {
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

func (rpctask *CreateVpcRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed")
		return &common.Error{
			Code: common.ENEWNETWORK.Code,
			Msg:  err.Error(),
		}
	}

	networkOpts := networks.CreateOpts{
		Name:        rpctask.Req.GetVpcName(),
		Description: rpctask.Req.GetVpcDesc(),
	}

	networkInfo, err := networks.Create(client, networkOpts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("networks create failed")
		return &common.Error{
			Code: common.ENETWORKSCREATE.Code,
			Msg:  err.Error(),
		}
	}

	_, ipNet, err := net.ParseCIDR(rpctask.Req.Subnet)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("parse CIDR failed")
		return &common.Error{
			Code: common.EPARSECIDR.Code,
			Msg:  err.Error(),
		}
	}

	netWork, broadCast := cidr.AddressRange(ipNet)
	gwIp := cidr.Inc(netWork)
	startIp := cidr.Inc(gwIp)
	endIp := cidr.Dec(broadCast)
	gw := gwIp.String()
	var enableDHCP = true

	subnetOpts := subnets.CreateOpts{
		NetworkID: networkInfo.ID,
		CIDR:      rpctask.Req.GetSubnet(),
		IPVersion: 4,
		AllocationPools: []subnets.AllocationPool{
			{
				Start: startIp.String(),
				End:   endIp.String(),
			},
		},
		GatewayIP:  &gw,
		EnableDHCP: &enableDHCP,
	}

	subnetInfo, err := subnets.Create(client, subnetOpts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("subnet create failed")
		return &common.Error{
			Code: common.ESUBNETCREATE.Code,
			Msg:  err.Error(),
		}
	}
	createSubnetTime := getCurTime()

	routerName := "router-" + rpctask.Req.GetVpcName()
	routerInfo, err := routers.Create(client, routers.CreateOpts{Name: routerName}).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("router create failed")
		return &common.Error{
			Code: common.EROUTERCREATE.Code,
			Msg:  err.Error(),
		}
	}
	createRouterTime := getCurTime()

	interfaceOpts := routers.AddInterfaceOpts{
		SubnetID: subnetInfo.ID,
	}

	interfaceInfo, err := routers.AddInterface(client, routerInfo.ID, interfaceOpts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("router add interface failed")
		return &common.Error{
			Code: common.EINTERFACEADD.Code,
			Msg:  err.Error(),
		}
	}
	createInterfaceTime := getCurTime()

	portInfo, err := ports.Get(client, interfaceInfo.PortID).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("ports get failed")
		return &common.Error{
			Code: common.EPORTSGET.Code,
			Msg:  err.Error(),
		}
	}

	vpcRouter := vpc.VpcRes_Vpc_Router{
		RouterId:          routerInfo.ID,
		RouterName:        routerInfo.Name,
		RouterCreatedTime: createRouterTime,
		Intfs:             make([]*vpc.VpcRes_Vpc_Router_Intf, 0),
	}

	vpcInterface := vpc.VpcRes_Vpc_Router_Intf{
		IntfId:          portInfo.ID,
		IntfName:        portInfo.Name,
		IntfIp:          portInfo.FixedIPs[0].IPAddress,
		SubnetId:        portInfo.FixedIPs[0].SubnetID,
		IntfCreatedTime: createInterfaceTime,
	}
	vpcRouter.Intfs = append(vpcRouter.Intfs, &vpcInterface)

	vpcSubnet := vpc.VpcRes_Vpc_Subnet{
		Subnet:            subnetInfo.CIDR,
		SubnetId:          subnetInfo.ID,
		SubnetCreatedTime: createSubnetTime,
	}

	rpctask.Res.Vpc = &vpc.VpcRes_Vpc{
		VpcId:          networkInfo.ID,
		VpcName:        networkInfo.Name,
		VpcDesc:        networkInfo.Description,
		Region:         rpctask.Req.GetRegion(),
		Subnet:         make([]*vpc.VpcRes_Vpc_Subnet, 0),
		VpcStatus:      networkInfo.Status,
		VpcCreatedTime: networkInfo.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		Router:         &vpcRouter,
	}
	rpctask.Res.Vpc.Subnet = append(rpctask.Res.Vpc.Subnet, &vpcSubnet)

	return common.EOK
}

func (rpctask *CreateVpcRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetVpcName() ||
		"" == rpctask.Req.GetVpcDesc() ||
		"" == rpctask.Req.GetRegion() ||
		"" == rpctask.Req.GetSubnet() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *CreateVpcRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
