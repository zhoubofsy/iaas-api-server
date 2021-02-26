package routesvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/proto/route"
	"time"
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
	return task.Res, task.Err
}

func (rs *RouteService) SetRoutes(ctx context.Context, req *route.SetRoutesReq) (*route.SetRoutesRes, error) {
	task := SetRoutesRPCTask{
		Req: req,
		Res: &route.SetRoutesRes{},
		Err: nil,
	}
	task.Run(ctx)
	return task.Res, task.Err
}

func getCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
