/*================================================================
*
*  文件名称：operate_firewall_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年04月02日
*
================================================================*/

package firewallsvc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	fg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/firewall"
)

type OperateFirewallRPCTask struct {
	Req *firewall.OperateFirewallReq
	Res *firewall.OperateFirewallRes
	Err *common.Error
}

var (
	DetachType string = "detach"
	AttachType string = "attach"
)

// Run call this func for doing task
func (rpctask *OperateFirewallRPCTask) Run(context.Context) {
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

func (rpctask *OperateFirewallRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

	ret, err := fg.Get(client, rpctask.Req.FirewallId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":         err,
			"firewall id": rpctask.Req.FirewallId,
		}).Error("get firewall group failed")
		return common.EFGETGROUP
	}

	if AttachType == rpctask.Req.OpsType {
		if len(ret.Ports) > 0 {
			log.WithFields(log.Fields{
				"group": ret,
			}).Error("group has bind already, can bot bind again")
			return common.EFGROUPBIND
		}
		_, err := fg.Update(client, rpctask.Req.FirewallId, fg.UpdateOpts{
			Ports: []string{rpctask.Req.PortId},
		}).Extract()
		if nil != err {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("attach firewall to port failed")
			return common.EFUPGROUP
		}
	} else if DetachType == rpctask.Req.OpsType {
		if len(ret.Ports) > 0 {
			err := DetachFirewallToPortsRawAPI(client, rpctask.Req.FirewallId)
			if nil != err {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("detach firewall to port failed")
				return common.EFUPGROUP
			}
		}
	}

	rpctask.Res.FirewallId = rpctask.Req.FirewallId
	rpctask.Res.FirewallAttachedPortId = rpctask.Req.PortId
	rpctask.Res.OpsType = rpctask.Req.OpsType
	rpctask.Res.OperatedTime = common.Now()

	return common.EOK
}

func (rpctask *OperateFirewallRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		(DetachType != rpctask.Req.GetOpsType() && AttachType != rpctask.Req.GetOpsType()) ||
		"" == rpctask.Req.GetFirewallId() ||
		"" == rpctask.Req.GetPortId() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *OperateFirewallRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
