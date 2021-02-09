package common

// basic error info for IAAS
type Error struct {
	Code int32
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}

// general error code
var (
	EOK         = &Error{Code: 200, Msg: "200 OK"}
	EUNAUTHORED = &Error{Code: 401, Msg: "401 Unauthorized"}

	EPARAM            = &Error{Code: 10000, Msg: "param ckeck failed"}
	EAPIAUTH          = &Error{Code: 10001, Msg: "api auth failed"}
	EGETOPSTACKCLIENT = &Error{Code: 10002, Msg: "get openstack client failed"}

	ESGNEWNETWORK  = &Error{Code: 20000, Msg: "security group new network v2 failed"}
	ESGCREATEGROUP = &Error{Code: 20001, Msg: "security group create group failed"}
	ESGGETGROUP    = &Error{Code: 20002, Msg: "security group get group failed"}
	ESGDELGROUP    = &Error{Code: 20003, Msg: "security group delete group failed"}
	ESGUPDATEGROUP = &Error{Code: 20004, Msg: "security group update group failed"}
	ESGOPERGROUP   = &Error{Code: 20005, Msg: "security group operate group failed"}

	EOSSCREATEUSER      = &Error{Code: 31001, Msg: "OSS create user failed"}
	EOSSCREATEBUCKET    = &Error{Code: 31002, Msg: "OSS create bucket failed"}
	EOSSSETQUOTAS       = &Error{Code: 31003, Msg: "OSS set quotas failed"}
	EOSSGETUSER         = &Error{Code: 31004, Msg: "OSS get user failed"}
	EOSSSETBUCKETPOLICY = &Error{Code: 31005, Msg: "OSS set bucket policy failed"}
	EOSSUNKNOWQUOTATYPE = &Error{Code: 31006, Msg: "OSS unknow quota type"}
	EOSSGETQUOTAS       = &Error{Code: 31007, Msg: "OSS get quotas failed"}
	EOSSGETBUCKET       = &Error{Code: 31008, Msg: "OSS get bucket failed"}
	EOSSNOPAGE          = &Error{Code: 31009, Msg: "OSS no page failed"}
	EOSSLISTBUCKETS     = &Error{Code: 31010, Msg: "OSS list buckets failed"}
)
