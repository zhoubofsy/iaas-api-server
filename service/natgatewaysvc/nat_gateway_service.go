/*================================================================
*
*  文件名称：nat_gateway_service.go
*  创 建 者: mongia
*  创建日期：2021年02月03日
*
================================================================*/

package natgatewaysvc

import (
	"time"

	"golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"iaas-api-server/proto/natgateway"
)

// NatGatewayService service for nat gateway
type NatGatewayService struct {
	natgateway.UnimplementedNatGatewayServiceServer
}

// CreateNatGateway create nat gateway
func (ngs *NatGatewayService) CreateNatGateway(ctx context.Context, req *natgateway.CreateNatGatewayReq) (*natgateway.NatGatewayRes, error) {
	task := &CreateNatGatewayRPCTask{
		Req: req,
		Res: &natgateway.NatGatewayRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

// GetNatGateway get nat gateway
func (ngs *NatGatewayService) GetNatGateway(ctx context.Context, req *natgateway.GetNatGatewayReq) (*natgateway.NatGatewayRes, error) {
	task := &GetNatGatewayRPCTask{
		Req: req,
		Res: &natgateway.NatGatewayRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

// DeleteNatGateway delete nat gateway
func (ngs *NatGatewayService) DeleteNatGateway(ctx context.Context, req *natgateway.DeleteNatGatewayReq) (*natgateway.DeleteNatGatewayRes, error) {
	task := &DeleteNatGatewayRPCTask{
		Req: req,
		Res: &natgateway.DeleteNatGatewayRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func getCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
