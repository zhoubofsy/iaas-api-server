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

type groupInfo struct {
	name       string
	desc       string
	inPolicyID string
	ePolicyID  string
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
		FirewallIngressPolicy: &firewall.FirewallPolicy{},
		FirewallEgressPolicy:  &firewall.FirewallPolicy{},
	}

	var wg sync.WaitGroup
	// 异步创建rule
	inRuleIDs := make([]string, len(rpctask.Req.GetFirewallIngressPolicyRules()))
	eRuleIDs := make([]string, len(rpctask.Req.GetFirewallEgressPolicyRules()))
	{
		for idx, rule := range rpctask.Req.GetFirewallIngressPolicyRules() {
			wg.Add(1)
			go create_firewall_rules(client, rule, idx, inRuleIDs, &wg)
		}
	}
	{
		for idx, rule := range rpctask.Req.GetFirewallEgressPolicyRules() {
			wg.Add(1)
			go create_firewall_rules(client, rule, idx, eRuleIDs, &wg)
		}
	}
	wg.Wait()

	// 异步创建policy
	info := groupInfo{}
	{
		wg.Add(1)
		go create_firewall_policy(client, inRuleIDs, &info.inPolicyID, &wg)
	}
	{
		wg.Add(1)
		go create_firewall_policy(client, eRuleIDs, &info.ePolicyID, &wg)
	}
	wg.Wait()

	// 创建防火墙
	group := create_firewall_group(client, info)
	if nil == group {
		log.Error("call firewall, create group failed")
		return common.EFCREATE

	}

	return common.EOK
}

func create_firewall_group(client *gophercloud.ServiceClient,
	info groupInfo) *fg.Group {

	ret, err := fg.Create(client, fg.CreateOpts{
		Name:                    info.name,
		Description:             info.desc,
		IngressFirewallPolicyID: info.inPolicyID,
		EgressFirewallPolicyID:  info.ePolicyID,
	}).Extract()

	if nil != err {
		log.WithField("err", err).Error("call firewall, create group failed")
		return nil
	}

	return ret
}

func create_firewall_policy(client *gophercloud.ServiceClient,
	rules []string,
	policyID *string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	ret, err := fp.Create(client, fp.CreateOpts{
		FirewallRules: rules,
	}).Extract()

	if nil != err {
		log.WithField("err", err).Error("call firewall, create policy failed")
		return
	}

	*policyID = ret.ID
}

func create_firewall_rules(client *gophercloud.ServiceClient,
	rule *firewall.FirewallRuleSet,
	idx int,
	ruleIDs []string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	options := fr.CreateOpts{
		Protocol:             fr.Protocol(rule.FilewallRuleProtocol),
		Action:               fr.Action(rule.FilewallRuleAction),
		Description:          rule.FilewallRuleDesc,
		SourceIPAddress:      rule.FilewallRuleSrcIp,
		SourcePort:           rule.FilewallRuleSrcPort,
		DestinationIPAddress: rule.FilewallRuleDstIp,
		DestinationPort:      rule.FilewallRuleDstPort,
	}

	ret, err := fr.Create(client, options).Extract()
	if nil != err {
		log.WithField("err", err).Error("call firewall, create rules failed ")
		return
	}

	ruleIDs[idx] = ret.ID
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
