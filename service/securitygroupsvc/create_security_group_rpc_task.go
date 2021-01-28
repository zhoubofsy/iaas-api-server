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
	req *securitygroup.CreateSecurityGroupReq
	res *securitygroup.SecurityGroupRes
	err error
}

// Run call this func for doing task
func (rpctask *CreateSecurityGroupRPCTask) Run(context.Context) {
	if err := rpctask.checkParam(); nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("check param failed.")
		return
	}

	if !common.APIAuth(rpctask.req.Apikey, rpctask.req.TenantId, rpctask.req.PlatformUserid) {
		log.Error("call common, api auth error")
		return
	}

	providers, err := common.GetOpenstackClient(rpctask.req.Apikey, rpctask.req.TenantId, rpctask.req.PlatformUserid)
	if nil != err {
		log.Error("call common, get openstack client error")
		return
	}

	rpctask.execute(providers)
}

func (rpctask *CreateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("new network v2 failed.")
		return
	}

	createOpts := sg.CreateOpts{
		Name:        rpctask.req.SecurityGroupName,
		Description: rpctask.req.SecurityGroupDesc,
	}

	group, err := sg.Create(client, createOpts).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("create security group failed.")
		return
	}

	//TODO 返回参数
	rpctask.res.SecurityGroup.SecurityGroupId = group.ID

}

func (rpctask *CreateSecurityGroupRPCTask) checkParam() error {
	return nil
}
