package routesvc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iaas-api-server/proto/route"
)

type RouteService struct {
	route.UnimplementedRouterServiceServer
}

func (rs *RouteService) GetRouter(ctx context.Context, req *route.GetRouterReq) (*route.GetRouterRes, error) {
	task := GetRouterRPCTask{
		Req: req,
		Res: &route.GetRouterRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}

func (rs *RouteService) SetRoutes(ctx context.Context, req *route.SetRoutesReq) (*route.SetRoutesRes, error) {
	task := SetRoutesRPCTask{
		Req: req,
		Res: &route.SetRoutesRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, status.Error(codes.OK, "success")
}
