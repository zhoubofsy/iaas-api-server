package vpcsvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/proto/vpc"
	"time"
)

type VpcService struct {
	vpc.UnimplementedVpcServiceServer
}

func (vpcs *VpcService) CreateVpc(ctx context.Context, req *vpc.CreateVpcReq) (*vpc.VpcRes, error) {
	task := CreateVpcRPCTask{
		Req: req,
		Res: &vpc.VpcRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

func (vpcs *VpcService) GetVpcInfo(ctx context.Context, req *vpc.GetVpcInfoReq) (*vpc.VpcRes, error) {
	task := GetVpcInfoRPCTask{
		Req: req,
		Res: &vpc.VpcRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

func (vpcs *VpcService) SetVpcInfo(ctx context.Context, req *vpc.SetVpcInfoReq) (*vpc.VpcRes, error) {
	task := SetVpcInfoRPCTask{
		Req: req,
		Res: &vpc.VpcRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

func getCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
