package nasdisksvc

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	netfloatip "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	log "github.com/sirupsen/logrus"
	"iaas-api-server/common"
	"iaas-api-server/common/config"
	"iaas-api-server/proto/nasdisk"
	"strconv"
)

type Authorization interface {
	Auth() bool
}

type OpenstackAPIAuthorization struct {
	Apikey         string
	TenantId       string
	PlatformUserid string
}

func (o *OpenstackAPIAuthorization) Auth() bool {
	return common.APIAuth(o.Apikey, o.TenantId, o.PlatformUserid)
}

type Op interface {
	// Predo
	Predo() error
	// Do
	Do() error
	// Done
	Done(error) (interface{}, error)
}

type BasicOp struct {
	conf NasDiskConfigure
}

type CreateNasDiskOp struct {
	BasicOp
	Req *nasdisk.CreateNasDiskReq
	Res *nasdisk.CreateNasDiskRes
	Ops *gophercloud.ProviderClient
}

func (o *CreateNasDiskOp) getAllowIPByNetworkID(networkID string) ([]string, error) {
	allowIP := []string{}
	client, err := openstack.NewNetworkV2(o.Ops, gophercloud.EndpointOpts{})
	if err != nil {
		return allowIP, common.ENETWORKCLIENT
	}
	networkInfo, err := networks.Get(client, networkID).Extract()
	if err != nil {
		return allowIP, common.ENETWORKSGET
	}
	routerName := "router-" + networkInfo.Name
	routerPages, err := routers.List(client, routers.ListOpts{Name: routerName}).AllPages()
	if err != nil {
		return allowIP, common.EROUTERLIST
	}
	routersInfo, err := routers.ExtractRouters(routerPages)
	if err != nil {
		return allowIP, common.EROUTEREXTRACT
	}
	if 1 != len(routersInfo) {
		return allowIP, common.EROUTERINFO
	}
	if 0 >= len(routersInfo[0].GatewayInfo.ExternalFixedIPs) {
		return allowIP, common.EROUTERINFO
	}
	pages, err := netfloatip.List(client, netfloatip.ListOpts{
		RouterID: routersInfo[0].ID,
	}).AllPages()
	if nil != err {
		return allowIP, common.EFLOATINGIPLIST
	}
	allFloatingIps, err := netfloatip.ExtractFloatingIPs(pages)
	if nil != err {
		return allowIP, common.EFLOATINGIPEXTRACT
	}
	allowIP = append(allowIP, routersInfo[0].GatewayInfo.ExternalFixedIPs[0].IPAddress)
	for _, ip := range allFloatingIps {
		allowIP = append(allowIP, ip.FloatingIP)
	}
	return allowIP, common.EOK
}

func (o *CreateNasDiskOp) Predo() error {
	if o.Req == nil || o.Req.ShareSizeInG <= 0 || o.Req.Region == "" || o.Req.PlatformUserid == "" || o.Req.ShareName == "" || o.Req.TenantId == "" || o.Req.Apikey == "" {
		return common.EPARAM
	}
	o.Res = new(nasdisk.CreateNasDiskRes)
	o.conf = GetNasDiskConfigure()
	ops, err := common.GetOpenstackClient(o.Req.Apikey, o.Req.TenantId, o.Req.PlatformUserid)
	if err != nil {
		return common.EGETOPSTACKCLIENT
	}
	o.Ops = ops

	return common.EOK
}

func (o *CreateNasDiskOp) Do() error {
	var (
		CEPHFS_DIR_FLAG = false
	)
	endpoint, user, passwd, err := o.conf.GetMGRConfig(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}
	cephfsid, rootPath, err := o.conf.GetCephFSConfig(o.Req.Region)
	if nil != err {
		return common.ENASGETCONFIG
	}
	nfsDomain, clusterID, userID, err := o.conf.GetGaneshaConfig(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}

	cephMgr := common.CephMgrRestful{Endpoint: endpoint, User: user, Passwd: passwd}
	dirPath := o.Req.PlatformUserid + "_" + o.Req.ShareName
	cephfsPath := rootPath + "/" + dirPath
	pseudoPath := "/" + dirPath

	maxSize := int(o.Req.ShareSizeInG) * 1024 * 1024 * 1024
	maxFiles, err := config.GetInt("nasdisk_share_quota_max_files")
	if nil != err {
		log.Error("[NASDISK] CreateNasDiskOp get nasdisk_share_quota_max_files failure")
		return common.ENASGETCONFIG
	}

	o.Res.NasDisk = &(nasdisk.CreateNasDiskRes_NasDisk{
		ShareId:      dirPath,
		ShareName:    o.Req.ShareName,
		ShareDesc:    o.Req.ShareDesc,
		ShareProto:   o.Req.ShareProto,
		ShareSizeInG: o.Req.ShareSizeInG,
		Region:       o.Req.Region,
		NetworkId:    o.Req.NetworkId,
		MountPoint:   nfsDomain + ":" + pseudoPath,
		CreatedTime:  common.Now()})
	var daemons []common.GaneshaDaemonInfo
	var dispatchDaemons []string
	// 1. 获取目录，判断目录是否存在
	dirs, err := cephMgr.ListCephFSDirectory(cephfsid, rootPath)
	for _, dir := range dirs {
		if dir == dirPath {
			err = common.ENASPATHEXISTED
			return err
		}
	}
	ips, err := o.getAllowIPByNetworkID(o.Req.NetworkId)
	if err != common.EOK {
		return err
	}
	// 2. 创建Cephfs目录
	err = cephMgr.MakeCephFSDirectory(cephfsid, cephfsPath)
	if err != common.EOK {
		return err
	}
	CEPHFS_DIR_FLAG = true
	for true {
		// 3. 设置Cephfs目录的配额
		err = cephMgr.SetCephFSQuotas(cephfsid, cephfsPath, maxSize, maxFiles)
		if err != common.EOK {
			break
		}
		// 4. 创建NFS-Ganesha Export
		daemons, err = cephMgr.ListGaneshaDaemons()
		if err != common.EOK {
			break
		}
		for idx, daemon := range daemons {
			log.Debug("daemon[", strconv.Itoa(idx), "]: ", daemon.DaemonID)
			if daemon.Status == 1 {
				dispatchDaemons = append(dispatchDaemons, daemon.DaemonID)
			}
		}
		err = cephMgr.CreateGaneshaExport(clusterID, userID, cephfsPath, pseudoPath, dispatchDaemons, ips)
		if err != common.EOK {
			break
		}
		return common.EOK
	}

	if CEPHFS_DIR_FLAG {
		// 删除 Cephfs 目录
		cephMgr.RemoveCephFSDirectory(cephfsid, cephfsPath)
	}
	return err
}

func (o *CreateNasDiskOp) Done(e error) (interface{}, error) {
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}

type DeleteNasDiskOp struct {
	BasicOp
	Req *nasdisk.DeleteNasDiskReq
	Res *nasdisk.DeleteNasDiskRes
}

func (o *DeleteNasDiskOp) Predo() error {
	if o.Req == nil || o.Req.ShareId == "" || o.Req.PlatformUserid == "" || o.Req.TenantId == "" || o.Req.Apikey == "" {
		return common.EPARAM
	}
	o.Res = new(nasdisk.DeleteNasDiskRes)
	o.conf = GetNasDiskConfigure()

	return common.EOK
}

func (o *DeleteNasDiskOp) Do() error {
	endpoint, user, passwd, err := o.conf.GetMGRConfig(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}
	cephfsid, rootPath, err := o.conf.GetCephFSConfig(o.Req.Region)
	if nil != err {
		return common.ENASGETCONFIG
	}
	_, clusterID, _, err := o.conf.GetGaneshaConfig(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}

	cephMgr := common.CephMgrRestful{Endpoint: endpoint, User: user, Passwd: passwd}
	//dirPath := o.Req.PlatformUserid + "_" + o.Req.ShareName
	dirPath := o.Req.ShareId
	cephfsPath := rootPath + "/" + dirPath
	pseudoPath := "/" + dirPath

	o.Res.ShareId = o.Req.ShareId
	o.Res.DeletedTime = common.Now()
	// 1. 获取Exports, 查找ExportID
	var exportID string
	exports, err := cephMgr.ListGaneshaExport()
	for _, export := range exports {
		if pseudoPath == export.Pseudo && export.Path == cephfsPath && clusterID == export.ClusterID {
			exportID = strconv.Itoa(export.ExportID)
			break
		}
	}
	if exportID != "" {
		// 2. 删除Ganesha Export
		err = cephMgr.DeleteGaneshaExport(clusterID, exportID)
		if err != common.EOK {
			log.Error("[NASDISK] DeleteNasDiskOp delete ganesha export failure.")
		}
	}
	// 3. 删除Cephfs 目录
	err = cephMgr.RemoveCephFSDirectory(cephfsid, cephfsPath)
	if err != common.EOK {
		log.Error("[NASDISK] DeleteNasDiskOp remove cephfs directory failure.")
	}
	return err
}

func (o *DeleteNasDiskOp) Done(e error) (interface{}, error) {
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}

func getGatewayByNetworkID(apiKey string, tenantID string, platformUserid string, networkID string) (string, error) {
	ops, err := common.GetOpenstackClient(apiKey, tenantID, platformUserid)
	if err != nil {
		return "", common.EGETOPSTACKCLIENT
	}
	client, err := openstack.NewNetworkV2(ops, gophercloud.EndpointOpts{})
	if err != nil {
		return "", common.ENETWORKCLIENT
	}
	networkInfo, err := networks.Get(client, networkID).Extract()
	if err != nil {
		return "", common.ENETWORKSGET
	}
	routerName := "router-" + networkInfo.Name
	routerPages, err := routers.List(client, routers.ListOpts{Name: routerName}).AllPages()
	if err != nil {
		return "", common.EROUTERLIST
	}
	routersInfo, err := routers.ExtractRouters(routerPages)
	if err != nil {
		return "", common.EROUTEREXTRACT
	}
	if 1 != len(routersInfo) {
		return "", common.EROUTERINFO
	}
	if 0 >= len(routersInfo[0].GatewayInfo.ExternalFixedIPs) {
		return "", common.EROUTERINFO
	}
	return routersInfo[0].GatewayInfo.ExternalFixedIPs[0].IPAddress, common.EOK
}

func UpdateGaneshaExportClient(addition bool, apiKey string, tenantID string, platformUserid string, networkID string, floatingIP string, selectRegion ...string) error {
	// 0. prepare CephMgrRest
	conf := GetNasDiskConfigure()
	region := "RegionOne"
	if len(selectRegion) == 1 {
		region = selectRegion[0]
	}
	endpoint, user, passwd, err := conf.GetMGRConfig(region)
	if err != nil {
		return common.ENASGETCONFIG
	}
	_, clusterID, _, err := conf.GetGaneshaConfig(region)
	if err != nil {
		return common.ENASGETCONFIG
	}

	// 1. get gateway by networkID
	gateway, err := getGatewayByNetworkID(apiKey, tenantID, platformUserid, networkID)
	if err != common.EOK {
		return err
	}
	cephMgr := common.CephMgrRestful{Endpoint: endpoint, User: user, Passwd: passwd}
	// 2. list all exports
	exports, err := cephMgr.ListGaneshaExport()
	if err != common.EOK {
		return err
	}

	for ie, export := range exports {
		for ic, clt := range export.Clients {
			found := false
			iExport := 0
			iClient := 0
			// 3. find export by networkID
			for _, addr := range clt.Addresses {
				if addr == gateway {
					found = true
					iExport = ie
					iClient = ic
					break
				}
			}
			if found {
				existed := false
				idx := 0
				fip := floatingIP
				for i, addr := range exports[iExport].Clients[iClient].Addresses {
					if addr == fip {
						existed = true
						idx = i
						break
					}
				}
				update := false
				if addition {
					if !existed {
						// append floating ip
						exports[iExport].Clients[iClient].Addresses = append(exports[iExport].Clients[iClient].Addresses, fip)
						update = true
					}
				} else {
					if existed {
						// delete floating ip
						exports[iExport].Clients[iClient].Addresses = append(exports[iExport].Clients[iClient].Addresses[:idx], exports[iExport].Clients[iClient].Addresses[idx+1:]...)
						update = true
					}
				}
				// 4. update export
				if update {
					err := cephMgr.PutGaneshaExport(clusterID, strconv.Itoa(exports[iExport].ExportID), exports[iExport])
					if err != common.EOK {
						return err
					}
				}
			}
		}
	}

	return common.EOK
}
