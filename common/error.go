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

	ENEWCPU            = &Error{Code: 40000, Msg: "nova new compute v2 failed"}
	ENFLVLIST          = &Error{Code: 40001, Msg: "nova flavor list failed"}
	ENFLVEXTRACT       = &Error{Code: 40002, Msg: "nova flavor extract failed"}
	ENFLVGET           = &Error{Code: 40003, Msg: "nova flavor get failed"}
	ENINSQUERYTENANT   = &Error{Code: 40100, Msg: "nova instance query tenant info failed"}
	ENINSCREATEVOLUME  = &Error{Code: 40101, Msg: "nova instance create volume failed"}
	ENINSCREATE        = &Error{Code: 40102, Msg: "nova instance create failed"}
	ENINSQUERY         = &Error{Code: 40103, Msg: "nova instance query failed"}
	ENINSGET           = &Error{Code: 40104, Msg: "nova instance get failed"}
	ENINSUPFLAVOR      = &Error{Code: 40105, Msg: "nova instance update flavor failed"}
	ENINSCONFIRMRESIZE = &Error{Code: 40106, Msg: "nova instance update flavor confirm failed"}
	ENINSDEL           = &Error{Code: 40107, Msg: "nova instance delete failed"}
	ENINSOPUNKNOWN     = &Error{Code: 40108, Msg: "nova instance operation not supported"}
	ENINSOP            = &Error{Code: 40109, Msg: "nova instance operate failed"}

	ENEWBLOCK         = &Error{Code: 50000, Msg: "cinder new block v3 failed"}
	ENEWVOLUME        = &Error{Code: 50001, Msg: "cinder create volume failed"}
	EDELETEVOLUME     = &Error{Code: 50002, Msg: "cinder delete volume failed"}
	ESHOWVOLUME       = &Error{Code: 50003, Msg: "cinder show volume info failed"}
	EEXTENDVOLUMESIZE = &Error{Code: 50004, Msg: "cinder extend volume size failed"}
	EVOLUMEUPDATE     = &Error{Code: 50005, Msg: "cinder update volume failed"}
	EVOLUMEATTACH     = &Error{Code: 50006, Msg: "cinder volume attach failed"}
	EVOLUMEDETACH     = &Error{Code: 50006, Msg: "cinder volume detach failed"}

	ETTGETTENANT        = &Error{Code: 90000, Msg: "tenant get failed"}
	ETTCREATETENANT     = &Error{Code: 90001, Msg: "tenant create failed"}
	ETTDELETETENANT     = &Error{Code: 90002, Msg: "tenant delete failed"}
	ETTISEMPTYTENANT    = &Error{Code: 90003, Msg: "tenant info is empty"}
	ETTGETTENANTNOTNULL = &Error{Code: 90004, Msg: "tenant info exits"}

	ETTGETMYSQLCLIENT = &Error{Code: 91001, Msg: "mysql client get failed"}
	EOSSGETCONFIG     = &Error{Code: 91005, Msg: "get oss config failed"}

	ETTEDITDOMAIN = &Error{Code: 95001, Msg: "openstack uodate domain failed"}

	ETTCREATEDOMAIN       = &Error{Code: 95001, Msg: "openstack create domain failed"}
	ETTCREATEPROJECT      = &Error{Code: 95002, Msg: "openstack create project failed"}
	ETTCREATEUSER         = &Error{Code: 95003, Msg: "openstack create user failed"}
	ETTCREATEUSERANDROLER = &Error{Code: 95004, Msg: "openstack create user and role relation failed"}

	ETTDELETEDOMAIN  = &Error{Code: 96001, Msg: "openstack delete domain failed"}
	ETTDELETEPROJECT = &Error{Code: 96002, Msg: "openstack delete project failed"}
	ETTDELETEUSER    = &Error{Code: 96003, Msg: "openstack delete user failed"}

	ETTGETIDENTITYCLIENT = &Error{Code: 97001, Msg: "openstack get client failed"}

	EIGGETIMAGE   = &Error{Code: 98000, Msg: "openstack get image failed"}
	EIGLISTIMAGES = &Error{Code: 98001, Msg: "openstack list image failed"}
)
