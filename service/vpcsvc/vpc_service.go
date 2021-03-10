package vpcsvc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iaas-api-server/proto/vpc"
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
	return task.Res, status.Error(codes.OK, "success")
}

func (vpcs *VpcService) GetVpcInfo(ctx context.Context, req *vpc.GetVpcInfoReq) (*vpc.VpcRes, error) {
	task := GetVpcInfoRPCTask{
		Req: req,
		Res: &vpc.VpcRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

func (vpcs *VpcService) SetVpcInfo(ctx context.Context, req *vpc.SetVpcInfoReq) (*vpc.VpcRes, error) {
	task := SetVpcInfoRPCTask{
		Req: req,
		Res: &vpc.VpcRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}
