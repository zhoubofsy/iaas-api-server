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
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
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
	var sharePoolAIP, sharePoolBIP string

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
	}

	wg.Wait()
	if "" == subnetACIDR || "" == subnetBCIDR {
		log.Error("get cidr by subnetid failed")
		return common.EPLGETCIDR
	}

	// TODO 异步创建interface（route -> share subnet）
	{
		wg.Add(1)
	}
	{
		wg.Add(1)
	}

	wg.Wait()

	return common.EOK
}

func getCIDRBySubnetID(client *gophercloud.ServiceClient, subnetID string, cidr *string, wg *sync.WaitGroup) {
	defer wg.Done()

	subnet, err := subnets.Get(client, subnetID).Extract()
	if nil != err {
		*cidr = ""
	} else {
		*cidr = subnet.CIDR
	}
}

func addRouteInterfaceToShareNet(client *gophercloud.ServiceClient, shareNetID string, routeID string, routeInterfaceIP *string, wg *sync.WaitGroup) {
	defer wg.Done()

	ifc, err := routers.AddInterface(client, routeID, &routers.AddInterfaceOpts{
		SubnetID: shareNetID,
	}).Extract()

	if nil != err {
		*routeInterfaceIP = ""
		return
	}
}

func getIPFromSubnetPool(client *gophercloud.ServiceClient, shareNetID string, wg *sync.WaitGroup) {
	defer wg.Done()

	subnetPool, err := subnetpools.Get(client, shareNetID).Extract()
	if nil != err {
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
}
