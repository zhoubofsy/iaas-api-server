/*================================================================
*
*  文件名称：delete_security_group_rpc_task.go
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

// DeleteSecurityGroupRPCTask use for get security group
type DeleteSecurityGroupRPCTask struct {
	Req *securitygroup.DeleteSecurityGroupReq
	Res *securitygroup.DeleteSecurityGroupRes
	Err *common.Error
}

// Run call this func in DeleteSecurityGroupRPCTask object
func (rpctask *DeleteSecurityGroupRPCTask) Run(context.Context) {
	defer func() {
		rpctask.Res.Code = rpctask.Err.Code
		rpctask.Res.Msg = rpctask.Err.Msg
		rpctask.Res.SecurityGroupId = rpctask.Req.GetSecurityGroupId()
		rpctask.Res.DeletedTime = getCurTime()
	}()

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

func (rpctask *DeleteSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	err = sg.Delete(client, rpctask.Req.GetSecurityGroupId()).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("delete security group failed")
		return &common.Error{
			Code: common.ESGDELGROUP.Code,
			Msg:  err.Error(),
		}
	}

	return common.EOK
}

func (rpctask *DeleteSecurityGroupRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetSecurityGroupId() {
		return errors.New("input params is wrong")
	}
	return nil
}
