/*================================================================
*
*  文件名称：operate_security_group_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年01月29日
*
================================================================*/

package securitygroupsvc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	//	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
)

// OperateSecurityGroupRPCTask use for get security group
type OperateSecurityGroupRPCTask struct {
	Req *securitygroup.OperateSecurityGroupReq
	Res *securitygroup.OperateSecurityGroupRes
	Err *common.Error
}

// Run call this func
func (rpctask *OperateSecurityGroupRPCTask) Run(context.Context) {
	defer func() {
		rpctask.Res.Code = rpctask.Err.Code
		rpctask.Res.Msg = rpctask.Err.Msg
	}()

	if err := rpctask.checkParam(); nil != err {
		log.WithFields(log.Fields{
			"err":             err,
			"apikey":          rpctask.Req.GetApikey(),
			"tenant_id":       rpctask.Req.GetTenantId(),
			"platform_userid": rpctask.Req.GetPlatformUserid(),
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

func (rpctask *OperateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err":             err,
			"apikey":          rpctask.Req.GetApikey(),
			"tenant_id":       rpctask.Req.GetTenantId(),
			"platform_userid": rpctask.Req.GetPlatformUserid(),
			"client":          client,
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	return common.EOK
}

func (rpctask *OperateSecurityGroupRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetSecurityGroupId() ||
		"" == rpctask.Req.GetOpsType() ||
		0 == len(rpctask.Req.GetInstanceIds()) {
		return errors.New("input param is wrong")
	}
	return nil
}
