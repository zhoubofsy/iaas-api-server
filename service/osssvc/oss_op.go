package osssvc

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"liyongcool.nat300.top/bozhou/go-radosgw/pkg/api"
	"liyongcool.nat300.top/iaas/iaas-api-server/common"
	"liyongcool.nat300.top/iaas/iaas-api-server/proto/oss"
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
	//common.APIAuth(o.Apikey, o.TenantId, o.PlantformUserid)
	return true
}

type Op interface {
	// Predo
	Predo() error
	// Do
	Do() error
	// Done
	Done(error) (interface{}, error)
}

/*
type BaseOp struct {
}

func (o *BaseOp) Predo() error {
	return status.Errorf(codes.Unimplemented, "method Predo not implemented")
}

func (o *BaseOp) Do() error {
	return status.Errorf(codes.Unimplemented, "method Do not implemented")
}

func (o *BaseOp) Done() error {
	return status.Errorf(codes.Unimplemented, "method Done not implemented")
}
*/
type BasicOp struct {
	conf OSSConfigure
}

type S3Provider struct {
	AccessKey string
	SecretKey string
}

func (o *S3Provider) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     o.AccessKey,
		SecretAccessKey: o.SecretKey,
	}, nil
}

func (o *S3Provider) IsExpired() bool {
	return false
}

type BucketOp struct {
	EndpointAddr string
	Access       string
	Secret       string
	S3Handler    *s3.S3
}

func (o *BucketOp) init() error {
	pathStyle := true
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String("default"),
			Endpoint:         &(o.EndpointAddr),
			S3ForcePathStyle: &pathStyle,
			Credentials:      credentials.NewCredentials(&S3Provider{AccessKey: o.Access, SecretKey: o.Secret}),
		},
	}))
	o.S3Handler = s3.New(sess)
	return nil
}

func (o *BucketOp) CreateBucket(name string) error {
	o.init()
	/*
		_, err := o.S3Handler.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(name),
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{
				LocationConstraint: aws.String(""),
			},
		})
	*/
	_, err := o.S3Handler.CreateBucket(&s3.CreateBucketInput{
		Bucket: &name,
	})
	return err
}

func (o *BucketOp) GetBucketInfo(name string) error {
	// Todo...
	return nil
}

type CreateUserAndBucketOp struct {
	BasicOp
	Req *oss.CreateUserAndBucketReq
	Res *oss.CreateUserAndBucketRes
}

func (o *CreateUserAndBucketOp) Predo() error {
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(oss.CreateUserAndBucketRes)
	o.conf = GetOSSConfigure()

	return common.EOK
}

func (o *CreateUserAndBucketOp) Do() error {
	// Get Endpoint via Region from config file
	endpoint := o.conf.GetEndpointByRegion(o.Req.Region)
	// Read Access & Secret keys
	access, secret := o.conf.GetRGWAdminAccessSecretKeys()
	maxObjs := o.Req.UserMaxObjects
	quotaSize := o.Req.UserMaxSizeInG * 1024 * 1024
	rgw, err := radosAPI.New(endpoint, access, secret)
	var user *radosAPI.User
	var statusCode int
	if err == nil {
		//create s3 user
		user, err, statusCode = rgw.CreateUser(radosAPI.UserConfig{Tenant: o.Req.TenantId, UID: o.Req.PlatformUserid, DisplayName: o.Req.PlatformUserid})
		if err == nil {
			//	set user quota(size)
			err = rgw.UpdateQuota(radosAPI.QuotaConfig{UID: o.Req.PlatformUserid, QuotaType: "user", MaxSizeKB: strconv.FormatInt(int64(quotaSize), 10), Enabled: "True"})
			if err != nil {
				return common.EOSSSETQUOTAS
			}
			// set bucket quota(object count)
			err = rgw.UpdateQuota(radosAPI.QuotaConfig{UID: o.Req.PlatformUserid, QuotaType: "bucket", MaxObjects: strconv.FormatInt(int64(maxObjs), 10), Enabled: "True"})
			if err != nil {
				return common.EOSSSETQUOTAS
			}
		}
	}
	if err != nil {
		if statusCode == 409 {
			// Read user info
			user, err = rgw.GetUser(o.Req.PlatformUserid)
		} else {
			return common.EOSSCREATEUSER
		}
	}

	//	create bucket
	bucketOperator := BucketOp{EndpointAddr: endpoint, Access: user.Keys[0].AccessKey, Secret: user.Keys[0].SecretKey}
	err = bucketOperator.CreateBucket(o.Req.BucketName)
	if err != nil {
		return common.EOSSCREATEBUCKET
	}
	o.Res.OssEndpoint = endpoint
	o.Res.OssAccessKey = user.Keys[0].AccessKey
	o.Res.OssSecretKey = user.Keys[0].SecretKey
	o.Res.OssUser = &(oss.OssUser{OssUid: user.UserID, OssUserCreatedTime: "", UserMaxSizeInG: 0, UserMaxObjects: 0, UserUseSizeInG: 0, UserUseObjects: 0, TotalBuckets: 0})
	o.Res.OssBucket = &(oss.OssBucket{BucketName: o.Req.BucketName, BucketPolicy: "", BucketUseSizeInG: 0, BucketUseObjects: 0, BelongToUid: user.UserID, BucketCreatedTime: ""})
	return common.EOK
}

func (o *CreateUserAndBucketOp) Done(e error) (interface{}, error) {
	//Translate error code
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}

type GetBucketInfoOp struct {
	BasicOp
	Req *oss.GetBucketInfoReq
	Res *oss.GetBucketInfoRes
}

func (o *GetBucketInfoOp) Predo() error {
	return nil
}

func (o *GetBucketInfoOp) Do() error {
	return nil
}

func (o *GetBucketInfoOp) Done(e error) (interface{}, error) {
	return o.Res, e
}

type ListBucketsInfoOp struct {
	BasicOp
	Req *oss.ListBucketsInfoReq
	Res *oss.ListBucketsInfoRes
}

func (o *ListBucketsInfoOp) Predo() error {
	return nil
}

func (o *ListBucketsInfoOp) Do() error {
	return nil
}

func (o *ListBucketsInfoOp) Done(e error) (interface{}, error) {
	return o.Res, e
}

type SetOssUserQuotaOp struct {
	BasicOp
	Req *oss.SetOssUserQuotaReq
	Res *oss.SetOssUserQuotaRes
}

func (o *SetOssUserQuotaOp) Predo() error {
	return nil
}

func (o *SetOssUserQuotaOp) Do() error {
	return nil
}

func (o *SetOssUserQuotaOp) Done(e error) (interface{}, error) {
	return o.Res, e
}

type RecoverKeyOp struct {
	BasicOp
	Req *oss.RecoverKeyReq
	Res *oss.RecoverKeyRes
}

func (o *RecoverKeyOp) Predo() error {
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(oss.RecoverKeyRes)
	o.conf = GetOSSConfigure()

	return common.EOK
}

func (o *RecoverKeyOp) Do() error {
	// Get Endpoint via Region from config file
	endpoint := o.conf.GetEndpointByRegion(o.Req.Region)
	// Read Access & Secret keys
	access, secret := o.conf.GetRGWAdminAccessSecretKeys()
	rgw, err := radosAPI.New(endpoint, access, secret)
	var user *radosAPI.User
	if err == nil {
		user, err = rgw.GetUser(o.Req.PlatformUserid)
	}

	if err != nil {
		return common.EOSSGETUSER
	}
	o.Res.OssAccessKey = user.Keys[0].AccessKey
	o.Res.OssSecretKey = user.Keys[0].SecretKey
	return common.EOK
}

func (o *RecoverKeyOp) Done(e error) (interface{}, error) {
	//Translate error code
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}
