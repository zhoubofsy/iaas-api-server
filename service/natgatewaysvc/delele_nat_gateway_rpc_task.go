/*================================================================
*
*  文件名称：delele_nat_gateway_rpc_task.go
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

// DeleteNatGatewayRPCTask create nat gateway
type DeleteNatGatewayRPCTask struct {
	Req *natgateway.DeleteNatGatewayReq
	Res *natgateway.DeleteNatGatewayRes
	Err *common.Error
}

// Run call this func
func (rpctask *DeleteNatGatewayRPCTask) Run(context.Context) {
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

func (rpctask *DeleteNatGatewayRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{
		Name: "neutron",
	})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	//TODO 此处路由外部网关需调用路由的update接口，不能使用removeInterface接口
	_, err = routers.Update(client, rpctask.Req.GetRouterId(), &routers.UpdateOpts{
		GatewayInfo: &routers.GatewayInfo{},
	}).Extract()

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("delete nat gateway failed.")
		return &common.Error{
			Code: common.ENGDELETE.Code,
			Msg:  err.Error(),
		}
	}

	return common.EOK
}

func (rpctask *DeleteNatGatewayRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetRouterId() ||
		"" == rpctask.Req.GetGatewayId() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *DeleteNatGatewayRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
	rpctask.Res.RouterId = rpctask.Req.GetRouterId()
	rpctask.Res.GatewayId = rpctask.Req.GetGatewayId()
	rpctask.Res.DeletedTime = getCurTime()
}
