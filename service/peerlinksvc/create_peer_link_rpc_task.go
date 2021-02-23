/*================================================================
*
*  文件名称：create_peer_link_rpc_task.go
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

	//	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/peerlink"
)

// CreatePeerLinkRPCTask rpc task
type CreatePeerLinkRPCTask struct {
	Req        *peerlink.PeerLinkReq
	Res        *peerlink.PeerLinkRes
	Err        *common.Error
	ShareNetID string
}

// Run rpc start
func (rpctask *CreatePeerLinkRPCTask) Run(context.Context) {
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

func (rpctask *CreatePeerLinkRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
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
	var availableIP []string = make([]string, 0)

	// 异步获取peera和peerb的子网cidr，后面设置需要用到, 从share net pool获取2个ip
	{
		wg.Add(1)
		go getCIDRBySubnetID(client, rpctask.Req.GetPeerASubnetid(), &subnetACIDR, &wg)
	}

	{
		wg.Add(1)
		go getCIDRBySubnetID(client, rpctask.Req.GetPeerBSubnetid(), &subnetBCIDR, &wg)
	}

	{
		wg.Add(1)
		go getIPFromSubnet(client, rpctask.ShareNetID, &availableIP, &wg)
	}

	wg.Wait()
	if "" == subnetACIDR || "" == subnetBCIDR || 0 == len(availableIP) {
		log.Error("prepare create peer link data failed")
		return common.EPLCREATEPREPARE
	}

	// TODO 异步创建interface（route -> share subnet）
	var peerA, peerB peerlink.PeerLinkRes_LinkConf
	{
		wg.Add(1)
		go addRouteInterfaceToShareNet(client,
			rpctask.ShareNetID,
			rpctask.Req.GetPeerARouterid(),
			availableIP[0],
			&peerA,
			&wg)
	}
	{
		wg.Add(1)
		go addRouteInterfaceToShareNet(client,
			rpctask.ShareNetID,
			rpctask.Req.GetPeerBRouterid(),
			availableIP[1],
			&peerB,
			&wg)
	}

	wg.Wait()
	if "" == peerA.IntfId || "" == peerB.IntfId {
		log.WithFields(log.Fields{
			"peerA": peerA,
			"peerB": peerB,
		}).Error("Create router interface failed")
		return common.EPLCREATEADDINTERFACE
	}

	// 为路由器添加路由表
	{
		wg.Add(1)
		go addRouteToRouter(client,
			subnetBCIDR,
			availableIP[1],
			rpctask.Req.GetPeerARouterid(),
			&peerA,
			&wg)
	}
	{
		wg.Add(1)
		go addRouteToRouter(client,
			subnetACIDR,
			availableIP[0],
			rpctask.Req.GetPeerBRouterid(),
			&peerB,
			&wg)

	}

	wg.Wait()
	if nil == peerA.RouteToPeer || nil == peerB.RouteToPeer {
		log.Error("add route to router failed")
		return common.EPLCREATEADDROUTE
	}

	rpctask.Res.LinkConfOnPeerA = &peerA
	rpctask.Res.LinkConfOnPeerB = &peerB

	return common.EOK
}

func addRouteToRouter(client *gophercloud.ServiceClient,
	cidr string,
	nexthop string,
	routerID string,
	peer *peerlink.PeerLinkRes_LinkConf,
	wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取路由器信息
	router, err := routers.Get(client, routerID).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"cidr":     cidr,
			"nexthop":  nexthop,
			"routerid": routerID,
		}).Error("Get router info failed")
		return
	}

	routes := append(router.Routes, routers.Route{
		DestinationCIDR: cidr,
		NextHop:         nexthop,
	})

	// 更新路由表
	router, err = routers.Update(client, routerID, routers.UpdateOpts{
		Routes: &routes,
	}).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":     err,
			"routeID": routerID,
			"cidr":    cidr,
			"nexthop": nexthop,
			"routes":  routes,
		}).Error("update route to router failed")
		return
	}

	peer.RouteToPeer = &peerlink.PeerLinkRes_LinkConf_Route{
		Destination: cidr,
		Nexthop:     nexthop,
	}
}

func addRouteInterfaceToShareNet(client *gophercloud.ServiceClient,
	shareNetID string,
	routeID string,
	routeInterfaceIP string,
	peer *peerlink.PeerLinkRes_LinkConf,
	wg *sync.WaitGroup) {
	defer wg.Done()

	// 首先必须创建端口
	port := ports.CreateOpts{
		NetworkID: "1ad28c94-3298-4cdb-ac5a-e6d08c133818",
		FixedIPs: []ports.IP{
			{
				SubnetID:  shareNetID,
				IPAddress: routeInterfaceIP,
			},
		},
	}

	pt, err := ports.Create(client, port).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":      err,
			"SubnetID": shareNetID,
			"ipaddr":   routeInterfaceIP,
		}).Error("create ports failed.")
		return
	}

	// 将创建的端口与路由绑定，千万别传入SubnetID，否则会出错
	ifc, err := routers.AddInterface(client, routeID, &routers.AddInterfaceOpts{
		PortID: pt.ID,
	}).Extract()

	if nil != err {
		log.WithFields(log.Fields{
			"err":        err,
			"routeID":    routeID,
			"shareNetID": shareNetID,
			"PortID":     pt.ID,
		}).Error("add interface error")
		return
	}

	peer.CreatedTime = getCurTime()
	peer.IntfId = ifc.PortID
	peer.IntfIp = routeInterfaceIP
}

func getIPFromSubnet(client *gophercloud.ServiceClient, subnetID string, availableIP *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	subnet, err := subnets.Get(client, subnetID).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":      err,
			"subnetID": subnetID,
		}).Error("get subnet by subnetID failed")
		return
	}

	if len(subnet.AllocationPools) == 0 {
		log.WithFields(log.Fields{
			"subnetID": subnetID,
		}).Error("subnet ip pool is empty")
		return
	}

	newPool := make([]subnets.AllocationPool, 0)
	switch subnet.IPVersion {
	case 4:
		for _, pools := range subnet.AllocationPools {
			if len(*availableIP) >= 2 { // 如果申请完了，那么丢到newpool里面用于更新子网ip池子
				newPool = append(newPool, pools)
				continue
			}
			*availableIP = append(*availableIP, pools.Start)
			if pools.Start == pools.End { // 首尾ip相同，
				continue
			}
			startIP := inetaton(pools.Start)
			endIP := inetaton(pools.End)
			nextIP := startIP + 1
			*availableIP = append(*availableIP, inetntoa(nextIP))

			if endIP != nextIP { // 如果池子超过2个可用ip，则后续还可以使用
				newPool = append(newPool, subnets.AllocationPool{
					Start: inetntoa(nextIP + 1),
					End:   pools.End,
				})
			}
		}
		break
	case 6:
		break
	default:
		break
	}

	// 更新子网池
	_, err = subnets.Update(client, subnetID, subnets.UpdateOpts{
		AllocationPools: newPool,
	}).Extract()
	if nil != err {
		*availableIP = nil
	}
}

func (rpctask *CreatePeerLinkRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetPeerARouterid() ||
		"" == rpctask.Req.GetPeerASubnetid() ||
		"" == rpctask.Req.GetPeerBRouterid() ||
		"" == rpctask.Req.GetPeerBSubnetid() ||
		"" == rpctask.ShareNetID {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *CreatePeerLinkRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"res": rpctask.Res,
	}).Info("request end")
}
