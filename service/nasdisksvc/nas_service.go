package nasdisksvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/proto/nasdisk"

	log "github.com/sirupsen/logrus"
	"iaas-api-server/common"
)

type NasDiskService struct {
	nasdisk.UnimplementedNasDiskServiceServer
}

func (o *NasDiskService) CreateNasDisk(ctx context.Context, r *nasdisk.CreateNasDiskReq) (*nasdisk.CreateNasDiskRes, error) {
	log.Info("[NasDiskService] CreateNasDisk request start. ")
	log.Debug("[NasDiskService] CreateNasDisk Apikey:")
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(CreateNasDiskOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*nasdisk.CreateNasDiskRes), err
}

func (o *NasDiskService) DeleteNasDisk(ctx context.Context, r *nasdisk.DeleteNasDiskReq) (*nasdisk.DeleteNasDiskRes, error) {
	log.Info("[NasDiskService] DeleteNasDisk request start. ")
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(DeleteNasDiskOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*nasdisk.DeleteNasDiskRes), err
}

func (o *NasDiskService) GetMountClients(ctx context.Context, r *nasdisk.GetMountClientsReq) (*nasdisk.GetMountClientsRes, error) {
	log.Info("[NasDiskService] GetMountClients request start. ")
	// Todo...
	log.Info("[NasDiskService] GetMountClients Unimplemented.")
	return nil, nil
}

func (o *NasDiskService) Process(auth Authorization, op Op) (interface{}, error) {
	defer log.Info("Request Done.")
	err := op.Predo()

	if err == common.EOK {
		if auth.Auth() == false {
			log.Info("[NasDiskService] Process Authorization Failure.")
			err = common.EUNAUTHORED
		} else {
			err = op.Do()
		}
	} else {
		log.Error("[NasDiskService] Process Predo Failure. ", err)
	}

	return op.Done(err)
}
