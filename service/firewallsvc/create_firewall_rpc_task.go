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
	"net"
	"strconv"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
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
			go create_firewall_rules(client, rule, idx, &inRuleIDs, rpctask.Res.Firewall.FirewallIngressPolicy, &wg)
		}
	}
	{
		for idx, rule := range rpctask.Req.GetFirewallEgressPolicyRules() {
			wg.Add(1)
			go create_firewall_rules(client, rule, idx, &eRuleIDs, rpctask.Res.Firewall.FirewallEgressPolicy, &wg)
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
		log.WithFields(log.Fields{
			"res": rpctask.Res,
		}).Error("create firewall failed")
		return common.EFCREATE
	}

	return common.EOK
}

func (rpctask *CreateFirewallRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetFirewallName() ||
		"" == rpctask.Req.GetFirewallDesc() {
		return errors.New("input params is wrong")
	}

	// 对每个rule严格校验
	for _, rule := range rpctask.Req.GetFirewallIngressPolicyRules() {
		protocol := fr.Protocol(rule.FirewallRuleProtocol)
		action := fr.Action(rule.FirewallRuleAction)
		if (fr.ProtocolAny == protocol || fr.ProtocolTCP == protocol || fr.ProtocolUDP == protocol || fr.ProtocolICMP == protocol) &&
			(fr.ActionDeny == action || fr.ActionAllow == action || fr.ActionReject == action) {
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

func (rpctask *CreateFirewallRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
