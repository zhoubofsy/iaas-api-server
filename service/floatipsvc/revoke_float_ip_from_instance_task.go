package floatipsvc

import (
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	computefloatip "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	netfloatip "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/floatip"
)

type RevokeFloatIpFromInstanceTask struct {
	Req *floatip.RevokeFloatIpFromInstanceReq
	Res *floatip.RevokeFloatIpFromInstanceRes
	Err *common.Error
}

func (rpctask *RevokeFloatIpFromInstanceTask) Run(context.Context) {
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

func (rpctask *RevokeFloatIpFromInstanceTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	netcl, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("get openstack network client failed")
		return &common.Error{
			Code: common.ENETWORKCLIENT.Code,
			Msg:  err.Error(),
		}
	}

	comcl, err := openstack.NewComputeV2(providers, gophercloud.EndpointOpts{})
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("get openstack compute client failed")
		return &common.Error{
			Code: common.ECOMPUTECLIENT.Code,
			Msg:  err.Error(),
		}
	}

	//取消floating ip关联并释放floating ip
	err = computefloatip.DisassociateInstance(comcl, rpctask.Req.InstanceId, computefloatip.DisassociateOpts{
		FloatingIP: rpctask.Req.FloatIp,
	}).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("floating ip disassociate instance failed")
		return &common.Error{
			Code: common.EFLOATINGIPDISASSOCIATE.Code,
			Msg:  err.Error(),
		}
	}

	pages, err := netfloatip.List(netcl, netfloatip.ListOpts{
		FloatingIP: rpctask.Req.FloatIp,
	}).AllPages()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("floating ip list failed")
		return &common.Error{
			Code: common.EFLOATINGIPLIST.Code,
			Msg:  err.Error(),
		}
	}

	allFloatingIps, err := netfloatip.ExtractFloatingIPs(pages)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("floating ip extract pages failed")
		return &common.Error{
			Code: common.EFLOATINGIPEXTRACT.Code,
			Msg:  err.Error(),
		}
	}

	floatingIpID := allFloatingIps[0].ID
	err = netfloatip.Delete(netcl, floatingIpID).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("floating ip delete failed")
		return &common.Error{
			Code: common.EFLOATINGIPDELETE.Code,
			Msg:  err.Error(),
		}
	}
	revokeTime := common.Now()

	rpctask.Res.RevokedTime = revokeTime

	return common.EOK
}

func (rpctask *RevokeFloatIpFromInstanceTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetInstanceId() ||
		"" == rpctask.Req.GetFloatIp() {
		return errors.New("input param is wrong")
	}
	return nil
}

func (rpctask *RevokeFloatIpFromInstanceTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
