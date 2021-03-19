/*================================================================
*
*  文件名称：update_security_group_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年01月29日
*
================================================================*/

package securitygroupsvc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	sr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"

	"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
)

// UpdateSecurityGroupRPCTask use for get security group
type UpdateSecurityGroupRPCTask struct {
	Req *securitygroup.UpdateSecurityGroupReq
	Res *securitygroup.SecurityGroupRes
	Err *common.Error
}

// Run first input
func (rpctask *UpdateSecurityGroupRPCTask) Run(context.Context) {
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

func (rpctask *UpdateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return &common.Error{
			Code: common.ENETWORKCLIENT.Code,
			Msg:  err.Error(),
		}
	}

	// 先根据安全组id获取当前已存在的安全组信息
	oldgroup, err := sg.Get(client, rpctask.Req.GetSecurityGroupId()).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("update secgroup, get old group failed.")
		return &common.Error{
			Code: common.ESGUPDATEGROUP.Code,
			Msg:  err.Error(),
		}
	}

	// TODO 下面的操作理论上得考虑事务性，
	// 添加新的安全组规则
	var existRule map[int]bool // 记录一下已经存在的而且又在本地更新的安全组规则
	for _, rule := range rpctask.Req.GetSecurityGroupRuleSets() {
		ropts := sr.CreateOpts{
			SecGroupID:     rpctask.Req.GetSecurityGroupId(),
			Direction:      sr.RuleDirection(rule.GetDirection()),
			PortRangeMin:   int(rule.GetPortRangeMin()),
			PortRangeMax:   int(rule.GetPortRangeMax()),
			Description:    rule.GetRuleDesc(),
			Protocol:       sr.RuleProtocol(rule.GetProtocol()),
			RemoteIPPrefix: rule.GetRemoteIpPrefix(),
			EtherType:      sr.EtherType4, // 兼容 ipv6 需要修改此处
		}

		_, err := sr.Create(client, ropts).Extract()
		if nil != err {
			log.WithFields(log.Fields{
				"err":  err,
				"rule": rule.String(),
			}).Warn("update security, insert new rules failed")

			for idx, oldrule := range oldgroup.Rules {
				if oldrule.Direction == rule.Direction &&
					oldrule.PortRangeMin == int(rule.PortRangeMin) &&
					oldrule.PortRangeMax == int(rule.PortRangeMax) &&
					oldrule.Protocol == rule.Protocol &&
					oldrule.RemoteIPPrefix == rule.RemoteIpPrefix { // 后续兼容ipv6，此处需要添加ipv6的判断
					existRule[idx] = true
					break
				}
			}
			continue
		}
	}

	// 删除旧的安全组规则
	for idx, rule := range oldgroup.Rules {
		if existRule[idx] { // 该旧规则在本次的更新又存在了，
			continue
		}
		err := sr.Delete(client, rule.ID).ExtractErr()
		if nil != err {
			log.WithFields(log.Fields{
				"err":    err,
				"ruleID": rule.ID,
			}).Warn("update security, delete old sec group rule failed")
		}
	}

	// 更新安全组
	uopts := sg.UpdateOpts{
		Name: rpctask.Req.GetSecurityGroupName(),
	}

	newgroup, err := sg.Update(client, rpctask.Req.GetSecurityGroupId(), uopts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("update security, update sec group failed")
		return &common.Error{
			Code: common.ESGUPDATEGROUP.Code,
			Msg:  err.Error(),
		}
	}

	rpctask.Res.SecurityGroup = &securitygroup.SecurityGroupRes_SecurityGroup{
		SecurityGroupId:   newgroup.ID,
		SecurityGroupName: newgroup.Name,
		SecurityGroupDesc: newgroup.Description,
		CreatedTime:       newgroup.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		UpdatedTime:       newgroup.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
	}
	if len(newgroup.Rules) > 0 {
		rpctask.Res.SecurityGroup.SecurityGroupRules = make([]*securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule, len(newgroup.Rules))
		for index, rule := range newgroup.Rules {
			rpctask.Res.SecurityGroup.SecurityGroupRules[index] = &securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule{
				RuleId:          rule.ID,
				RuleDesc:        rule.Description,
				Direction:       rule.Direction,
				Protocol:        rule.Protocol,
				PortRangeMin:    int32(rule.PortRangeMin),
				PortRangeMax:    int32(rule.PortRangeMax),
				RemoteIpPrefix:  rule.RemoteIPPrefix,
				SecurityGroupId: rule.SecGroupID,
				CreatedTime:     newgroup.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				UpdatedTime:     newgroup.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
			}
		}
	}

	return common.EOK
}

func (rpctask *UpdateSecurityGroupRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetSecurityGroupId() ||
		("" == rpctask.Req.GetSecurityGroupName() && "" == rpctask.Req.GetSecurityGroupDesc() && 0 == len(rpctask.Req.GetSecurityGroupRuleSets())) {
		errors.New("input param is wrong")
	}
	return nil
}

func (rpctask *UpdateSecurityGroupRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
