package common

import (
	sdk "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
)

// CreateProviderClient is for creating an openstack provider client
//   @idenEndPoint:   "http://192.168.1.168/identity"
//   @username:       "admin"
//   @password:       "your password"
//   @domain:         "default"
//   @tenantID:       "6e60f8565fd54caa8b18f4e4cb501fb4"
func CreateProviderClient(idenEndPoint string, username string, password string,
	domain string, tenantID string) (*sdk.ProviderClient, error) {
	opts := sdk.AuthOptions{
		IdentityEndpoint: idenEndPoint,
		Username:         username,
		Password:         password,
		DomainName:       domain,
		TenantID:         tenantID,
	}

	return openstack.AuthenticatedClient(opts)
}
