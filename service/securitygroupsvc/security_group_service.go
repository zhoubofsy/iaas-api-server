/*================================================================
*
*  文件名称：securitygroup_service.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package securitygroupsvc

import (
	"golang.org/x/net/context"
	//"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
	//"unicode"
)

// SecurityGroupService service for security group
type SecurityGroupService struct {
	securitygroup.UnimplementedSecurityGroupServiceServer
}

// CreateSecurityGroup create security group
func (sgs *SecurityGroupService) CreateSecurityGroup(ctx context.Context, req *securitygroup.CreateSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	task := CreateSecurityGroupRPCTask{
		req: req,
		res: &securitygroup.SecurityGroupRes{},
		err: nil,
	}
	task.Run(ctx)
	return task.res, task.err
}

// GetSecurityGroup get security group
func (sgs *SecurityGroupService) GetSecurityGroup(context.Context, *securitygroup.GetSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	return &securitygroup.SecurityGroupRes{}, nil
}

// UpdateSecurityGroup update security group
func (sgs *SecurityGroupService) UpdateSecurityGroup(context.Context, *securitygroup.UpdateSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	return &securitygroup.SecurityGroupRes{}, nil
}

// DeleteSecurityGroup delete security group
func (sgs *SecurityGroupService) DeleteSecurityGroup(context.Context, *securitygroup.DeleteSecurityGroupReq) (*securitygroup.DeleteSecurityGroupRes, error) {
	return &securitygroup.DeleteSecurityGroupRes{}, nil
}

// OperateSecurityGroup operate security group
func (sgs *SecurityGroupService) OperateSecurityGroup(context.Context, *securitygroup.OperateSecurityGroupReq) (*securitygroup.OperateSecurityGroupRes, error) {
	return &securitygroup.OperateSecurityGroupRes{}, nil
}
