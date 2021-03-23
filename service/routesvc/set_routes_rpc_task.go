package routesvc

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/route"
)

type SetRoutesRPCTask struct {
	Req *route.SetRoutesReq
	Res *route.SetRoutesRes
	Err *common.Error
}

func (rpctask *SetRoutesRPCTask) Run(context.Context) {
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

func (rpctask *SetRoutesRPCTask) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})
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

	jsTmp := `
{
    "router" : {
      "routes" : {{.routesInfo}}
   }
}`

	routesInfo := "["
	for i := 0; i < len(rpctask.Req.Routes); i++ {
		if i > 0 {
			routesInfo = routesInfo + fmt.Sprintf(",{\"destination\":\"%s\", \"nexthop\":\"%s\"}",
				rpctask.Req.Routes[i].Destination, rpctask.Req.Routes[i].Nexthop)
		} else {
			routesInfo = routesInfo + fmt.Sprintf("{\"destination\":\"%s\", \"nexthop\":\"%s\"}",
				rpctask.Req.Routes[i].Destination, rpctask.Req.Routes[i].Nexthop)
		}
	}
	routesInfo = routesInfo + "]"

	mp := map[string]string{
		"routesInfo": routesInfo,
	}

	jsbody, _ := common.CreateJsonByTmpl(jsTmp, mp)

	//判断操作类型拼接URL
	var url string
	if "add" == rpctask.Req.GetSetType() {
		url = client.ResourceBase + "routers/" + rpctask.Req.GetRouterId() + "/add_extraroutes"
	} else {
		url = client.ResourceBase + "routers/" + rpctask.Req.GetRouterId() + "/remove_extraroutes"
	}

	res, err := common.CallRawAPI(url, "PUT", jsbody, client.TokenID)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers set failed")
		return &common.Error{
			Code: common.EROUTERSET.Code,
			Msg:  err.Error(),
		}
	}
	setTime := common.Now()

	result, err := simplejson.NewJson(res)
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers simplejson failed")
		return &common.Error{
			Code: common.EROUTERSET.Code,
			Msg:  err.Error(),
		}
	}

	//增加原生接口调用后返回错误信息的处理
	resErr, err := result.Get("NeutronError").Get("message").String()
	if nil == err {
		log.WithFields(log.Fields{
			"err": "routers set failed",
			"req": rpctask.Req.String(),
		}).Error("message: ", resErr)
		return &common.Error{
			Code: common.EROUTERSET.Code,
			Msg:  resErr,
		}
	}

	resId, err := result.Get("router").Get("id").String()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers get id failed")
		return &common.Error{
			Code: common.EROUTERSET.Code,
			Msg:  err.Error(),
		}
	}

	routes, err := result.Get("router").Get("routes").Array()
	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("routers get routes failed")
		return &common.Error{
			Code: common.EROUTERSET.Code,
			Msg:  err.Error(),
		}
	}

	resRoutes := make([]*route.Route, 0)
	for _, row := range routes {
		if eachMap, ok := row.(map[string]interface{}); ok {
			if dn, ok := eachMap["destination"].(string); ok {
				if nh, ok := eachMap["nexthop"].(string); ok {
					resRoutes = append(resRoutes, &route.Route{
						Destination: dn,
						Nexthop:     nh,
					})
				}
			}
		}
	}

	rpctask.Res = &route.SetRoutesRes{
		RouterId:      resId,
		SetType:       rpctask.Req.GetSetType(),
		SetTime:       setTime,
		CurrentRoutes: resRoutes,
	}

	return common.EOK
}

func (rpctask *SetRoutesRPCTask) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() ||
		"" == rpctask.Req.GetRouterId() ||
		"" == rpctask.Req.GetSetType() ||
		nil == rpctask.Req.GetRoutes() {
		return errors.New("input param is wrong")
	}
	return nil
}

func (rpctask *SetRoutesRPCTask) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
