/*================================================================
*
*  文件名称：delete_firewall_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年04月02日
*
================================================================*/

package firewallsvc

import (
	"errors"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	fg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/firewall"
)

type DeleteFirewallRPCTask struct {
	Req *firewall.DeleteFirewallReq
	Res *firewall.DeleteFirewallRes
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *DeleteFirewallRPCTask) Run(context.Context) {
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

func (rpctask *DeleteFirewallRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
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

	// 确保没有端口跟当前安全组绑定才能删除
	if len(ret.Ports) > 0 {
		log.WithFields(log.Fields{
			"group": ret,
		}).Error("group has bind already, can not delete")
		return common.EFGROUPBIND
	}

	fw := &firewall.Firewall{
		FirewallId: ret.ID,
		FirewallIngressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyId:    ret.IngressFirewallPolicyID,
			FirewallPolicyRules: make([]*firewall.FirewallRule, 0),
		},
		FirewallEgressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyId:    ret.EgressFirewallPolicyID,
			FirewallPolicyRules: make([]*firewall.FirewallRule, 0),
		},
	}

	//获取ruleID
	var wg sync.WaitGroup
	{
		wg.Add(1)
		go getFirewallPolicy(client, fw.FirewallIngressPolicy, &wg)
	}
	{
		wg.Add(1)
		go getFirewallPolicy(client, fw.FirewallIngressPolicy, &wg)
	}
	wg.Wait()

	// 删除rule
	{
		for _, ruleID := range fw.FirewallIngressPolicy.FirewallPolicyRules {
			wg.Add(1)
			go delFirewallRule(client, ruleID.FirewallRuleId, &wg)
		}
	}
	{
		for _, ruleID := range fw.FirewallEgressPolicy.FirewallPolicyRules {
			wg.Add(1)
			go delFirewallRule(client, ruleID.FirewallRuleId, &wg)
		}
	}
	wg.Wait()

	// 删除policy
	{
		wg.Add(1)
		go delFirewallPolicy(client, fw.FirewallIngressPolicy.FirewallPolicyId, &wg)
	}
	{
		wg.Add(1)
		go delFirewallPolicy(client, fw.FirewallEgressPolicy.FirewallPolicyId, &wg)
	}
	wg.Wait()

	// 删除group
	err = fg.Delete(client, rpctask.Req.FirewallId).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err":        err,
			"firewallID": rpctask.Req.FirewallId,
		}).Error("delete firewall group id failed")
		return common.EFDELGROUP
	}

	rpctask.Res.DeletedTime = common.Now()
	rpctask.Res.FirewallId = rpctask.Req.FirewallId

	return common.EOK
}

func (rpctask *DeleteFirewallRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetFirewallId() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *DeleteFirewallRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
