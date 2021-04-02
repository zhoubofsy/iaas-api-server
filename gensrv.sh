#!/bin/bash

#=================================================================
#  
#  文件名称：gensrv.sh
#  创 建 者: mongia
#  创建日期：2021年04月02日
#  
#=================================================================

set -e -x

if [ $# != 7 ]; then
    echo ""
    echo "usage: ./gensvr.sh RpcTaskFilename ModuleName ServiceName MethodName ReqName ResName AuthorName"
    echo ""
    echo "example: ./gensrv.sh create_firewall_rpc_task.go firewall FirewallService CreateFirewall CreateFirewallReq FirewallRes mongia"
    echo ""
    exit 0
fi

filename=$1
modulename=$2
servicename=$3
taskname=$4RPCTask
methodname=$4
req=$5
res=$6
author=$7
time=`date +%Y年%m月%d日`

cat > $filename << EOF
/*================================================================
*
*  文件名称：${filename}
*  创 建 者: ${author}
*  创建日期：${time}
*
================================================================*/

package ${modulename}svc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/${modulename}"
)

type ${taskname} struct {
	Req *${modulename}.${req}
	Res *${modulename}.${res}
	Err *common.Error
}

// Run call this func for doing task
func (rpctask *${taskname}) Run(context.Context) {
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

func (rpctask *${taskname}) execute(providers *gophercloud.ProviderClient) *common.Error {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{})

	if nil != err {
		log.WithFields(log.Fields{
			"err": err,
			"req": rpctask.Req.String(),
		}).Error("new network v2 failed.")
		return common.ENETWORKCLIENT
	}

    //TODO code here
    log.WithFields(log.Fields{
		"client": client,
	}).Info("remove this code")

	return common.EOK
}

func (rpctask *${taskname}) checkParam() error {
	if "" == rpctask.Req.GetApikey() ||
		"" == rpctask.Req.GetTenantId() ||
		"" == rpctask.Req.GetPlatformUserid() {
		return errors.New("input params is wrong")
	}
	return nil
}

func (rpctask *${taskname}) setResult() {
	rpctask.Res.Code = rpctask.Err.Code
	rpctask.Res.Msg = rpctask.Err.Msg

	log.WithFields(log.Fields{
		"req": rpctask.Req,
		"res": rpctask.Res,
	}).Info("request end")
}
EOF

mv ${filename} service/${modulename}svc

if [ ! -e "service/${modulename}svc/${modulename}_service.go" ]; then
cat > service/${modulename}svc/${modulename}_service.go << EOF
/*================================================================
*
*  文件名称：${modulename}_service.go
*  创 建 者: ${author}
*  创建日期：${time}
*
================================================================*/

package ${modulename}svc

import (
	"golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"iaas-api-server/proto/${modulename}"
)

type ${servicename} struct {
    ${modulename}.Unimplemented${servicename}Server
}
EOF
fi

cat >> service/${modulename}svc/${modulename}_service.go << EOF

func (pthis *${servicename}) ${methodname}(ctx context.Context, req *${modulename}.${req}) (*${modulename}.${res}, error) {
    task := ${taskname}{
        Req: req,
        Res: &${modulename}.${res}{},
        Err: nil,
    }

    task.Run(ctx)

    return task.Res, status.Error(codes.OK, "success")
}
EOF
