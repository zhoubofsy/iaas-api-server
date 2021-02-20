/*================================================================
*
*  文件名称：get_nat_gateway_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年02月03日
*
================================================================*/
package natgatewaysvc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	routers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/natgateway"
)

// GetNatGatewayRPCTask create nat gateway
type GetNatGatewayRPCTask struct {
	Req *natgateway.GetNatGatewayReq
	Res *natgateway.NatGatewayRes
	Err *common.Error
}

// Run call this func
func (rpctask *GetNatGatewayRPCTask) Run(context.Context) {
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

func (rpctask *GetNatGatewayRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

	router, err := routers.Get(client, rpctask.Req.GetRouterId()).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get router failed.")
		return &common.Error{
			Code: common.ENGGET.Code,
			Msg:  err.Error(),
		}
	}

	rpctask.Res.NatGateway = &natgateway.NatGatewayRes_NatGateway{
		GatewayId:         rpctask.Req.GetGatewayId(),
		RouterId:          rpctask.Req.GetRouterId(),
		ExternalNetworkId: router.GatewayInfo.NetworkID,
		EnableSnat:        *(router.GatewayInfo.EnableSNAT),
		CreatedTime:       getCurTime(),
	}

	// TODO ExternalFixedIPs 可能有多个，通常只会有一个，但是不能排除多个的情况
	if len(router.GatewayInfo.ExternalFixedIPs) > 0 {
		rpctask.Res.NatGateway.ExternalFixedIp = router.GatewayInfo.ExternalFixedIPs[0].IPAddress
	}

	return common.EOK
}

func (rpctask *GetNatGatewayRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetRouterId() ||
		"" == rpctask.Req.GetGatewayId() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *GetNatGatewayRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
