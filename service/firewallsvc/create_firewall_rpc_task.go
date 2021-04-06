/*================================================================
*
*  文件名称：create_firewall_rpc_task.go
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
	fp "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
	fr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/firewall"
)

type CreateFirewallRPCTask struct {
	Req *firewall.CreateFirewallReq
	Res *firewall.FirewallRes
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *CreateFirewallRPCTask) Run(context.Context) {
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

func (rpctask *CreateFirewallRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

	rpctask.Res.Firewall = &firewall.Firewall{
		FirewallIngressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyRules: make([]*firewall.FirewallRule, len(rpctask.Req.GetFirewallIngressPolicyRules())),
		},
		FirewallEgressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyRules: make([]*firewall.FirewallRule, len(rpctask.Req.GetFirewallEgressPolicyRules())),
		},
	}

	var wg sync.WaitGroup
	// 异步创建rule
	inRuleIDs := make([]string, len(rpctask.Req.GetFirewallIngressPolicyRules()))
	eRuleIDs := make([]string, len(rpctask.Req.GetFirewallEgressPolicyRules()))
	{
		for idx, rule := range rpctask.Req.GetFirewallIngressPolicyRules() {
			wg.Add(1)
			go create_firewall_rules(client, rule, idx, inRuleIDs, rpctask.Res.Firewall.FirewallIngressPolicy.FirewallPolicyRules, &wg)
		}
	}
	{
		for idx, rule := range rpctask.Req.GetFirewallEgressPolicyRules() {
			wg.Add(1)
			go create_firewall_rules(client, rule, idx, eRuleIDs, rpctask.Res.Firewall.FirewallEgressPolicy.FirewallPolicyRules, &wg)
		}
	}
	wg.Wait()

	// 异步创建policy
	{
		wg.Add(1)
		go create_firewall_policy(client, inRuleIDs, rpctask.Res.Firewall.FirewallIngressPolicy, &wg)
	}
	{
		wg.Add(1)
		go create_firewall_policy(client, eRuleIDs, rpctask.Res.Firewall.FirewallEgressPolicy, &wg)
	}
	wg.Wait()

	// 创建防火墙
	rpctask.Res.Firewall.FirewallName = rpctask.Req.FirewallName
	rpctask.Res.Firewall.FirewallDesc = rpctask.Req.FirewallDesc
	create_firewall_group(client, rpctask.Res.Firewall)

	if rpctask.Res.Firewall.UpdatedTime == "" {
		return common.EFCREATE
	}

	return common.EOK
}

func create_firewall_group(client *gophercloud.ServiceClient,
	group *firewall.Firewall) {

	ret, err := fg.Create(client, fg.CreateOpts{
		Name:                    group.FirewallName,
		Description:             group.FirewallDesc,
		IngressFirewallPolicyID: group.FirewallIngressPolicy.FirewallPolicyId,
		EgressFirewallPolicyID:  group.FirewallEgressPolicy.FirewallPolicyId,
	}).Extract()

	if nil != err {
		log.WithField("err", err).Error("call firewall, create group failed")
		return
	}

	group.FirewallId = ret.ID
	group.FirewallAttachedPortId = ""
	group.FirewallStatus = ret.Status
	group.CreatedTime = common.Now()
	group.UpdatedTime = common.Now()
}

func create_firewall_policy(client *gophercloud.ServiceClient,
	rules []string,
	policy *firewall.FirewallPolicy,
	wg *sync.WaitGroup) {
	defer wg.Done()

	ret, err := fp.Create(client, fp.CreateOpts{
		FirewallRules: rules,
	}).Extract()

	if nil != err {
		log.WithField("err", err).Error("call firewall, create policy failed")
		return
	}

	*policy = firewall.FirewallPolicy{
		FirewallPolicyId:   ret.ID,
		FirewallPolicyName: ret.Name,
		FirewallPolicyDesc: ret.Description,
	}
}

func create_firewall_rules(client *gophercloud.ServiceClient,
	rule *firewall.FirewallRuleSet,
	idx int,
	ruleIDs []string,
	rules []*firewall.FirewallRule,
	wg *sync.WaitGroup) {
	defer wg.Done()

	options := fr.CreateOpts{
		Protocol:             fr.Protocol(rule.FirewallRuleProtocol),
		Action:               fr.Action(rule.FirewallRuleAction),
		Description:          rule.FirewallRuleDesc,
		SourceIPAddress:      rule.FirewallRuleSrcIp,
		SourcePort:           rule.FirewallRuleSrcPort,
		DestinationIPAddress: rule.FirewallRuleDstIp,
		DestinationPort:      rule.FirewallRuleDstPort,
	}

	ret, err := fr.Create(client, options).Extract()
	if nil != err {
		log.WithField("err", err).Error("call firewall, create rules failed ")
		return
	}

	ruleIDs[idx] = ret.ID
	rules[idx] = &firewall.FirewallRule{
		FirewallRuleId:       ret.ID,
		FirewallRuleName:     ret.Name,
		FirewallRuleDesc:     ret.Description,
		FirewallRuleAction:   ret.Action,
		FirewallRuleProtocol: ret.Protocol,
		FirewallRuleSrcIp:    ret.SourceIPAddress,
		FirewallRuleSrcPort:  ret.SourcePort,
		FirewallRuleDstIp:    ret.DestinationIPAddress,
		FirewallRuleDstPort:  ret.DestinationPort,
	}
}

func (rpctask *CreateFirewallRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *CreateFirewallRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
