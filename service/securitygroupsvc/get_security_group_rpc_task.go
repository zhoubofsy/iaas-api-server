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
	//	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
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

	if !common.APIAuth(rpctask.Req.Apikey, rpctask.Req.TenantId, rpctask.Req.PlatformUserid) {
		log.Error("call common, api auth error")
		rpctask.Err = common.EAPIAUTH
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
	//TODO call sdk api

	return common.EOK
}

func (rpctask *GetSecurityGroupRPCTask) checkParam() error {
	return nil
}
