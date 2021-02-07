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
		IdentityEndpoint: "", // idenEndPoint,
		Username:         "", // username,
		Password:         "", // password,
		DomainName:       "", // domain,
		TenantID:         "", // tenantID,
	}

	return openstack.AuthenticatedClient(opts)
}
