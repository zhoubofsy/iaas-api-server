/*================================================================
*
*  文件名称：get_security_group_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年01月29日
*
================================================================*/

package securitygroupsvc

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
)

// GetSecurityGroupRPCTask use for get security group
type GetSecurityGroupRPCTask struct {
	Req *securitygroup.GetSecurityGroupReq
	Res *securitygroup.SecurityGroupRes
	Err *common.Error
}

// Run first input
func (rpctask *GetSecurityGroupRPCTask) Run(context.Context) {
	if err := rpctask.checkParam(); nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("check param failed.")
		rpctask.Err = common.EPARAM
		return
	}

	providers, err := common.GetOpenstackClient(rpctask.Req.Apikey, rpctask.Req.TenantId, rpctask.Req.PlatformUserid)
	if nil != err {
		log.Error("call common, get openstack client error")
		rpctask.Err = common.EGETOPSTACKCLIENT
		return
	}

	rpctask.Err = rpctask.execute(providers)
}

func (rpctask *GetSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	log.WithFields(log.Fields{
		"client": client,
	}).Info("client")

	groups, err := sg.Get(client, rpctask.Req.SecurityGroupId).Extract()
	if nil != err {
		return common.ESGGETGROUP
	}

	rpctask.Res.SecurityGroup.UpdatedTime = groups.UpdatedAt.String()
	rpctask.Res.SecurityGroup.CreatedTime = groups.CreatedAt.String()
	rpctask.Res.SecurityGroup.SecurityGroupId = groups.ID
	rpctask.Res.SecurityGroup.SecurityGroupName = groups.Name
	rpctask.Res.SecurityGroup.SecurityGroupDesc = groups.Description

	if len(groups.Rules) > 0 {
		rpctask.Res.SecurityGroup.SecurityGroupRules = make([]*securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule, len(groups.Rules))
		for index, rule := range groups.Rules {
			rpctask.Res.SecurityGroup.SecurityGroupRules[index] = &securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule{
				RuleId:          rule.ID,
				RuleDesc:        rule.Description,
				Direction:       rule.Direction,
				Protocol:        rule.Protocol,
				PortRangeMin:    int32(rule.PortRangeMin),
				PortRangeMax:    int32(rule.PortRangeMax),
				RemoteIpPrefix:  rule.RemoteIPPrefix,
				SecurityGroupId: rule.SecGroupID,
				CreatedTime:     groups.CreatedAt.String(),
				UpdatedTime:     groups.UpdatedAt.String(),
			}
		}
	}

	return common.EOK
}

func (rpctask *GetSecurityGroupRPCTask) checkParam() error {
	return nil
}
