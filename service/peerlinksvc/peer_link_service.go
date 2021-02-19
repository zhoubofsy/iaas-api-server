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

	"iaas-api-server/proto/peerlink"
)

// PeerLinkService service for peer link
type PeerLinkService struct {
	peerlink.UnimplementedPeerLinkServiceServer
}

//type rpcTask interface {
//	Run(context.Context)
//	exeute(providers *gophercloud.ProviderClient) *common.Error
//	setResult()
//	checkParam() error
//}
//
//// Run rpc start
//func (rpctask *rpcTask) Run(context.Context) {
//	defer rpctask.setResult()
//
//	if err := rpctask.checkParam(); nil != err {
//		log.WithFields(log.Fields{
//			"err": err,
//			"req": rpctask.Req.String(),
//		}).Error("check param failed.")
//		rpctask.Err = common.EPARAM
//		return
//	}
//
//	providers, err := common.GetOpenstackClient(rpctask.Req.Apikey, rpctask.Req.TenantId, rpctask.Req.PlatformUserid)
//	if nil != err {
//		log.WithFields(log.Fields{
//			"err": err,
//			"req": rpctask.Req.String(),
//		}).Error("call common, get openstack client error")
//		rpctask.Err = common.EGETOPSTACKCLIENT
//		return
//	}
//
//	rpctask.Err = rpctask.execute(providers)
//}

// CreatePeerLink create peer link
func (pls *PeerLinkService) CreatePeerLink(ctx context.Context, req *peerlink.PeerLinkReq) (*peerlink.PeerLinkRes, error) {
	task := &CreatePeerLinkRPCTask{
		Req: req,
		Res: &peerlink.PeerLinkRes{},
		Err: nil,
	}

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
