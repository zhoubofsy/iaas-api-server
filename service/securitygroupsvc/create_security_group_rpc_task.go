/*================================================================
*
*  文件名称：create_securitygroup_task.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package securitygroupsvc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	sr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
)

// CreateSecurityGroupRPCTask use for create securty group
type CreateSecurityGroupRPCTask struct {
	Req *securitygroup.CreateSecurityGroupReq
	Res *securitygroup.SecurityGroupRes
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *CreateSecurityGroupRPCTask) Run(context.Context) {
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

func (rpctask *CreateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

	gopts := sg.CreateOpts{
		Name:        rpctask.Req.GetSecurityGroupName(),
		Description: rpctask.Req.GetSecurityGroupDesc(),
	}

	// TODO 事务保证安全组跟安全组规则都是ok的
	group, err := sg.Create(client, gopts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("create security group failed.")
		return &common.Error{
			Code: common.ESGCREATEGROUP.Code,
			Msg:  err.Error(),
		}
	}

	rpctask.Res.SecurityGroup = &securitygroup.SecurityGroupRes_SecurityGroup{
		SecurityGroupId:    group.ID,
		SecurityGroupName:  group.Name,
		SecurityGroupDesc:  group.Description,
		CreatedTime:        group.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		UpdatedTime:        group.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
		SecurityGroupRules: make([]*securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule, 0),
	}

	// 创建安全组会默认创建2条安全组规则，删除这默认的2条规则
	if len(group.Rules) > 0 {
		for _, rule := range group.Rules {
			err = sr.Delete(client, rule.ID).ExtractErr()
			if nil != err {
				log.WithFields(log.Fields{
					"err":             err,
					"ruleID":          rule.ID,
					"securitygroupid": group.ID,
				}).Error("remove default security group rules error")
			}
		}
	}

	// 创建安全组规则
	if nil != rpctask.Req.GetSecurityGroupRuleSets() {
		for _, rule := range rpctask.Req.GetSecurityGroupRuleSets() {
			ropts := sr.CreateOpts{
				Direction:      sr.RuleDirection(rule.GetDirection()),
				Description:    rule.GetRuleDesc(),
				Protocol:       sr.RuleProtocol(rule.GetProtocol()),
				PortRangeMin:   int(rule.GetPortRangeMin()),
				PortRangeMax:   int(rule.GetPortRangeMax()),
				RemoteIPPrefix: rule.GetRemoteIpPrefix(),
				SecGroupID:     group.ID,
				//TODO 网络类型，ipv4，ipv6，proto后续加了加上，默认设置ipv4
				EtherType: sr.EtherType4,
			}

			// TODO 此处得考虑是否考虑事务性
			rl, err := sr.Create(client, ropts).Extract()
			if nil != err {
				log.WithFields(log.Fields{
					"err":  err,
					"rule": rule.String(),
				}).Error("create security group rule failed.")
				continue
			}

			cur := common.Now()
			rpctask.Res.SecurityGroup.SecurityGroupRules = append(rpctask.Res.SecurityGroup.SecurityGroupRules, &securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule{
				RuleId:          rl.ID,
				RuleDesc:        rl.Description,
				Direction:       rl.Direction,
				Protocol:        rl.Protocol,
				PortRangeMin:    int32(rl.PortRangeMin),
				PortRangeMax:    int32(rl.PortRangeMax),
				RemoteIpPrefix:  rl.RemoteIPPrefix,
				SecurityGroupId: rl.SecGroupID,
				UpdatedTime:     cur,
				CreatedTime:     cur,
			})
		}
	}

	return common.EOK
}

func (rpctask *CreateSecurityGroupRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetSecurityGroupName() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *CreateSecurityGroupRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
