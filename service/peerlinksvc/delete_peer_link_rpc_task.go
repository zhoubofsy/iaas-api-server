/*================================================================
*
*  文件名称：delete_peer_link_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年02月19日
*
================================================================*/

package peerlinksvc

import (
	"errors"
	"sort"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/peerlink"
)

// DeletePeerLinkRPCTask rpc task
type DeletePeerLinkRPCTask struct {
	Req        *peerlink.PeerLinkReq
	Res        *peerlink.DeletePeerLinkRes
	Err        *common.Error
	ShareNetID string
}

// Run rpc start
func (rpctask *DeletePeerLinkRPCTask) Run(context.Context) {
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

func (rpctask *DeletePeerLinkRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
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

	// TODO 以下接口都可能出现失败，聚合操作保证事务性，后续得修改
	// 异步获取peera和peerb的子网cidr，后面设置需要用到, 从share net pool获取2个ip
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

	wg.Wait()
	if "" == subnetACIDR || "" == subnetBCIDR {
		log.Error("get cidr by subnetid failed")
		return common.EPLDELETEPREPARE
	}

	// TODO 从路由器里面删除接口跟路由表, 获取对应的router ip，把ip还给sharenet子网池
	var routerIP = []int64{0, 0} // 用来存放ip的int值，方便后续快速归还ip
	{
		wg.Add(1)
		go removeRouteFromRouter(client,
			subnetACIDR,
			rpctask.Req.GetPeerBRouterid(),
			rpctask.ShareNetID,
			&routerIP[0],
			&wg)
	}
	{
		wg.Add(1)
		go removeRouteFromRouter(client,
			subnetBCIDR,
			rpctask.Req.GetPeerARouterid(),
			rpctask.ShareNetID,
			&routerIP[1],
			&wg)
	}

	wg.Wait()

	// TODO 上面的逻辑没有拿到router从sharenet获取的ip，那么从router的interface获取ip
	if 0 == routerIP[0] {
		var portsA ports.Port
		wg.Add(1)
		go getPortByRouterIDAndNetID(client,
			rpctask.Req.GetPeerARouterid(),
			rpctask.ShareNetID,
			&portsA,
			&wg)
		if 0 != len(portsA.FixedIPs) {
			routerIP[0] = inetaton(portsA.FixedIPs[0].IPAddress)
		}
	}

	if 0 == routerIP[1] {
		var portsB ports.Port
		wg.Add(1)
		go getPortByRouterIDAndNetID(client,
			rpctask.Req.GetPeerBRouterid(),
			rpctask.ShareNetID,
			&portsB,
			&wg)
		if 0 != len(portsB.FixedIPs) {
			routerIP[1] = inetaton(portsB.FixedIPs[0].IPAddress)
		}
	}

	wg.Wait()

	// 归还ip给子网池
	sort.Slice(routerIP, func(i, j int) bool { return routerIP[i] < routerIP[j] })
	{
		wg.Add(1)
		go giveBackIPToSubnet(client, rpctask.ShareNetID, routerIP, len(routerIP), &wg)
	}

	wg.Wait()

	return common.EOK
}

func giveBackIPToSubnet(client *gophercloud.ServiceClient,
	subnetID string,
	ip []int64,
	length int,
	wg *sync.WaitGroup) {
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
		}).Info("subnet ip pool is empty")
	}

	// 构造归还的子网ip pool，减少后续循环次数
	newPool := make([]subnets.AllocationPool, 0)
	for i := 0; i < length; i++ {
		switch subnet.IPVersion {
		case 4:
			for _, pools := range subnet.AllocationPools { // 遍历现在的子网池
				start := inetaton(pools.Start)
				end := inetaton(pools.End)

				if (ip[i] + 1) == start { // 归还的ip在当前子网池的头部
					newPool = append(newPool, subnets.AllocationPool{
						Start: inetntoa(ip[i]),
						End:   pools.End,
					})
				} else if (end + 1) == ip[i] { //  归还的ip在当前子网池的尾部
					newPool = append(newPool, subnets.AllocationPool{
						Start: pools.Start,
						End:   inetntoa(ip[i]),
					})
				} else { // 当前归还ip不在头尾
					newPool = append(newPool, pools)
				}
			}
			break
		case 6: // 暂时没考虑ipv6
			log.Error("not surpport ipv6")
			break
		default:
			log.Error("ipversion is not 4 or 6")
			break
		}
	}
}

func removeRouteFromRouter(client *gophercloud.ServiceClient,
	cidr string,
	routerID string,
	shareNetID string,
	routerIP *int64,
	wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取路由器信息
	router, err := routers.Get(client, routerID).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"cidr":     cidr,
			"routerid": routerID,
		}).Error("Get router info failed")
		return
	}

	routes := make([]routers.Route, 0)
	for _, route := range router.Routes {
		if route.DestinationCIDR != cidr {
			routes = append(routes, route)
		} else {
			*routerIP = inetaton(route.NextHop)
		}
	}

	// 更新路由表，把指定的对端路由表删除掉
	router, err = routers.Update(client, routerID, routers.UpdateOpts{
		Routes: &routes,
	}).Extract()

	if nil != err {
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"cidr":     cidr,
			"routerid": routerID,
		}).Error("detele route from router failed")
		return
	}

	// 删除跟share网络的接口
	_, err = routers.RemoveInterface(client, routerID, routers.RemoveInterfaceOpts{
		SubnetID: shareNetID,
	}).Extract()

	if nil != err {
		log.WithFields(log.Fields{
			"err":      err,
			"cidr":     cidr,
			"routerid": routerID,
			"SubnetID": shareNetID,
		}).Error("remove interface from router failed")
	}
}

func (rpctask *DeletePeerLinkRPCTask) checkParam() error {
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

func (rpctask *DeletePeerLinkRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg
	rpctask.Res.PeerARouterid = rpctask.Req.PeerARouterid
	rpctask.Res.PeerASubnetid = rpctask.Req.PeerASubnetid
	rpctask.Res.PeerBRouterid = rpctask.Req.PeerBRouterid
	rpctask.Res.PeerBSubnetid = rpctask.Req.PeerBSubnetid
	rpctask.Res.DeletedTime = common.Now()

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
