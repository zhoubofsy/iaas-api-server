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
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"iaas-api-server/proto/securitygroup"
)

// SecurityGroupService service for security group
type SecurityGroupService struct {
	securitygroup.UnimplementedSecurityGroupServiceServer
}

// CreateSecurityGroup create security group
func (sgs *SecurityGroupService) CreateSecurityGroup(ctx context.Context, req *securitygroup.CreateSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	task := CreateSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.SecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

// GetSecurityGroup get security group
func (sgs *SecurityGroupService) GetSecurityGroup(ctx context.Context, req *securitygroup.GetSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	task := GetSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.SecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

// UpdateSecurityGroup update security group
func (sgs *SecurityGroupService) UpdateSecurityGroup(ctx context.Context, req *securitygroup.UpdateSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	task := UpdateSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.SecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

// DeleteSecurityGroup delete security group
func (sgs *SecurityGroupService) DeleteSecurityGroup(ctx context.Context, req *securitygroup.DeleteSecurityGroupReq) (*securitygroup.DeleteSecurityGroupRes, error) {
	task := DeleteSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.DeleteSecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

// OperateSecurityGroup operate security group
func (sgs *SecurityGroupService) OperateSecurityGroup(ctx context.Context, req *securitygroup.OperateSecurityGroupReq) (*securitygroup.OperateSecurityGroupRes, error) {
	task := OperateSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.OperateSecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}
