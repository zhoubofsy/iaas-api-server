/*================================================================
*
*  文件名称：update_firewall_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年04月02日
*
================================================================*/

package firewallsvc

import (
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	fg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	fr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/firewall"
)

type UpdateFirewallRPCTask struct {
	Req *firewall.UpdateFirewallReq
	Res *firewall.FirewallRes
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *UpdateFirewallRPCTask) Run(context.Context) {
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

func (rpctask *UpdateFirewallRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

	oldfw, err := fg.Get(client, rpctask.Req.FirewallId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":         err,
			"firewall id": rpctask.Req.FirewallId,
		}).Error("get firewall group failed")
	}

	// 1: 解绑
	if len(oldfw.Ports) > 0 {
		_, err := fg.Update(client, rpctask.Req.FirewallId, fg.UpdateOpts{
			Ports: []string{},
		}).Extract()
		if nil != err {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("update firewall info, detach firewall to port failed")
			return common.EFUPGROUP
		}
	}

	// 2： 创建新的policy，rule
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

	// 3: 绑定新的policy和重新绑定port
	ops := fg.UpdateOpts{
		IngressFirewallPolicyID: &rpctask.Res.Firewall.FirewallIngressPolicy.FirewallPolicyId,
		EgressFirewallPolicyID:  &rpctask.Res.Firewall.FirewallEgressPolicy.FirewallPolicyId,
	}
	if len(oldfw.Ports) > 0 {
		ops.Ports = oldfw.Ports
	}
	// 更新到新的防火墙规则
	_, err = fg.Update(client, rpctask.Req.FirewallId, ops).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("attach new firewall to port failed")

		// 新的policy绑定失败，恢复绑定到原来的
		if len(oldfw.Ports) > 0 {
			_, err = fg.Update(client, rpctask.Req.FirewallId, fg.UpdateOpts{
				Ports: oldfw.Ports,
			}).Extract()
			if nil != err {
				log.WithFields(log.Fields{
					"err":   err,
					"oldfw": oldfw,
				}).Error("new policy attach to port failed, restore old policy failed")
			}
		}
		return common.EFUPGROUP
	}
	rpctask.Res.Firewall.FirewallName = rpctask.Req.FirewallName
	rpctask.Res.Firewall.FirewallDesc = rpctask.Req.FirewallDesc

	// 4: 删除老的policy
	rpctask.deleteOldPolicy(client, oldfw)

	return common.EOK
}

func (rpctask *UpdateFirewallRPCTask) deleteOldPolicy(client *gophercloud.ServiceClient,
	oldfw *fg.Group) {
	fw := &firewall.Firewall{
		FirewallId: oldfw.ID,
		FirewallIngressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyId:    oldfw.IngressFirewallPolicyID,
			FirewallPolicyRules: make([]*firewall.FirewallRule, 0),
		},
		FirewallEgressPolicy: &firewall.FirewallPolicy{
			FirewallPolicyId:    oldfw.EgressFirewallPolicyID,
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
}

func (rpctask *UpdateFirewallRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() {
		return errors.New("input params is wrong")
	}

	for _, rule := range rpctask.Req.GetFirewallIngressPolicyRules() {
		if (fr.ProtocolAny == fr.Protocol(rule.FirewallRuleProtocol) ||
			fr.ProtocolTCP == fr.Protocol(rule.FirewallRuleProtocol) ||
			fr.ProtocolUDP == fr.Protocol(rule.FirewallRuleProtocol) ||
			fr.ProtocolICMP == fr.Protocol(rule.FirewallRuleProtocol)) &&
			(fr.ActionDeny == fr.Action(rule.FirewallRuleAction) ||
				fr.ActionAllow == fr.Action(rule.FirewallRuleAction) ||
				fr.ActionReject == fr.Action(rule.FirewallRuleAction)) {
			continue
		}
		sourceIP := net.ParseIP(rule.FirewallRuleSrcIp)
		sourcePort, srcerr := strconv.Atoi(rule.FirewallRuleSrcPort)
		dstIP := net.ParseIP(rule.FirewallRuleDstIp)
		dstPort, dsterr := strconv.Atoi(rule.FirewallRuleDstPort)

		if nil != sourceIP && nil != dstIP &&
			nil == srcerr && nil == dsterr &&
			sourcePort > 0 && dstPort > 0 {
			continue
		}

		return errors.New("params is unregular")
	}

	return nil
}

func (rpctask *UpdateFirewallRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
