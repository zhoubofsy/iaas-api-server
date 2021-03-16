package nasdisksvc

import (
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
}

func (o *CreateNasDiskOp) Predo() error {
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(nasdisk.CreateNasDiskRes)
	o.conf = GetNasDiskConfigure()

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
	clusterID, userID, err := o.conf.GetGaneshaConfig(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}

	cephMgr := common.CephMgrRestful{Endpoint: endpoint, User: user, Passwd: passwd}
	dirPath := o.Req.PlatformUserid + "_" + o.Req.ShareName
	cephfsPath := rootPath + "/" + dirPath
	pseudoPath := "/" + dirPath

	maxSizeInG, err := config.GetInt("nasdisk_share_quota_max_size_in_g")
	if nil != err {
		log.Error("[NASDISK] CreateNasDiskOp get nasdisk_max_size_in_g failure")
		return common.ENASGETCONFIG
	}
	if o.Req.ShareSizeInG > 0 {
		maxSizeInG = int(o.Req.ShareSizeInG)
	}
	maxSize := maxSizeInG * 1024 * 1024 * 1024
	maxFiles, err := config.GetInt("nasdisk_share_quota_max_files")
	if nil != err {
		log.Error("[NASDISK] CreateNasDiskOp get nasdisk_share_quota_max_files failure")
		return common.ENASGETCONFIG
	}

	var daemons []common.GaneshaDaemonInfo
	var dispatchDaemons []string
	// 1. 获取目录，判断目录是否存在
	dirs, err := cephMgr.ListCephFSDirectory(cephfsid, rootPath)
	for _, dir := range dirs {
		if dir == dirPath {
			err = common.ENASPATHEXISTED
			goto CREATE_FAILED
		}
	}
	// 2. 创建Cephfs目录
	err = cephMgr.MakeCephFSDirectory(cephfsid, cephfsPath)
	if err != common.EOK {
		goto CREATE_FAILED
	}
	CEPHFS_DIR_FLAG = true
	// 3. 设置Cephfs目录的配额
	err = cephMgr.SetCephFSQuotas(cephfsid, cephfsPath, maxSize, maxFiles)
	if err != common.EOK {
		goto CREATE_FAILED
	}
	// 4. 创建NFS-Ganesha Export
	daemons, err = cephMgr.ListGaneshaDaemons()
	if err != common.EOK {
		goto CREATE_FAILED
	}
	for idx, daemon := range daemons {
		log.Debug("daemon[", strconv.Itoa(idx), "]: ", daemon.DaemonID)
		if daemon.Status == 1 {
			dispatchDaemons = append(dispatchDaemons, daemon.DaemonID)
		}
	}
	err = cephMgr.CreateGaneshaExport(clusterID, userID, cephfsPath, pseudoPath, dispatchDaemons)
	if err != common.EOK {
		goto CREATE_FAILED
	}
	return common.EOK

CREATE_FAILED:
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
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(nasdisk.DeleteNasDiskRes)
	o.conf = GetNasDiskConfigure()

	return common.EOK
}

func (o *DeleteNasDiskOp) Do() error {
	/*
		endpoint, user, passwd, err := o.conf.GetMGRConfig(o.Req.Region)
		if err != nil {
			return common.ENASGETCONFIG
		}
		cephfsid, rootPath, err := o.conf.GetCephFSConfig(o.Req.Region)
		if nil != err {
			return common.ENASGETCONFIG
		}
		clusterID, userID, err := o.conf.GetGaneshaConfig(o.Req.Region)
		if err != nil {
			return common.ENASGETCONFIG
		}

		cephMgr := common.CephMgrRestful{Endpoint: endpoint, User: user, Passwd: passwd}
		dirPath := o.Req.PlatformUserid + "_" + o.Req.ShareName
		cephfsPath := rootPath + "/" + dirPath
		pseudoPath := "/" + dirPath
		// 1. 获取Exports, 查找ExportID
		var exportID string
		exports, err := cephMgr.ListGaneshaExport()
		for idx, export := range exports {
			if pseudoPath == export.Pseudo && export.Path == cephfsPath && clusterID == export.ClusterID {
				exportID = strconv.Itoa(export.ExportID)
				break
			}
		}
		// 2. 删除Ganesha Export
		err = cephMgr.DeleteGaneshaExport(clusterID, exportID)
		if err != common.EOK {
			log.Error("[NASDISK] DeleteNasDiskOp delete ganesha export failure.")
		}
		// 3. 删除Cephfs 目录
		err = cephMgr.RemoveCephFSDirectory(cephfsid, cephfsPath)
		if err != common.EOK {
			log.Error("[NASDISK] DeleteNasDiskOp remove cephfs directory failure.")
		}
		return err
	*/
	return common.EOK
}

func (o *DeleteNasDiskOp) Done(e error) (interface{}, error) {
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}
