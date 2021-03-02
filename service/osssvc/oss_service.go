package osssvc

import (
	"iaas-api-server/common"
	"iaas-api-server/proto/oss"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type OSSService struct {
	oss.UnimplementedOSSServiceServer
}

func (o *OSSService) CreateUserAndBucket(ctx context.Context, r *oss.CreateUserAndBucketReq) (*oss.CreateUserAndBucketRes, error) {
	log.Info("[OSSService] CreateUserAndBucket request start. ", *r)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(CreateUserAndBucketOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.CreateUserAndBucketRes), err
}

func (o *OSSService) GetBucketInfo(ctx context.Context, r *oss.GetBucketInfoReq) (*oss.GetBucketInfoRes, error) {
	log.Info("[OSSService] GetBucketInfo request start. ", *r)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(GetBucketInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.GetBucketInfoRes), err
}

func (o *OSSService) ListBucketsInfo(ctx context.Context, r *oss.ListBucketsInfoReq) (*oss.ListBucketsInfoRes, error) {
	log.Info("[OSSService] ListBucketsInfo request start. ", *r)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(ListBucketsInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.ListBucketsInfoRes), err
}

func (o *OSSService) SetOssUserQuota(ctx context.Context, r *oss.SetOssUserQuotaReq) (*oss.SetOssUserQuotaRes, error) {
	log.Info("[OSSService] SetOssUserQuota request start. ", *r)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(SetOssUserQuotaOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.SetOssUserQuotaRes), err
}

func (o *OSSService) RecoverKey(ctx context.Context, r *oss.RecoverKeyReq) (*oss.RecoverKeyRes, error) {
	log.Info("[OSSService] RecoverKey request start. ", *r)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(RecoverKeyOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.RecoverKeyRes), err
}

func (o *OSSService) Process(auth Authorization, op Op) (interface{}, error) {
	defer log.Info("Request Done.")
	err := op.Predo()

	if err == common.EOK {
		if auth.Auth() == false {
			log.Info("[OSSService] Process Authorization Failure.")
			err = common.EUNAUTHORED
		} else {
			err = op.Do()
		}
	} else {
		log.Error("[OSSService] Process Predo Failure. ", err)
	}

	return op.Done(err)
}