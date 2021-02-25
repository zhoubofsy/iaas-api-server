package osssvc

import (
	"iaas-api-server/common"
	"iaas-api-server/proto/oss"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/zhoubofsy/go-radosgw/pkg/api"
)

func TransUTCTime(t string) string {
	if t == "" {
		return ""
	}
	ut, _ := time.Parse(time.RFC3339, t)
	zt := ut.In(time.Local)
	return zt.Format("2006-01-02 15:04:05")
}

type Authorization interface {
	Auth() bool
}

type OpenstackAPIAuthorization struct {
	Apikey         string
	TenantId       string
	PlatformUserid string
}

func (o *OpenstackAPIAuthorization) Auth() bool {
	return common.APIAuth(o.Apikey, o.TenantId,o.PlatformUserid)
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
	AdmAccess    string
	AdmSecret    string
	RGWHandler   *radosAPI.API
	buckets      []*s3.Bucket
}

type Quota struct {
	Enabled    bool
	MaxSize    int
	MaxObjects int
}
type BucketInfo struct {
	Name        string
	UsedSize    int // kb
	UsedObjects int
	Owner       string
	CreatedTime string
	BucketQuota Quota
}

func (o *BucketOp) Init() error {
	var err error
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
	o.RGWHandler, err = radosAPI.New(o.EndpointAddr, o.AdmAccess, o.AdmSecret)
	return err
}

func (o *BucketOp) CreateBucket(name string) error {
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

func (o *BucketOp) GetBucketInfo(name string) (*BucketInfo, error) {
	// Get Bucket Info , Include : CreateTime, Size of used, Objects os used, Owner
	buckets, err := o.RGWHandler.GetBucket(radosAPI.BucketConfig{Bucket: name, Stats: true})
	if err != nil {
		return nil, err
	}
	bkt := buckets[0]
	return &BucketInfo{
		Name:        bkt.Stats.Bucket,
		Owner:       bkt.Stats.Owner,
		CreatedTime: bkt.Stats.Mtime,
		UsedSize:    bkt.Stats.Usage.RgwMain.SizeKb,
		UsedObjects: bkt.Stats.Usage.RgwMain.NumObjects,
		BucketQuota: Quota{
			Enabled:    bkt.Stats.BucketQuota.Enabled,
			MaxSize:    bkt.Stats.BucketQuota.MaxSizeKb,
			MaxObjects: bkt.Stats.BucketQuota.MaxObjects,
		},
	}, err
}

func (o *BucketOp) GetPolicy(name string) (string, error) {
	res, err := o.S3Handler.GetBucketPolicy(&s3.GetBucketPolicyInput{Bucket: aws.String(name)})
	if err != nil {
		return "", err
	}
	return *(res.Policy), err
}

func (o *BucketOp) SetPolicy(name string, policy string) error {
	if policy == "" {
		return nil
	}
	_, err := o.S3Handler.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(name),
		Policy: aws.String(policy),
	})
	return err
}

func (o *BucketOp) ListBucketsInit() error {
	output, err := o.S3Handler.ListBuckets(&s3.ListBucketsInput{})
	if err == nil {
		o.buckets = output.Buckets
	}
	return err
}

func (o *BucketOp) ListBucketsCount() int {
	if o.buckets == nil {
		return 0
	}
	return len(o.buckets)
}

func (o *BucketOp) ListBucketsPage(num int, size int) ([]*s3.Bucket, error) {
	startIdx := (num - 1) * size
	endIdx := startIdx + size
	if startIdx >= len(o.buckets) {
		return nil, common.EOSSNOPAGE
	}
	return o.buckets[startIdx:endIdx], nil
}

type UserOp struct {
	Access       string
	Secret       string
	EndpointAddr string
	RGWHandler   *radosAPI.API
}

type UserInfo struct {
	Uid          string
	Display      string
	AccessKey    string
	SecretKey    string
	UsedSize     int //kb
	UsedObjects  int
	BucketsQuota Quota
	UserQuota    Quota
	CreatedTime  string
}

func (o *UserOp) Init() error {
	var err error
	o.RGWHandler, err = radosAPI.New(o.EndpointAddr, o.Access, o.Secret)
	return err
}

func (o *UserOp) CreateUser(uid string, display string) (error, int) {
	_, err, status := o.RGWHandler.CreateUser(radosAPI.UserConfig{UID: uid, DisplayName: display})
	if err != nil {
		if status != 409 {
			return common.EOSSCREATEUSER, status
		}
	}

	return err, status
}

func (o *UserOp) GetUserInfo(uid string) (*UserInfo, error) {
	var userInfo UserInfo
	user, err := o.RGWHandler.GetUserInfo(radosAPI.UserInfoConfig{UID: uid, Stats: "true", Sync: "true"})
	if err != nil {
		return nil, common.EOSSGETUSER
	}
	bktQuota, err := o.GetQuota(uid, "bucket")
	if err != nil {
		return nil, common.EOSSGETQUOTAS
	}
	userQuota, err := o.GetQuota(uid, "user")
	if err != nil {
		return nil, common.EOSSGETQUOTAS
	}

	userInfo.Uid = user.UID
	userInfo.Display = user.Display
	userInfo.AccessKey = user.Keys[0].AccessKey
	userInfo.SecretKey = user.Keys[0].SecretKey
	userInfo.UsedSize = user.Stats.SizeKB / 1024 / 1024
	userInfo.UsedObjects = user.Stats.NumObjects
	userInfo.BucketsQuota = *bktQuota
	userInfo.UserQuota = *userQuota
	userInfo.CreatedTime = ""
	return &userInfo, nil
}

func (o *UserOp) setQuota(uid string, qtype string, sizeKB int, numObjects int) error {
	if qtype != "user" && qtype != "bucket" {
		return common.EOSSUNKNOWQUOTATYPE
	}
	needUpdate := false
	var quotaConfig radosAPI.QuotaConfig
	quotaConfig.UID = uid
	quotaConfig.QuotaType = qtype

	if sizeKB != 0 {
		quotaConfig.MaxSizeKB = strconv.Itoa(sizeKB)
		needUpdate = true
	}

	if numObjects != 0 {
		quotaConfig.MaxObjects = strconv.Itoa(numObjects)
		needUpdate = true
	}

	if needUpdate {
		quotaConfig.Enabled = "True"
		o.RGWHandler.UpdateQuota(quotaConfig)
	}

	return nil
}

func (o *UserOp) GetQuota(uid string, qtype string) (*Quota, error) {
	var quota Quota
	var quotaConfig radosAPI.QuotaConfig
	quotaConfig.UID = uid
	quotaConfig.QuotaType = qtype

	q, err := o.RGWHandler.GetQuotas(quotaConfig)
	if err != nil {
		return nil, err
	}
	if qtype == "user" {
		quota.Enabled = q.Enabled
		quota.MaxSize = q.MaxSizeKb
		quota.MaxObjects = q.MaxObjects
	} else {
		quota.Enabled = q.Enabled
		quota.MaxSize = q.MaxSizeKb
		quota.MaxObjects = q.MaxObjects
	}
	return &quota, err
}

func (o *UserOp) SetUserSizeQuota(uid string, sizeKB int) error {
	return o.setQuota(uid, "user", sizeKB, 0)
}

func (o *UserOp) SetUserObjectsQuota(uid string, numObjects int) error {
	return o.setQuota(uid, "user", 0, numObjects)
}

func (o *UserOp) SetBucketsSizeQuotaInUser(uid string, sizeKB int) error {
	return o.setQuota(uid, "bucket", sizeKB, 0)
}

func (o *UserOp) SetBucketsObjectsQuotaInUser(uid string, numObjects int) error {
	return o.setQuota(uid, "bucket", 0, numObjects)
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
	endpoint, err := o.conf.GetEndpointByRegion(o.Req.Region)
	if nil != err {
		return err
	}
	// Read Access & Secret keys
	access, secret, err := o.conf.GetRGWAdminAccessSecretKeys(o.Req.Region)
	if nil != err {
		return err
	}
	maxObjs := o.Req.UserMaxObjects
	maxSize := o.Req.UserMaxSizeInG * 1024 * 1024

	userOperator := UserOp{EndpointAddr: endpoint, Access: access, Secret: secret}
	userOperator.Init()
	err, status := userOperator.CreateUser(o.Req.PlatformUserid, o.Req.PlatformUserid)
	if err != nil {
		return err
	}
	if status != 409 {
		err = userOperator.SetUserSizeQuota(o.Req.PlatformUserid, int(maxSize))
		if err != nil {
			return err
		}
		err = userOperator.SetBucketsObjectsQuotaInUser(o.Req.PlatformUserid, int(maxObjs))
		if err != nil {
			return err
		}
	}

	// Get User Info
	userInfo, err := userOperator.GetUserInfo(o.Req.PlatformUserid)
	if err != nil {
		return err
	}
	//	create bucket
	bucketOperator := BucketOp{EndpointAddr: endpoint, Access: userInfo.AccessKey, Secret: userInfo.SecretKey, AdmAccess: access, AdmSecret: secret}
	bucketOperator.Init()
	err = bucketOperator.CreateBucket(o.Req.BucketName)
	if err != nil {
		return common.EOSSCREATEBUCKET
	} else {
		if nil != bucketOperator.SetPolicy(o.Req.BucketName, o.Req.BucketPolicy) {
			return common.EOSSSETBUCKETPOLICY
		}
	}
	bucketInfo, err := bucketOperator.GetBucketInfo(o.Req.BucketName)
	if err != nil {
		return common.EOSSGETBUCKET
	}
	bucketPolicy, err := bucketOperator.GetPolicy(bucketInfo.Name)
	err = bucketOperator.ListBucketsInit()

	o.Res.OssEndpoint = endpoint
	o.Res.OssAccessKey = userInfo.AccessKey
	o.Res.OssSecretKey = userInfo.SecretKey
	o.Res.OssUser = &(oss.OssUser{
		OssUid:             userInfo.Uid,
		OssUserCreatedTime: time.Now().Format("2006-01-02 15:04:05"), // Current Time
		UserMaxSizeInG:     int32(userInfo.UserQuota.MaxSize),
		UserMaxObjects:     int32(userInfo.UserQuota.MaxObjects),
		UserUseSizeInG:     int32(userInfo.UsedSize),
		UserUseObjects:     int32(userInfo.UsedObjects),
		TotalBuckets:       int32(bucketOperator.ListBucketsCount())})
	o.Res.OssBucket = &(oss.OssBucket{
		BucketName:        bucketInfo.Name,
		BucketPolicy:      bucketPolicy,
		BucketUseSizeInG:  int32(bucketInfo.UsedSize / 1024 / 1024),
		BucketUseObjects:  int32(bucketInfo.UsedObjects),
		BelongToUid:       bucketInfo.Owner,
		BucketCreatedTime: TransUTCTime(bucketInfo.CreatedTime)})
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
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(oss.GetBucketInfoRes)
	o.Res.OssBucket = new(oss.OssBucket)
	o.conf = GetOSSConfigure()

	return common.EOK
}

func (o *GetBucketInfoOp) Do() error {
	// Get Endpoint via Region from config file
	endpoint, err := o.conf.GetEndpointByRegion(o.Req.Region)
	if nil != err {
		return err
	}
	// Read Access & Secret keys
	access, secret, err := o.conf.GetRGWAdminAccessSecretKeys(o.Req.Region)
	if nil != err {
		return err
	}

	userOperator := UserOp{EndpointAddr: endpoint, Access: access, Secret: secret}
	userOperator.Init()

	userInfo, err := userOperator.GetUserInfo(o.Req.PlatformUserid)
	if err != nil {
		return common.EOSSGETUSER
	}
	bucketOperator := BucketOp{EndpointAddr: endpoint, Access: userInfo.AccessKey, Secret: userInfo.SecretKey, AdmAccess: access, AdmSecret: secret}
	bucketOperator.Init()

	bucketInfo, err := bucketOperator.GetBucketInfo(o.Req.BucketName)
	if err != nil {
		return common.EOSSGETBUCKET
	}
	bucketPolicy, err := bucketOperator.GetPolicy(bucketInfo.Name)

	o.Res.OssBucket.BucketName = bucketInfo.Name
	o.Res.OssBucket.BucketPolicy = bucketPolicy
	o.Res.OssBucket.BucketUseSizeInG = int32(bucketInfo.UsedSize / 1024 / 1024)
	o.Res.OssBucket.BucketUseObjects = int32(bucketInfo.UsedObjects)
	o.Res.OssBucket.BelongToUid = bucketInfo.Owner
	o.Res.OssBucket.BucketCreatedTime = TransUTCTime(bucketInfo.CreatedTime)
	return common.EOK
}

func (o *GetBucketInfoOp) Done(e error) (interface{}, error) {
	//Translate error code
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}

type ListBucketsInfoOp struct {
	BasicOp
	Req *oss.ListBucketsInfoReq
	Res *oss.ListBucketsInfoRes
}

func (o *ListBucketsInfoOp) Predo() error {
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	if o.Req.PageNumber <= 0 || o.Req.PageSize <= 0 {
		return common.EPARAM
	}
	o.Res = new(oss.ListBucketsInfoRes)
	o.Res.OssBuckets = make([]*oss.OssBucket, o.Req.PageSize, o.Req.PageSize)
	o.conf = GetOSSConfigure()

	return common.EOK
}

func (o *ListBucketsInfoOp) Do() error {
	// Get Endpoint via Region from config file
	endpoint, err := o.conf.GetEndpointByRegion(o.Req.Region)
	if nil != err {
		return err
	}
	// Read Access & Secret keys
	access, secret, err := o.conf.GetRGWAdminAccessSecretKeys(o.Req.Region)
	if nil != err {
		return err
	}

	userOperator := UserOp{EndpointAddr: endpoint, Access: access, Secret: secret}
	userOperator.Init()
	userInfo, err := userOperator.GetUserInfo(o.Req.OssUid)
	if err != nil {
		return common.EOSSGETUSER
	}
	bucketOperator := BucketOp{EndpointAddr: endpoint, Access: userInfo.AccessKey, Secret: userInfo.SecretKey, AdmAccess: access, AdmSecret: secret}
	bucketOperator.Init()
	err = bucketOperator.ListBucketsInit()
	if err != nil {
		return common.EOSSLISTBUCKETS
	}
	buckets, err := bucketOperator.ListBucketsPage(int(o.Req.PageNumber), int(o.Req.PageSize))
	if err != nil {
		return err
	}
	for idx, bucket := range buckets {
		bucketInfo, _ := bucketOperator.GetBucketInfo(*bucket.Name)
		bucketPolicy, _ := bucketOperator.GetPolicy(*bucket.Name)

		ossBucket := new(oss.OssBucket)
		ossBucket.BucketName = bucketInfo.Name
		ossBucket.BucketPolicy = bucketPolicy
		ossBucket.BucketUseSizeInG = int32(bucketInfo.UsedSize / 1024 / 1024)
		ossBucket.BucketUseObjects = int32(bucketInfo.UsedObjects)
		ossBucket.BelongToUid = bucketInfo.Owner
		ossBucket.BucketCreatedTime = TransUTCTime(bucketInfo.CreatedTime)
		o.Res.OssBuckets[idx] = ossBucket
	}
	return common.EOK
}

func (o *ListBucketsInfoOp) Done(e error) (interface{}, error) {
	//Translate error code
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
	return o.Res, e
}

type SetOssUserQuotaOp struct {
	BasicOp
	Req *oss.SetOssUserQuotaReq
	Res *oss.SetOssUserQuotaRes
}

func (o *SetOssUserQuotaOp) Predo() error {
	// check params
	if o.Req == nil {
		return common.EPARAM
	}
	o.Res = new(oss.SetOssUserQuotaRes)
	o.Res.OssUser = new(oss.OssUser)
	o.conf = GetOSSConfigure()

	return common.EOK
}

func (o *SetOssUserQuotaOp) Do() error {
	// Get Endpoint via Region from config file
	endpoint, err := o.conf.GetEndpointByRegion(o.Req.Region)
	if nil != err {
		return err
	}
	// Read Access & Secret keys
	access, secret, err := o.conf.GetRGWAdminAccessSecretKeys(o.Req.Region)
	if nil != err {
		return err
	}
	userOperator := UserOp{EndpointAddr: endpoint, Access: access, Secret: secret}
	userOperator.Init()

	maxSize := int(o.Req.UserMaxSizeInG * 1024 * 1024)
	maxObjs := int(o.Req.UserMaxObjects)
	err = userOperator.SetUserSizeQuota(o.Req.OssUid, maxSize)
	if err != nil {
		return common.EOSSSETQUOTAS
	}
	err = userOperator.SetBucketsObjectsQuotaInUser(o.Req.OssUid, maxObjs)
	if err != nil {
		return common.EOSSSETQUOTAS
	}
	userInfo, err := userOperator.GetUserInfo(o.Req.OssUid)
	if err != nil {
		return common.EOSSGETUSER
	}
	userQuota, err := userOperator.GetQuota(o.Req.OssUid, "user")

	bucketOperator := BucketOp{EndpointAddr: endpoint, Access: userInfo.AccessKey, Secret: userInfo.SecretKey, AdmAccess: access, AdmSecret: secret}
	bucketOperator.Init()
	bucketOperator.ListBucketsInit()

	o.Res.OssUser.OssUid = userInfo.Uid
	o.Res.OssUser.OssUserCreatedTime = userInfo.CreatedTime
	o.Res.OssUser.UserMaxSizeInG = int32(userQuota.MaxSize / 1024 / 1024)
	o.Res.OssUser.UserMaxObjects = int32(userQuota.MaxObjects)
	o.Res.OssUser.UserUseSizeInG = int32(userInfo.UsedSize / 1024 / 1024)
	o.Res.OssUser.UserUseObjects = int32(userInfo.UsedObjects)
	o.Res.OssUser.TotalBuckets = int32(bucketOperator.ListBucketsCount())

	return common.EOK
}

func (o *SetOssUserQuotaOp) Done(e error) (interface{}, error) {
	//Translate error code
	o.Res.Msg = e.Error()
	if e == common.EOK {
		o.Res.Code = common.EOK.Code
		return o.Res, nil
	}
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
	endpoint, err := o.conf.GetEndpointByRegion(o.Req.Region)
	if nil != err {
		return err
	}
	// Read Access & Secret keys
	access, secret, err := o.conf.GetRGWAdminAccessSecretKeys(o.Req.Region)
	if nil != err {
		return err
	}
	userOperator := UserOp{EndpointAddr: endpoint, Access: access, Secret: secret}
	userOperator.Init()
	userInfo, err := userOperator.GetUserInfo(o.Req.OssUid)
	if err != nil {
		return common.EOSSGETUSER
	}

	o.Res.OssAccessKey = userInfo.AccessKey
	o.Res.OssSecretKey = userInfo.SecretKey
	o.Res.OssEndpoint = endpoint
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
