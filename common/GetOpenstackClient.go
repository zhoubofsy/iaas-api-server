package common

import (
	sdk "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
)

// GetOpenstackClient is for creating an openstack provider client
func GetOpenstackClient(apikey string, tenantID string, platformUserID string,
	resourceID ...string) (*sdk.ProviderClient, error) {
	// TODO:
	//   1. auth
	//   2. get tenant info

	//opts := sdk.AuthOptions{
	//	IdentityEndpoint: "", // idenEndPoint,
	//	Username:         "", // username,
	//	Password:         "", // password,
	//	DomainName:       "", // domain,
	//	TenantID:         "", // tenantID,
	//}

	//获取provider
	opts := sdk.AuthOptions{
		IdentityEndpoint: "http://120.92.19.57:5000/identity",
		Username:         "admin",
		Password:         "ADMIN_PASS",
		DomainName:       "default",
		TenantID:         "b37bb68ac46943bdb134a7861553380a",
	}

	return openstack.AuthenticatedClient(opts)
}
