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
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	//sr "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
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
	wg  *sync.WaitGroup
}

// Run call this func for doing task
func (rpctask *CreateSecurityGroupRPCTask) Run(context.Context) {
	defer func() {
		rpctask.Res.Code = rpctask.Err.Code
		rpctask.Res.Msg = rpctask.Err.Msg
	}()

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

func (rpctask *CreateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	gopts := sg.CreateOpts{
		Name:        rpctask.Req.SecurityGroupName,
		Description: rpctask.Req.SecurityGroupDesc,
	}

	group, err := sg.Create(client, gopts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("create security group failed.")
		return &common.Error{
			Code: common.ESGCREATEGROUP.Code,
			Msg:  err.Error(),
		}
	}

	rpctask.Res.SecurityGroup.SecurityGroupId = group.ID
	rpctask.Res.SecurityGroup.SecurityGroupName = group.Name
	rpctask.Res.SecurityGroup.SecurityGroupDesc = group.Description
	rpctask.Res.SecurityGroup.CreatedTime = group.CreatedAt.String()
	rpctask.Res.SecurityGroup.UpdatedTime = group.UpdatedAt.String()

	//	if nil != rpctask.Req.GetSecurityGroupRuleSets() {
	//		for _, rule := range rpctask.Req.GetSecurityGroupRuleSets() {
	//			ropts := sr.CreateOpts{
	//				Direction:      sr.RuleDirection(rule.GetDirection()),
	//				Description:    rule.GetRuleDesc(),
	//				Protocol:       sr.RuleProtocol(rule.GetProtocol()),
	//				PortRangeMin:   int(rule.GetPortRangeMin()),
	//				PortRangeMax:   int(rule.GetPortRangeMax()),
	//				RemoteIPPrefix: rule.GetRemoteIpPrefix(),
	//				SecGroupID:     group.ID,
	//				//TODO 网络类型，ipv4，ipv6，proto后续加了加上，默认设置ipv4
	//				EtherType: sr.EtherType4,
	//			}
	//
	//			rl, err := sr.Create(client, ropts).Extract()
	//			if nil != err {
	//				log.WithFields(log.Fields{
	//					"err":            err,
	//					"Direction":      rule.GetDirection(),
	//					"Description":    rule.GetRuleDesc(),
	//					"Protocol":       rule.GetProtocol(),
	//					"PortRangeMin":   rule.GetPortRangeMin(),
	//					"PortRangeMax":   rule.GetPortRangeMax(),
	//					"RemoteIPPrefix": rule.GetRemoteIpPrefix(),
	//					"SecGroupID":     group.ID,
	//					//TODO 网络类型，ipv4，ipv6，proto后续加了加上，默认设置ipv4
	//					"EtherType": sr.EtherType4,
	//				}).Error("create security group rule failed.")
	//				continue
	//			}
	//		}
	//	}
	return common.EOK
}

func (rpctask *CreateSecurityGroupRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() {
		return errors.New("input params is wrong")
	}
	return nil
}
