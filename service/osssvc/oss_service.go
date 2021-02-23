package osssvc

import (
	"golang.org/x/net/context"
	"iaas-api-server/common"
	"iaas-api-server/proto/oss"
)

type OSSService struct {
	oss.UnimplementedOSSServiceServer
}

func (o *OSSService) CreateUserAndBucket(ctx context.Context, r *oss.CreateUserAndBucketReq) (*oss.CreateUserAndBucketRes, error) {
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(CreateUserAndBucketOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.CreateUserAndBucketRes), err
}

func (o *OSSService) GetBucketInfo(ctx context.Context, r *oss.GetBucketInfoReq) (*oss.GetBucketInfoRes, error) {
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(GetBucketInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.GetBucketInfoRes), err
}

func (o *OSSService) ListBucketsInfo(ctx context.Context, r *oss.ListBucketsInfoReq) (*oss.ListBucketsInfoRes, error) {
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(ListBucketsInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.ListBucketsInfoRes), err
}

func (o *OSSService) SetOssUserQuota(ctx context.Context, r *oss.SetOssUserQuotaReq) (*oss.SetOssUserQuotaRes, error) {
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(SetOssUserQuotaOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.SetOssUserQuotaRes), err
}

func (o *OSSService) RecoverKey(ctx context.Context, r *oss.RecoverKeyReq) (*oss.RecoverKeyRes, error) {
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(RecoverKeyOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.RecoverKeyRes), err
}

func (o *OSSService) Process(auth Authorization, op Op) (interface{}, error) {
	err := op.Predo()

	if err == common.EOK {
		if auth.Auth() == false {
			err = common.EUNAUTHORED
		} else {
			err = op.Do()
		}
	}

	return op.Done(err)
}
