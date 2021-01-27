/*================================================================
*
*  文件名称：security_group_service.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package security_group_service

import (
	"golang.org/x/net/context"
	"iaas-api-server/proto/security_group"
	//"unicode"
)

type SecurityGroupService struct {
	security_group.UnimplementedSecurityGroupServiceServer
}

func (sgs *SecurityGroupService) CreateSecurityGroup(context.Context, *security_group.CreateSecurityGroupReq) (*security_group.SecurityGroup, error) {
	return &security_group.SecurityGroup{}, nil
}

func (sgs *SecurityGroupService) GetSecurityGroup(context.Context, *security_group.GetSecurityGroupReq) (*security_group.SecurityGroup, error) {
	return &security_group.SecurityGroup{}, nil
}

func (sgs *SecurityGroupService) UpdateSecurityGroup(context.Context, *security_group.UpdateSecurityGroupReq) (*security_group.SecurityGroup, error) {
	return &security_group.SecurityGroup{}, nil
}

func (sgs *SecurityGroupService) DeleteSecurityGroup(context.Context, *security_group.DeleteSecurityGroupReq) (*security_group.DeleteSecurityGroupRes, error) {
	return &security_group.DeleteSecurityGroupRes{}, nil
}

func (sgs *SecurityGroupService) OperateSecurityGroup(context.Context, *security_group.OperateSecurityGroupReq) (*security_group.OperateSecurityGroupRes, error) {
	return &security_group.OperateSecurityGroupRes{}, nil
}
