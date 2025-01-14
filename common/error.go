package common

// Error basic error info for IAAS
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
	ENETWORKCLIENT    = &Error{Code: 10003, Msg: "get openstack network client failed"}
	ECOMPUTECLIENT    = &Error{Code: 10004, Msg: "get openstack compute client failed"}
	EPARSECIDR        = &Error{Code: 10005, Msg: "parse CIDR failed"}

	ESGNEWNETWORK  = &Error{Code: 20000, Msg: "security group new network v2 failed"}
	ESGCREATEGROUP = &Error{Code: 20001, Msg: "security group create group failed"}
	ESGGETGROUP    = &Error{Code: 20002, Msg: "security group get group failed"}
	ESGDELGROUP    = &Error{Code: 20003, Msg: "security group delete group failed"}
	ESGUPDATEGROUP = &Error{Code: 20004, Msg: "security group update group failed"}
	ESGOPERGROUP   = &Error{Code: 20005, Msg: "security group operate group failed"}

	ENETWORKSCREATE = &Error{Code: 21000, Msg: "networks create failed"}
	ENETWORKSGET    = &Error{Code: 21001, Msg: "networks get failed"}
	ENETWORKSSET    = &Error{Code: 21002, Msg: "networks update failed"}

	ESUBNETCREATE = &Error{Code: 22000, Msg: "subnet create failed"}
	ESUBNETGET    = &Error{Code: 22001, Msg: "subnets get failed"}

	EROUTERCREATE   = &Error{Code: 23000, Msg: "routers create failed"}
	EROUTERLIST     = &Error{Code: 23001, Msg: "routers list failed"}
	EROUTEREXTRACT  = &Error{Code: 23002, Msg: "routers extract failed"}
	EROUTERINFO     = &Error{Code: 23003, Msg: "routersInfo is not unique"}
	EROUTERGET      = &Error{Code: 23004, Msg: "routers get failed"}
	EROUTERSET      = &Error{Code: 23005, Msg: "routers set failed"}
	EROUTERNOGATWAY = &Error{Code: 23006, Msg: "The router has no external gateway"}

	EPORTSGET     = &Error{Code: 24000, Msg: "ports get failed"}
	EPORTSLIST    = &Error{Code: 24001, Msg: "ports list failed"}
	EPORTSEXTRACT = &Error{Code: 24002, Msg: "ports extract failed"}

	EINTERFACEADD = &Error{Code: 25000, Msg: "router add interface failed"}

	EFLOATINGIPCREATE       = &Error{Code: 26000, Msg: "floating ip create failed"}
	EFLOATINGIPASSOCIATE    = &Error{Code: 26001, Msg: "floating ip associate instance failed"}
	EFLOATINGIPDISASSOCIATE = &Error{Code: 26002, Msg: "floating ip disassociate instance failed"}
	EFLOATINGIPLIST         = &Error{Code: 26003, Msg: "floating ip list failed"}
	EFLOATINGIPEXTRACT      = &Error{Code: 26004, Msg: "floating ip extract pages failed"}
	EFLOATINGIPDELETE       = &Error{Code: 26005, Msg: "floating ip delete failed"}

	ENGCREATE = &Error{Code: 30000, Msg: "nat gateway create failed"}
	ENGDELETE = &Error{Code: 30001, Msg: "nat gateway delete failed"}
	ENGGET    = &Error{Code: 30002, Msg: "nat gateway get failed"}

	ENEWCPU      = &Error{Code: 40000, Msg: "nova new compute v2 failed"}
	ENFLVLIST    = &Error{Code: 40001, Msg: "nova flavor list failed"}
	ENFLVEXTRACT = &Error{Code: 40002, Msg: "nova flavor extract failed"}
	ENFLVGET     = &Error{Code: 40003, Msg: "nova flavor get failed"}

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
	ENINSSTATUS        = &Error{Code: 40110, Msg: "nova instance update flavor invalid status"}

	ENEWBLOCK         = &Error{Code: 50000, Msg: "cinder new block v3 failed"}
	ENEWVOLUME        = &Error{Code: 50001, Msg: "cinder create volume failed"}
	EDELETEVOLUME     = &Error{Code: 50002, Msg: "cinder delete volume failed"}
	ESHOWVOLUME       = &Error{Code: 50003, Msg: "cinder show volume info failed"}
	EEXTENDVOLUMESIZE = &Error{Code: 50004, Msg: "cinder extend volume size failed"}
	EVOLUMEUPDATE     = &Error{Code: 50005, Msg: "cinder update volume failed"}
	EVOLUMEATTACH     = &Error{Code: 50006, Msg: "cinder volume attach failed"}
	EVOLUMEDETACH     = &Error{Code: 50006, Msg: "cinder volume detach failed"}

	EOSSCREATEUSER      = &Error{Code: 51001, Msg: "OSS create user failed"}
	EOSSCREATEBUCKET    = &Error{Code: 51002, Msg: "OSS create bucket failed"}
	EOSSSETQUOTAS       = &Error{Code: 51003, Msg: "OSS set quotas failed"}
	EOSSGETUSER         = &Error{Code: 51004, Msg: "OSS get user failed"}
	EOSSSETBUCKETPOLICY = &Error{Code: 51005, Msg: "OSS set bucket policy failed"}
	EOSSUNKNOWQUOTATYPE = &Error{Code: 51006, Msg: "OSS unknow quota type"}
	EOSSGETQUOTAS       = &Error{Code: 51007, Msg: "OSS get quotas failed"}
	EOSSGETBUCKET       = &Error{Code: 51008, Msg: "OSS get bucket failed"}
	EOSSNOPAGE          = &Error{Code: 51009, Msg: "OSS no page failed"}
	EOSSLISTBUCKETS     = &Error{Code: 51010, Msg: "OSS list buckets failed"}

	ENASPATHEXISTED = &Error{Code: 52001, Msg: "Path of NAS is existed"}

	ETTGETTENANT        = &Error{Code: 90000, Msg: "tenant get failed"}
	ETTCREATETENANT     = &Error{Code: 90001, Msg: "tenant create failed"}
	ETTDELETETENANT     = &Error{Code: 90002, Msg: "tenant delete failed"}
	ETTISEMPTYTENANT    = &Error{Code: 90003, Msg: "tenant info is empty"}
	ETTGETTENANTNOTNULL = &Error{Code: 90004, Msg: "tenant info exits"}

	ETTGETMYSQLCLIENT           = &Error{Code: 91001, Msg: "mysql client get failed"}
	ETTGETENATSEQ               = &Error{Code: 91002, Msg: "mysql get seq failed"}
	ETTTRANS                    = &Error{Code: 91003, Msg: "seq transform failed"}
	EOSSGETCONFIG               = &Error{Code: 91005, Msg: "get oss config failed"}
	EPARSE                      = &Error{Code: 91006, Msg: "parse failed"}
	ECEPHMGRMKDIR               = &Error{Code: 91007, Msg: "Ceph MGR API make directory failed"}
	ECEPHMGRRMDIR               = &Error{Code: 91008, Msg: "Ceph MGR API remove directory failed"}
	ENASGETCONFIG               = &Error{Code: 91009, Msg: "get NAS config failed"}
	EIO                         = &Error{Code: 91010, Msg: "I/O failed"}
	ECEPHMGRSETQUOTA            = &Error{Code: 91011, Msg: "Ceph MGR API set quotas failed"}
	ECEPHMGRLISTGANESHADAEMON   = &Error{Code: 91012, Msg: "Ceph MGR API list ganesha daemons"}
	ECEPHMGRLISTGANESHAEXPORT   = &Error{Code: 91012, Msg: "Ceph MGR API list ganesha exports"}
	ECEPHMGRCREATEGANESHAEXPORT = &Error{Code: 91013, Msg: "Ceph MGR API create ganesha exports"}
	ECEPHMGRDELETEGANESHAEXPORT = &Error{Code: 91014, Msg: "Ceph MGR API delete ganesha exports"}
	ENASDISKGETCONFIG           = &Error{Code: 91015, Msg: "get nas disk config failed"}
	ECEPHMGRPUTGANESHAEXPORT    = &Error{Code: 91016, Msg: "Ceph MGR API put ganesha exports"}

	ETTEDITDOMAIN = &Error{Code: 94001, Msg: "openstack update domain failed"}

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

	EPLGETPREPARE         = &Error{Code: 100000, Msg: "openstack peerlink get, prepare param failed"}
	EPLDELETEPREPARE      = &Error{Code: 100001, Msg: "openstack peerlink delete, prepare param failed"}
	EPLCREATEPREPARE      = &Error{Code: 100002, Msg: "openstack peerlink create, prepare param failed"}
	EPLCREATEADDROUTE     = &Error{Code: 100003, Msg: "openstack peerlink create, add route to router failed"}
	EPLCREATEADDINTERFACE = &Error{Code: 100004, Msg: "openstack peerlink create, add interface to router failed"}
	EPLGETIPPOOL          = &Error{Code: 100005, Msg: "openstack get ip pool failed"}
	EPLGETIPPOOLNONE      = &Error{Code: 100006, Msg: "openstack get ip pool none"}

	EFCREATE    = &Error{Code: 101000, Msg: "openstack firewall create failed"}
	EFGETGROUP  = &Error{Code: 101000, Msg: "openstack firewall get group failed"}
	EFGROUPBIND = &Error{Code: 101001, Msg: "openstack firewall group has bind already"}
	EFDELGROUP  = &Error{Code: 101002, Msg: "openstack firewall group delete failed"}
	EFUPGROUP   = &Error{Code: 101003, Msg: "openstack firewall group update failed"}
	EFSETSTATUS = &Error{Code: 101004, Msg: "openstack firewall group set status failed"}
)
