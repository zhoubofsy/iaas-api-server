/*================================================================
*
*  文件名称：peer_link_service.go
*  创 建 者: mongia
*  创建日期：2021年02月19日
*
================================================================*/

package peerlinksvc

import (
	"fmt"
	"math/big"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"iaas-api-server/common"
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
	sharedNetID, _ := config.GetString("SharedNetID")
	task.SharedNetID = sharedNetID

	sharedSubnetID, _ := config.GetString("SharedSubnetID")
	task.SharedSubnetID = sharedSubnetID

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
	sharedSubnetID, _ := config.GetString("SharedSubnetID")
	task.SharedSubnetID = sharedSubnetID

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
	sharedSubnetID, _ := config.GetString("SharedSubnetID")
	task.SharedSubnetID = sharedSubnetID

	task.Run(ctx)

	return task.Res, status.Error(codes.OK, "success")
}

func getCIDRBySubnetID(client *gophercloud.ServiceClient,
	subnetID string,
	cidr *string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	subnet, err := subnets.Get(client, subnetID).Extract()
	if nil != err {
		*cidr = ""
	} else {
		*cidr = subnet.CIDR
	}
}

func getRouterByRouterID(client *gophercloud.ServiceClient,
	routerID string,
	router *routers.Router,
	wg *sync.WaitGroup) {
	defer wg.Done()

	ret, err := routers.Get(client, routerID).Extract()
	if nil == err {
		*router = *ret
	}
}

func getPortByRouterIDAndNetID(client *gophercloud.ServiceClient,
	routerID string,
	netID string,
	port *ports.Port,
	wg *sync.WaitGroup) {
	defer wg.Done()

	allPages, err := ports.List(client, ports.ListOpts{
		DeviceID: routerID,
	}).AllPages()

	if nil != err {
		return
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if nil != err {
		return
	}

	for _, pt := range allPorts {
		if len(pt.FixedIPs) > 0 && pt.FixedIPs[0].SubnetID == netID {
			*port = pt
			return
		}
	}
}

func inetntoa(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func inetaton(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func addRouteToRouterByRawAPI(client *gophercloud.ServiceClient, routerID string, destination string, nexthop string) error {
	url := client.ResourceBase + "routers/" + routerID + "/add_extraroutes"

	jsTemplate := `{
	"router": {
		"routes": [{
			"destination": "{{.Destination}}",
			"nexthop": "{{.NextHop}}"
		}]
	}
}`
	mp := map[string]string{
		"Destination": destination,
		"NextHop":     nexthop,
	}

	jsonReq, err := common.CreateJsonByTmpl(jsTemplate, mp)
	if nil != err {
		return err
	}

	_, err = common.CallRawAPI(url, "PUT", jsonReq, client.TokenID)
	if nil != err {
		return err
	}

	return nil
}

func delRouteToRouterByRawAPI(client *gophercloud.ServiceClient, routerID string, destination string, nexthop string) error {
	url := client.ResourceBase + "routers/" + routerID + "/remove_extraroutes"

	jsTemplate := `{
	"router": {
		"routes": [{
			"destination": "{{.Destination}}",
			"nexthop": "{{.NextHop}}"
		}]
	}
}`
	mp := map[string]string{
		"Destination": destination,
		"NextHop":     nexthop,
	}

	jsonReq, err := common.CreateJsonByTmpl(jsTemplate, mp)
	if nil != err {
		return err
	}

	_, err = common.CallRawAPI(url, "PUT", jsonReq, client.TokenID)
	if nil != err {
		return err
	}

	return nil
}
