/*================================================================
*
*  文件名称：get_peer_link_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年02月19日
*
================================================================*/

package peerlinksvc

import (
	"errors"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/peerlink"
)

// GetPeerLinkRPCTask rpc task
type GetPeerLinkRPCTask struct {
	Req        *peerlink.PeerLinkReq
	Res        *peerlink.PeerLinkRes
	Err        *common.Error
	ShareNetID string
}

// Run rpc start
func (rpctask *GetPeerLinkRPCTask) Run(context.Context) {
	defer rpctask.setResult()

	if err := rpctask.checkParam(); nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("check param failed.")
		rpctask.Err = common.EPARAM
		return
	}

	providers, err := common.GetOpenstackClient(rpctask.Req.Apikey, rpctask.Req.TenantId, rpctask.Req.PlatformUserid)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("call common, get openstack client error")
		rpctask.Err = common.EGETOPSTACKCLIENT
		return
	}

	rpctask.Err = rpctask.execute(providers)
}

func (rpctask *GetPeerLinkRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

	var wg sync.WaitGroup
	var subnetACIDR, subnetBCIDR string
	var routerA, routerB routers.Router
	var portsA, portsB ports.Port

	// TODO 以下接口都可能出现失败，聚合操作保证事务性，后续得修改
	// 异步获取peera和peerb的子网cidr，后面设置需要用到
	{
		wg.Add(1)
		go getCIDRBySubnetID(client,
			rpctask.Req.GetPeerASubnetid(),
			&subnetACIDR,
			&wg)
	}

	{
		wg.Add(1)
		go getCIDRBySubnetID(client,
			rpctask.Req.GetPeerBSubnetid(),
			&subnetBCIDR,
			&wg)
	}

	// 获取路由器信息
	{
		wg.Add(1)
		go getRouterByRouterID(client,
			rpctask.Req.GetPeerARouterid(),
			&routerA,
			&wg)
	}

	{
		wg.Add(1)
		go getRouterByRouterID(client,
			rpctask.Req.GetPeerBRouterid(),
			&routerB,
			&wg)
	}

	// 根据路由器ID，从端口列表中找出该路由器的所有port，并取出子网ID是share网络的port
	{
		wg.Add(1)
		go getPortByRouterIDAndNetID(client,
			rpctask.Req.GetPeerARouterid(),
			rpctask.ShareNetID,
			&portsA,
			&wg)
	}

	{
		wg.Add(1)
		go getPortByRouterIDAndNetID(client,
			rpctask.Req.GetPeerBRouterid(),
			rpctask.ShareNetID,
			&portsB,
			&wg)
	}

	wg.Wait()
	if "" == subnetACIDR || "" == subnetBCIDR ||
		"" == routerA.ID || "" == routerB.ID ||
		"" == portsA.ID || "" == portsB.ID {
		log.Error("get cidr or router or ports failed")
		return common.EPLGETPREPARE
	}

	// 设置路由器A的对端路由B的路由表信息
	for _, route := range routerA.Routes {
		if route.DestinationCIDR == subnetBCIDR {
			rpctask.Res.LinkConfOnPeerA = &peerlink.PeerLinkRes_LinkConf{
				CreatedTime: getCurTime(),
				RouteToPeer: &peerlink.PeerLinkRes_LinkConf_Route{
					Destination: route.DestinationCIDR,
					Nexthop:     route.NextHop,
				},
				IntfId: portsA.ID,
				IntfIp: portsA.FixedIPs[0].IPAddress,
			}
		}
	}

	// 设置路由器B的对端路由A的路由表信息
	for _, route := range routerB.Routes {
		if route.DestinationCIDR == subnetACIDR {
			rpctask.Res.LinkConfOnPeerB = &peerlink.PeerLinkRes_LinkConf{
				CreatedTime: getCurTime(),
				RouteToPeer: &peerlink.PeerLinkRes_LinkConf_Route{
					Destination: route.DestinationCIDR,
					Nexthop:     route.NextHop,
				},
				IntfId: portsB.ID,
				IntfIp: portsB.FixedIPs[0].IPAddress,
			}
		}
	}

	return common.EOK
}

func (rpctask *GetPeerLinkRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetPeerARouterid() ||
		"" == rpctask.Req.GetPeerASubnetid() ||
		"" == rpctask.Req.GetPeerBRouterid() ||
		"" == rpctask.Req.GetPeerBSubnetid() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *GetPeerLinkRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
}
