/*================================================================
*
*  文件名称：create_securitygroup_task.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
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

// CreateSecurityGroupRPCTask use for create securty group
type CreateSecurityGroupRPCTask struct {
	Req *securitygroup.CreateSecurityGroupReq
	Res *securitygroup.SecurityGroupRes
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *CreateSecurityGroupRPCTask) Run(context.Context) {
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
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	createOpts := sg.CreateOpts{
		Name:        rpctask.Req.SecurityGroupName,
		Description: rpctask.Req.SecurityGroupDesc,
	}

	group, err := sg.Create(client, createOpts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("create security group failed.")
		return common.ESGCREATEGROUP
	}

	//TODO 返回参数
	rpctask.Res.SecurityGroup.SecurityGroupId = group.ID
	rpctask.Res.SecurityGroup.SecurityGroupName = group.Name
	rpctask.Res.SecurityGroup.SecurityGroupDesc = group.Description
	rpctask.Res.SecurityGroup.CreatedTime = group.CreatedAt.String()
	rpctask.Res.SecurityGroup.UpdatedTime = group.UpdatedAt.String()

	return common.EOK
}

func (rpctask *CreateSecurityGroupRPCTask) checkParam() error {
	return nil
}
