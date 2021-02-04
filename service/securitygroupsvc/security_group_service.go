/*================================================================
*
*  文件名称：securitygroup_service.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package securitygroupsvc

import (
	"time"

	"golang.org/x/net/context"
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
	return task.Res, task.Err
}

// GetSecurityGroup get security group
func (sgs *SecurityGroupService) GetSecurityGroup(ctx context.Context, req *securitygroup.GetSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	task := GetSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.SecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

// UpdateSecurityGroup update security group
func (sgs *SecurityGroupService) UpdateSecurityGroup(ctx context.Context, req *securitygroup.UpdateSecurityGroupReq) (*securitygroup.SecurityGroupRes, error) {
	task := UpdateSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.SecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

// DeleteSecurityGroup delete security group
func (sgs *SecurityGroupService) DeleteSecurityGroup(ctx context.Context, req *securitygroup.DeleteSecurityGroupReq) (*securitygroup.DeleteSecurityGroupRes, error) {
	task := DeleteSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.DeleteSecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

// OperateSecurityGroup operate security group
func (sgs *SecurityGroupService) OperateSecurityGroup(ctx context.Context, req *securitygroup.OperateSecurityGroupReq) (*securitygroup.OperateSecurityGroupRes, error) {
	task := OperateSecurityGroupRPCTask{
		Req: req,
		Res: &securitygroup.OperateSecurityGroupRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

func getCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}