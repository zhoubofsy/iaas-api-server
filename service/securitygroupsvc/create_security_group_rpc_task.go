/*================================================================
*
*  文件名称：create_securitygroup_task.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package securitygroupsvc

import (
	"golang.org/x/net/context"
	//	"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
)

type CreateSecurityGroupRPCTask struct {
	req *securitygroup.CreateSecurityGroupReq
	res *securitygroup.SecurityGroupRes
	err error
}

func (this *CreateSecurityGroupRPCTask) Run(context.Context) {
	if err := this.checkParam(); nil != err {
		this.err = err
		return
	}
}

func (this *CreateSecurityGroupRPCTask) checkParam() error {
	return nil
}
