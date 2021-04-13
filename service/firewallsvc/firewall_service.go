/*================================================================
*
*  文件名称：firewall_service.go
*  创 建 者: mongia
*  创建日期：2021年04月02日
*
================================================================*/

package firewallsvc

import (
	"sync"

	"golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"github.com/gophercloud/gophercloud"
	fg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	fp "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
	fr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	log "github.com/sirupsen/logrus"

	"iaas-api-server/common"
	"iaas-api-server/proto/firewall"
)

type FirewallService struct {
	firewall.UnimplementedFirewallServiceServer
}

func (pthis *FirewallService) GetFirewall(ctx context.Context, req *firewall.GetFirewallReq) (*firewall.FirewallRes, error) {
	task := GetFirewallRPCTask{
		Req: req,
		Res: &firewall.FirewallRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) CreateFirewall(ctx context.Context, req *firewall.CreateFirewallReq) (*firewall.FirewallRes, error) {
	task := CreateFirewallRPCTask{
		Req: req,
		Res: &firewall.FirewallRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) UpdateFirewall(ctx context.Context, req *firewall.UpdateFirewallReq) (*firewall.FirewallRes, error) {
	task := UpdateFirewallRPCTask{
		Req: req,
		Res: &firewall.FirewallRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) DeleteFirewall(ctx context.Context, req *firewall.DeleteFirewallReq) (*firewall.DeleteFirewallRes, error) {
	task := DeleteFirewallRPCTask{
		Req: req,
		Res: &firewall.DeleteFirewallRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) OperateFirewall(ctx context.Context, req *firewall.OperateFirewallReq) (*firewall.OperateFirewallRes, error) {
	task := OperateFirewallRPCTask{
		Req: req,
		Res: &firewall.OperateFirewallRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func getFirewallRule(client *gophercloud.ServiceClient,
	rule *firewall.FirewallRule,
	wg *sync.WaitGroup) {
	defer wg.Done()

	ret, err := fr.Get(client, rule.FirewallRuleId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":  err,
			"rule": rule,
		}).Error("query firewall rule failed")
		return
	}

	rule.FirewallRuleName = ret.Name
	rule.FirewallRuleDesc = ret.Description
	rule.FirewallRuleAction = ret.Action
	rule.FirewallRuleProtocol = ret.Protocol
	rule.FirewallRuleSrcIp = ret.SourceIPAddress
	rule.FirewallRuleSrcPort = ret.SourcePort
	rule.FirewallRuleDstIp = ret.DestinationIPAddress
	rule.FirewallRuleDstPort = ret.DestinationPort
}

func getFirewallPolicy(client *gophercloud.ServiceClient,
	policy *firewall.FirewallPolicy,
	wg *sync.WaitGroup) {
	defer wg.Done()

	ret, err := fp.Get(client, policy.FirewallPolicyId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":    err,
			"policy": policy,
		}).Error("query firewall policy failed")
		return
	}

	policy.FirewallPolicyName = ret.Name
	policy.FirewallPolicyDesc = ret.Description
	for _, ruleID := range ret.Rules {
		policy.FirewallPolicyRules = append(policy.FirewallPolicyRules, &firewall.FirewallRule{
			FirewallRuleId: ruleID,
		})
	}
}

func delFirewallRule(client *gophercloud.ServiceClient,
	ruleID string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	err := fr.Delete(client, ruleID).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err":    err,
			"ruleID": ruleID,
		}).Error("delete firewall rule failed")
	}
}

func delFirewallPolicy(client *gophercloud.ServiceClient,
	policyID string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	err := fp.Delete(client, policyID).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err":      err,
			"policyID": policyID,
		}).Error("delete firewall policy failed")
	}
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
		log.WithFields(log.Fields{
			"err":   err,
			"group": group,
		}).Error("call firewall, create group failed")
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
		log.WithFields(log.Fields{
			"err":    err,
			"policy": policy,
		}).Error("call firewall, create policy failed")
		return
	}

	policy.FirewallPolicyName = ret.Name
	policy.FirewallPolicyId = ret.ID
	policy.FirewallPolicyDesc = ret.Description
}

func create_firewall_rules(client *gophercloud.ServiceClient,
	rule *firewall.FirewallRuleSet,
	idx int,
	ruleIDs *[]string,
	policy *firewall.FirewallPolicy,
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
		log.WithFields(log.Fields{
			"err":  err,
			"rule": rule,
		}).Error("call firewall, create rules failed ")
		return
	}

	(*ruleIDs)[idx] = ret.ID
	policy.FirewallPolicyRules[idx] = &firewall.FirewallRule{
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

func DetachFirewallToPortsRawAPI(client *gophercloud.ServiceClient, firewallID string) error {
	url := client.ResourceBase + "fwaas/firewall_groups/" + firewallID

	jsTemplate := `{
	"firewall_group": {
		"ports": []
	}
}`

	mp := map[string]string{}

	jsonReq, err := common.CreateJsonByTmpl(jsTemplate, mp)
	if nil != err {
		return err
	}

	_, err = common.CallRawAPI(url, "PUT", jsonReq, client.TokenID)
	if nil != err {
		return err
	}

	return nil
}
