/*================================================================
*
*  文件名称：create_security_group_task.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package security_group_service

import (
	"golang.org/x/net/context"
	//	"iaas-api-server/common"
	"iaas-api-server/proto/security_group"
)

type CreateSecurityGroupTask struct {
	req *security_group.CreateSecurityGroupReq
	res *security_group.SecurityGroup
	err error
}

func (this *CreateSecurityGroupTask) Run(context.Context) {
	if err := this.checkParam(); nil != err {
		this.err = err
		return
	}
}

func (this *CreateSecurityGroupTask) checkParam() error {
	return nil
}
