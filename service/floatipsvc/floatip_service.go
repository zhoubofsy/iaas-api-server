package floatipsvc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iaas-api-server/common/config"
	"iaas-api-server/proto/floatip"
)

type FloatIpService struct {
	floatip.UnimplementedFloatIpServiceServer
}

func (fis *FloatIpService) BindFloatIpToInstance(ctx context.Context, req *floatip.BindFloatIpToInstanceReq) (*floatip.BindFloatIpToInstanceRes, error) {
	task := BindFloatIpToInstanceTask{
		Req: req,
		Res: &floatip.BindFloatIpToInstanceRes{},
		Err: nil,
	}
	PublicNetID, _ := config.GetString("PublicNetID")
	task.PublicNetID = PublicNetID

	PublicSubnetID, _ := config.GetString("PublicSubnetID")
	task.PublicSubnetID = PublicSubnetID

	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

func (fis *FloatIpService) RevokeFloatIpFromInstance(ctx context.Context, req *floatip.RevokeFloatIpFromInstanceReq) (*floatip.RevokeFloatIpFromInstanceRes, error) {
	task := RevokeFloatIpFromInstanceTask{
		Req: req,
		Res: &floatip.RevokeFloatIpFromInstanceRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}
