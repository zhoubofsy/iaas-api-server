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
	log.Info("[OSSService] CreateUserAndBucket request start. ")
	log.Debug("[OSSService] CreateUserAndBucketReq Apikey:", r.Apikey, ",TenantId:", r.TenantId, ",PlatformUserid:", r.PlatformUserid, ",Region:", r.Region, ",BucketName:", r.BucketName, ",StorageType:", r.StorageType, ",UserMaxSizeInG:", r.UserMaxSizeInG, ",UserMaxObjects:", r.UserMaxObjects, ",BucketPolicy:", r.BucketPolicy)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(CreateUserAndBucketOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.CreateUserAndBucketRes), err
}

func (o *OSSService) GetBucketInfo(ctx context.Context, r *oss.GetBucketInfoReq) (*oss.GetBucketInfoRes, error) {
	log.Info("[OSSService] GetBucketInfo request start. ")
	log.Debug("[OSSService] GetBucketInfoReq Apikey:", r.Apikey, ",TenantId:", r.TenantId, ",PlatformUserid:", r.PlatformUserid, ",Region:", r.Region, ",OssUid:", r.OssUid, ",BucketName:", r.BucketName)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(GetBucketInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.GetBucketInfoRes), err
}

func (o *OSSService) ListBucketsInfo(ctx context.Context, r *oss.ListBucketsInfoReq) (*oss.ListBucketsInfoRes, error) {
	log.Info("[OSSService] ListBucketsInfo request start. ")
	log.Debug("[OSSService] ListBucketsInfoReq Apikey:", r.Apikey, ",TenantId:", r.TenantId, ",PlatformUserid:", r.PlatformUserid, ",Region:", r.Region, ",OssUid:", r.OssUid, ",PageNumber:", r.PageNumber, ",PageSize:", r.PageSize)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(ListBucketsInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.ListBucketsInfoRes), err
}

func (o *OSSService) SetOssUserQuota(ctx context.Context, r *oss.SetOssUserQuotaReq) (*oss.SetOssUserQuotaRes, error) {
	log.Info("[OSSService] SetOssUserQuota request start. ")
	log.Debug("[OSSService] SetOssUserQuotaReq Apikey:", r.Apikey, ",TenantId:", r.TenantId, ",PlatformUserid:", r.PlatformUserid, ",Region:", r.Region, ",OssUid:", r.OssUid, ",UserMaxSizeInG:", r.UserMaxSizeInG, ",UserMaxObjects:", r.UserMaxObjects)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(SetOssUserQuotaOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.SetOssUserQuotaRes), err
}

func (o *OSSService) RecoverKey(ctx context.Context, r *oss.RecoverKeyReq) (*oss.RecoverKeyRes, error) {
	log.Info("[OSSService] RecoverKey request start. ")
	log.Debug("[OSSService] RecoverKeyReq Apikey:", r.Apikey, ",TenantId:", r.TenantId, ",PlatformUserid:", r.PlatformUserid, ",Region:", r.Region, ",OssUid:", r.OssUid)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(RecoverKeyOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.RecoverKeyRes), err
}

func (o *OSSService) GetUserInfo(ctx context.Context, r *oss.GetUserInfoReq) (*oss.GetUserInfoRes, error) {
	log.Info("[OSSService] GetUserInfo request start. ")
	log.Debug("[OSSService] GetUserInfo Apikey:", r.Apikey, ",TenantId:", r.TenantId, ",PlatformUserid:", r.PlatformUserid, ",Region:", r.Region, ",OssUid:", r.OssUid)
	auth := &OpenstackAPIAuthorization{Apikey: r.Apikey, TenantId: r.TenantId, PlatformUserid: r.PlatformUserid}
	op := new(GetUserInfoOp)

	op.Req = r
	res, err := o.Process(auth, op)
	return res.(*oss.GetUserInfoRes), err
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
