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

	nsg "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
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

var (
	attachInstance string = "attach"
	detachInstance string = "detach"
)

// Run call this func
func (rpctask *OperateSecurityGroupRPCTask) Run(context.Context) {
	defer func() {
		rpctask.Res.Code = rpctask.Err.Code
		rpctask.Res.Msg = rpctask.Err.Msg
		rpctask.Res.OperateedTime = getCurTime()
		rpctask.Res.OpsType = rpctask.Req.GetOpsType()
		rpctask.Res.SecurityGroupId = rpctask.Req.GetSecurityGroupId()
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

func (rpctask *OperateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	netclient, neterr := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})
	novaclient, novaerr := openstack.NewComputeV2(providers, gophercloud.EndpointOpts{})

	if nil != neterr || nil != novaerr {
		log.WithFields(log.Fields{
			"neutron err": neterr,
			"nova err":    novaerr,
			"req":         rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ESGNEWNETWORK
	}

	// 获取安全组，保证要操作的安全组存在
	_, err := sg.Get(netclient, rpctask.Req.GetSecurityGroupId()).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("operate security, get sec group failed")
		return &common.Error{
			Code: common.ESGOPERGROUP.Code,
			Msg:  err.Error(),
		}
	}

	// 根据传入的操作类型进行绑定跟解绑操作
	// TODO 给多个实例绑定安全组，可能中间绑定失败，理论上得考虑事务性
	if attachInstance == rpctask.Req.GetOpsType() {
		for _, instanceID := range rpctask.Req.GetInstanceIds() {
			err := nsg.AddServer(novaclient, instanceID, rpctask.Req.GetSecurityGroupId()).ExtractErr()
			if nil != err {
				log.WithFields(log.Fields{
					"err":        err,
					"instanceid": instanceID,
					"secgroupid": rpctask.Req.GetSecurityGroupId(),
				}).Warn("operate security, attach instance failed")
				continue
			}
		}
	} else if detachInstance == rpctask.Req.GetOpsType() {
		for _, instanceID := range rpctask.Req.GetInstanceIds() {
			err := nsg.RemoveServer(novaclient, instanceID, rpctask.Req.GetSecurityGroupId()).ExtractErr()
			if nil != err {
				log.WithFields(log.Fields{
					"err":        err,
					"instanceid": instanceID,
					"secgroupid": rpctask.Req.GetSecurityGroupId(),
				}).Warn("operate security, attach instance failed")
				continue
			}
		}
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
