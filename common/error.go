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

	ETTGETTENANT     = &Error{Code: 90000, Msg: "tenant get failed"}
	ETTCREATETENANT  = &Error{Code: 90001, Msg: "tenant create failed"}
	ETTDELETETENANT  = &Error{Code: 90002, Msg: "tenant delete failed"}
	ETTISEMPTYTENANT = &Error{Code: 90003, Msg: "tenant info is empty"}

	ETTCREATEDOMAIN  = &Error{Code: 95001, Msg: "openstack create domain failed"}
	ETTCREATEPROJECT = &Error{Code: 95002, Msg: "openstack create project failed"}
	ETTCREATEUSER    = &Error{Code: 95003, Msg: "openstack create user failed"}

	ETTDELETEDOMAIN  = &Error{Code: 96001, Msg: "openstack delete domain failed"}
	ETTDELETEPROJECT = &Error{Code: 96002, Msg: "openstack delete project failed"}
	ETTDELETEUSER    = &Error{Code: 96003, Msg: "openstack delete user failed"}

	ETTGETIDENTITYCLIENT = &Error{Code: 97001, Msg: "openstack delete user failed"}
)
