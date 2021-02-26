package routesvc

import (
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/route"
)

type GetRouterRPCTask struct {
	Req *route.GetRouterReq
	Res *route.GetRouterRes
	Err *common.Error
}

func (rpctask *GetRouterRPCTask) Run(context.Context) {
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

func (rpctask *GetRouterRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
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

	router, err := routers.Get(client, rpctask.Req.RouterId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers get failed")
		return &common.Error{
			Code: common.EROUTERGET.Code,
			Msg:  err.Error(),
		}
	}

	pages, err := ports.List(client, ports.ListOpts{DeviceID: router.ID}).AllPages()
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

	portsInfo, err := ports.ExtractPorts(pages)
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

	rpctask.Res.Router = &route.GetRouterRes_Router{
		RouterId:          router.ID,
		RouterName:        router.Name,
		RouterCreatedTime: getCurTime(),
		Intfs:             make([]*route.GetRouterRes_Router_Intf, 0),
	}

	for _, v := range portsInfo {
		rpctask.Res.Router.Intfs = append(rpctask.Res.Router.Intfs, &route.GetRouterRes_Router_Intf{
			IntfId:          v.ID,
			IntfName:        v.Name,
			IntfIp:          v.FixedIPs[0].IPAddress,
			SubnetId:        v.FixedIPs[0].SubnetID,
			IntfCreatedTime: "",
		})
	}

	return common.EOK
}

func (rpctask *GetRouterRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetRouterId() ||
		"" == rpctask.Req.GetTenantId() {
		return errors.New("input param is wrong")
	}
	return nil
}

func (rpctask *GetRouterRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
