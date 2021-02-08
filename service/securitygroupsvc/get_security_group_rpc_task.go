/*================================================================
*
*  文件名称：get_security_group_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年01月29日
*
================================================================*/

package securitygroupsvc

import (
	"errors"

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

func (rpctask *GetSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	group, err := sg.Get(client, rpctask.Req.SecurityGroupId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("call sdk get security group failed.")
		return &common.Error{
			Code: common.ESGGETGROUP.Code,
			Msg:  err.Error(),
		}
	}

	//TODO 时间返回后续修改为接口需要的格式
	rpctask.Res.SecurityGroup = &securitygroup.SecurityGroupRes_SecurityGroup{
		SecurityGroupId:   group.ID,
		SecurityGroupName: group.Name,
		SecurityGroupDesc: group.Description,
		CreatedTime:       group.CreatedAt.String(),
		UpdatedTime:       group.UpdatedAt.String(),
	}
	if len(group.Rules) > 0 {
		cur := getCurTime()
		rpctask.Res.SecurityGroup.SecurityGroupRules = make([]*securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule, len(group.Rules))
		for index, rule := range group.Rules {
			rpctask.Res.SecurityGroup.SecurityGroupRules[index] = &securitygroup.SecurityGroupRes_SecurityGroup_SecurityGroupRule{
				RuleId:          rule.ID,
				RuleDesc:        rule.Description,
				Direction:       rule.Direction,
				Protocol:        rule.Protocol,
				PortRangeMin:    int32(rule.PortRangeMin),
				PortRangeMax:    int32(rule.PortRangeMax),
				RemoteIpPrefix:  rule.RemoteIPPrefix,
				SecurityGroupId: rule.SecGroupID,
				CreatedTime:     cur,
				UpdatedTime:     cur,
			}
		}
	}

	return common.EOK
}

func (rpctask *GetSecurityGroupRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetSecurityGroupId() ||
		"" == rpctask.Req.GetTenantId() {
		return errors.New("imput param is wrong")
	}
	return nil
}

func (rpctask *GetSecurityGroupRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
