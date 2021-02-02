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

	opts := sdk.AuthOptions{
		IdentityEndpoint: "http://192.168.66.131/identity",   // idenEndPoint,
		Username:         "admin",                            // username,
		Password:         "secret",                           // password,
		DomainName:       "default",                          // domain,
		TenantID:         "e851733194d5460c9d3c21b801fe8831", // tenantID,
	}

	return openstack.AuthenticatedClient(opts)
}
