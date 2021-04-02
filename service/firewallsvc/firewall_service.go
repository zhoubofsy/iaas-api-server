/*================================================================
*
*  文件名称：firewall_service.go
*  创 建 者: mongia
*  创建日期：2021年04月02日
*
================================================================*/

package firewallsvc

import (
	"golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"iaas-api-server/proto/firewall"
)

type FirewallService struct {
    firewall.UnimplementedFirewallServiceServer
}

func (pthis *FirewallService) GetFirewall(ctx context.Context, req *firewall.GetFirewallReq) (*firewall.FirewallRes, error) {
    task := GetFirewallRPCTask{
        Req: req,
        Res: &firewall.FirewallRes{},
        Err: nil,
    }

    task.Run(ctx)

    return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) CreateFirewall(ctx context.Context, req *firewall.CreateFirewallReq) (*firewall.FirewallRes, error) {
    task := CreateFirewallRPCTask{
        Req: req,
        Res: &firewall.FirewallRes{},
        Err: nil,
    }

    task.Run(ctx)

    return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) UpdateFirewall(ctx context.Context, req *firewall.UpdateFirewallReq) (*firewall.FirewallRes, error) {
    task := UpdateFirewallRPCTask{
        Req: req,
        Res: &firewall.FirewallRes{},
        Err: nil,
    }

    task.Run(ctx)

    return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) DeleteFirewall(ctx context.Context, req *firewall.DeleteFirewallReq) (*firewall.DeleteFirewallRes, error) {
    task := DeleteFirewallRPCTask{
        Req: req,
        Res: &firewall.DeleteFirewallRes{},
        Err: nil,
    }

    task.Run(ctx)

    return task.Res, status.Error(codes.OK, "success")
}

func (pthis *FirewallService) OperateFirewall(ctx context.Context, req *firewall.OperateFirewallReq) (*firewall.OperateFirewallRes, error) {
    task := OperateFirewallRPCTask{
        Req: req,
        Res: &firewall.OperateFirewallRes{},
        Err: nil,
    }

    task.Run(ctx)

    return task.Res, status.Error(codes.OK, "success")
}
