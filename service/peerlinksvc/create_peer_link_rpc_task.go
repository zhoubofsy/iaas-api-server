/*================================================================
*
*  文件名称：create_peer_link_rpc_task.go
*  创 建 者: mongia
*  创建日期：2021年02月19日
*
================================================================*/

package peerlinksvc

import (
	"encoding/json"
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
			"subnetACIDR": subnetACIDR,
			"subnetBCIDR": subnetBCIDR,
			"ipaddr":      availableIP,
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
		NetworkID: "5f55f08b-d55a-45a2-8f3f-c5b3eb295856",
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

	peer.CreatedTime = common.Now()
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

	ipPool, errSQL := common.QuerySharedSubnetUsedIP(subnetID)
	if nil != errSQL && common.EPLGETIPPOOLNONE != errSQL {
		log.WithFields(log.Fields{
			"subnetID": subnetID,
		}).Error("get used ip from mysql failed")
		return
	}
	var usedIP []int64
	var newIP []int64 = make([]int64, 0)

	if common.EPLGETIPPOOLNONE != errSQL {
		err = json.Unmarshal([]byte(ipPool.UsedIP), &usedIP)
		if nil != err {
			log.WithFields(log.Fields{
				"subnetID": subnetID,
			}).Error("parse used ip from json to vector failed")
			return
		}
	}

	// 暂时不考虑子网的ip池被人修改的情况
	switch subnet.IPVersion {
	case 4:
		pools := subnet.AllocationPools[0]
		startIP := inetaton(pools.Start)
		endIP := inetaton(pools.End)
		for i := startIP; i <= endIP; i++ {
			if len(*availableIP) >= 2 {
				break
			}

			// 从已使用的ip池查询当前ip是否已经使用了，这个位置ip应该不会太多，直接遍历
			used := false
			for _, ip := range usedIP {
				if ip == i { // 当前ip已经被使用了，则需要下一跳的ip
					used = true
					break
				}
			}
			if used { // 当前ip已经被使用了，则需要下一跳的ip
				continue
			}
			*availableIP = append(*availableIP, inetntoa(i)) // 分配ip
			newIP = append(newIP, i)
		}
		break
	case 6: // TODO 暂时没考虑ipv6
		log.Error("not surpport ipv6")
		break
	default:
		log.Error("ipversion is not 4 or 6")
		break
	}

	// TODO ip 可以用一个长的0101010的串去标识哪个ip已经使用
	// 将已用的ip记录到数据库
	newUsedIP, err := json.Marshal(append(newIP, usedIP...))
	if nil != err {
		log.WithFields(log.Fields{
			"err":       err,
			"newUsedIP": newIP,
		}).Error("parse ip vector to string failed")
		*availableIP = nil
		return
	}

	// 不存在数据，则插入，否则更新
	var ret bool = false
	if errSQL == common.EPLGETIPPOOLNONE {
		ret = common.CreateSharedSubnetUsedIP(subnetID, string(newUsedIP))
	} else {
		ret = common.UpdateSharedSubnetUsedIP(subnetID, string(newUsedIP))
	}
	if !ret {
		log.WithFields(log.Fields{
			"err":       err,
			"newUsedIP": newIP,
		}).Error("update used ip to DB failed")
		*availableIP = nil
		return
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
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
