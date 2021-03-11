package nasdisksvc

import (
	"iaas-api-server/common"
	"iaas-api-server/proto/nasdisk"
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
		CEPHFS_DIR_FLAG  = false
		GANESHA_EXP_FLAG = false
	)
	endpoint, err := o.conf.GetMGRRestfulAddr(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}
	user, passwd, err := o.conf.GetMGRUserPasswd(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}
	cephfsid, err := o.conf.GetCephfsID(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}
	rootPath, err := o.conf.GetRootPath(o.Req.Region)
	if err != nil {
		return common.ENASGETCONFIG
	}

	cephMgr := common.CephMgrRestful{Endpoint: endpoint, User: user, Passwd: passwd}
	dirPath := o.Req.PlatformUserid + o.Req.ShareId
	cephfsPath := rootPath + "/" + o.Req.PlatformUserid
	maxSize := 0
	maxFiles := 0
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
	GANESHA_EXP_FLAG = true
	return common.EOK

CREATE_FAILED:
	if GANESHA_EXP_FLAG {
		// 删除 NFS-Ganesha Export
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
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(nasdisk.DeleteNasDiskRes)
	o.conf = GetNasDiskConfigure()

	return common.EOK
}

func (o *DeleteNasDiskOp) Do() error {
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
