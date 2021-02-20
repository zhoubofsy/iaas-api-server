/*================================================================
*
*  文件名称：peer_link_service.go
*  创 建 者: mongia
*  创建日期：2021年02月19日
*
================================================================*/

package peerlinksvc

import (
	"time"

	"golang.org/x/net/context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"iaas-api-server/common/config"
	"iaas-api-server/proto/peerlink"
)

// PeerLinkService service for peer link
type PeerLinkService struct {
	peerlink.UnimplementedPeerLinkServiceServer
}

// CreatePeerLink create peer link
func (pls *PeerLinkService) CreatePeerLink(ctx context.Context, req *peerlink.PeerLinkReq) (*peerlink.PeerLinkRes, error) {
	task := &CreatePeerLinkRPCTask{
		Req: req,
		Res: &peerlink.PeerLinkRes{},
		Err: nil,
	}
	shareNetID, _ := config.GetString("ShareNetID")
	task.ShareNetID = shareNetID

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

// GetPeerLink get peer link
func (pls *PeerLinkService) GetPeerLink(ctx context.Context, req *peerlink.PeerLinkReq) (*peerlink.PeerLinkRes, error) {
	task := &GetPeerLinkRPCTask{
		Req: req,
		Res: &peerlink.PeerLinkRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

// DeletePeerLink delete peer link
func (pls *PeerLinkService) DeletePeerLink(ctx context.Context, req *peerlink.PeerLinkReq) (*peerlink.DeletePeerLinkRes, error) {
	task := &DeletePeerLinkRPCTask{
		Req: req,
		Res: &peerlink.DeletePeerLinkRes{},
		Err: nil,
	}

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func getCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type routeInfo struct {
	routeID string
	cidr    string
}

//type baseRPCTask interface {
//	getOCIRBySubnetID(client *gophercloud.ServiceClient, subnetID string) (string, error)
//}
//
//func (base *baseRPCTask) getOCIRBySubnetID(client *gophercloud.ServiceClient, subnetID string) (string, error) {
//	subnet, err := subnets.Get(client, subnetID).Extract()
//	if nil != err {
//		return "", err
//	}
//
//	return subnet.CIDR, nil
//}
