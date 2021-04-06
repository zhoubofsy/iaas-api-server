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
	fp "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
	fr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	log "github.com/sirupsen/logrus"

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
