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
)
