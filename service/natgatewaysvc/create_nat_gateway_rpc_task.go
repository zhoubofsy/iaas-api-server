/*================================================================
*
*  文件名称：create_nat_gateway_rpc_task.go
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

// CreateNatGatewayRPCTask create nat gateway
type CreateNatGatewayRPCTask struct {
	Req *natgateway.CreateNatGatewayReq
	Res *natgateway.NatGatewayRes
	Err *common.Error
}

// Run call this func
func (rpctask *CreateNatGatewayRPCTask) Run(context.Context) {
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

func (rpctask *CreateNatGatewayRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	enableSnat := true
	router, err := routers.Update(client, rpctask.Req.GetRouterId(), routers.UpdateOpts{
		GatewayInfo: &routers.GatewayInfo{
			NetworkID:  rpctask.Req.GetExternalNetworkId(),
			EnableSNAT: &enableSnat,
		},
	}).Extract()

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("create nat gateway failed.")
		return &common.Error{
			Code: common.ENGCREATE.Code,
			Msg:  err.Error(),
		}
	}

	rpctask.Res.NatGateway = &natgateway.NatGatewayRes_NatGateway{
		RouterId:          router.ID,
		ExternalNetworkId: router.GatewayInfo.NetworkID,
		EnableSnat:        *(router.GatewayInfo.EnableSNAT),
		CreatedTime:       getCurTime(),
	}

	if len(router.GatewayInfo.ExternalFixedIPs) > 0 {
		rpctask.Res.NatGateway.ExternalFixedIp = router.GatewayInfo.ExternalFixedIPs[0].IPAddress
		rpctask.Res.NatGateway.GatewayId = router.GatewayInfo.ExternalFixedIPs[0].SubnetID
	}

	return common.EOK
}

func (rpctask *CreateNatGatewayRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetRouterId() ||
		"" == rpctask.Req.GetExternalNetworkId() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *CreateNatGatewayRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
