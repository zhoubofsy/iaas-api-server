/*================================================================
*
*  文件名称：get_firewall_rpc_task.go
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

type GetFirewallRPCTask struct {
	Req *firewall.GetFirewallReq
	Res *firewall.FirewallRes
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *GetFirewallRPCTask) Run(context.Context) {
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

func (rpctask *GetFirewallRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
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
		log.WithField("err", err).Error("get firewall group failed")
		return common.EFGETGROUP
	}

	rpctask.Res.Firewall = &firewall.Firewall{
		FirewallId:             ret.ID,
		FirewallName:           ret.Name,
		FirewallDesc:           ret.Description,
		FirewallAttachedPortId: "",
		FirewallStatus:         ret.Status,
		FirewallIngressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyId:    ret.IngressFirewallPolicyID,
			FirewallPolicyRules: make([]*firewall.FirewallRule, 0),
		},
		FirewallEgressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyId:    ret.EgressFirewallPolicyID,
			FirewallPolicyRules: make([]*firewall.FirewallRule, 0),
		},
		CreatedTime: "",
		UpdatedTime: "",
	}
	if len(ret.Ports) > 0 {
		rpctask.Res.Firewall.FirewallAttachedPortId = ret.Ports[0]
	}

	var wg sync.WaitGroup
	// 获取policy
	{
		wg.Add(1)
		go getFirewallPolicy(client, rpctask.Res.Firewall.FirewallIngressPolicy, &wg)
	}
	{
		wg.Add(1)
		go getFirewallPolicy(client, rpctask.Res.Firewall.FirewallEgressPolicy, &wg)
	}
	wg.Wait()

	// 获取rule
	{
		for _, rule := range rpctask.Res.Firewall.FirewallIngressPolicy.FirewallPolicyRules {
			wg.Add(1)
			go getFirewallRule(client, rule, &wg)
		}
	}
	{
		for _, rule := range rpctask.Res.Firewall.FirewallEgressPolicy.FirewallPolicyRules {
			wg.Add(1)
			go getFirewallRule(client, rule, &wg)
		}
	}

	wg.Wait()

	return common.EOK
}

func (rpctask *GetFirewallRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetFirewallId() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *GetFirewallRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
