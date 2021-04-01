package floatipsvc

import (
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	computefloatip "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	netfloatip "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/floatip"
	nd "iaas-api-server/service/nasdisksvc"
)

type BindFloatIpToInstanceTask struct {
	Req            *floatip.BindFloatIpToInstanceReq
	Res            *floatip.BindFloatIpToInstanceRes
	Err            *common.Error
	PublicNetID    string
	PublicSubnetID string
}

func (rpctask *BindFloatIpToInstanceTask) Run(context.Context) {
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

func (rpctask *BindFloatIpToInstanceTask) execute(providers *gophercloud.ProviderClient) *common.Error {
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

	//查找路由器上是否存在外部网关，不存在则返回错误信息
	router, err := routers.Get(netcl, rpctask.Req.VpcRouterId).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers get failed")
		return &common.Error{
			Code: common.EROUTERGET.Code,
			Msg:  err.Error(),
		}
	}

	if "" == router.GatewayInfo.NetworkID {
		log.WithFields(log.Fields{
			"err": "The router has no external gateway",
			"req": rpctask.Req.String(),
		}).Error("The router has no external gateway")
		return &common.Error{
			Code: common.EROUTERNOGATWAY.Code,
			Msg:  "The router has no external gateway",
		}
	}

	//生成floating ip并关联vm
	floatingIp, err := netfloatip.Create(netcl, netfloatip.CreateOpts{
		FloatingNetworkID: rpctask.PublicNetID,
		SubnetID:          rpctask.PublicSubnetID,
	}).Extract()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("floating ip create failed")
		return &common.Error{
			Code: common.EFLOATINGIPCREATE.Code,
			Msg:  err.Error(),
		}
	}

	err = computefloatip.AssociateInstance(comcl, rpctask.Req.InstanceId, computefloatip.AssociateOpts{
		FloatingIP: floatingIp.FloatingIP,
	}).ExtractErr()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("floating ip associate instance failed")

		rollBackErr := netfloatip.Delete(netcl, floatingIp.ID).ExtractErr()
		if nil != rollBackErr {
			log.Info("rollback delete float ip err: ", rollBackErr)
		}

		return &common.Error{
			Code: common.EFLOATINGIPASSOCIATE.Code,
			Msg:  err.Error(),
		}
	}

	//通知Ganesha更新Export,增加浮动ip
	go nd.UpdateGaneshaExportClient(true, rpctask.Req.Apikey, rpctask.Req.TenantId,
		rpctask.Req.PlatformUserid, router.GatewayInfo.ExternalFixedIPs[0].IPAddress, floatingIp.FloatingIP)

	rpctask.Res.FloatIp = floatingIp.FloatingIP
	rpctask.Res.BindedTime = floatingIp.CreatedAt.Local().Format("2006-01-02 15:04:05")

	return common.EOK
}

func (rpctask *BindFloatIpToInstanceTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetInstanceId() ||
		"" == rpctask.Req.GetVpcRouterId() {
		return errors.New("input param is wrong")
	}
	return nil
}

func (rpctask *BindFloatIpToInstanceTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
