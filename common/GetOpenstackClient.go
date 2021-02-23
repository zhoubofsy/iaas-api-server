package common

import (
	sdk "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	//	"os"
)

// GetOpenstackClient is for creating an openstack provider client
func GetOpenstackClient(apikey string, tenantID string, platformUserID string,
	resourceID ...string) (*sdk.ProviderClient, error) {
	//	authResult := APIAuth(apikey, tenantID, platformUserID, resourceID...)
	//	if !authResult {
	//		return nil, EAPIAUTH
	//	}
	//	tenantInfo, err := QueryTenantInfoByTenantIdAndApikey(tenantID, apikey)
	//	if err != nil {
	//		return nil, err
	//	}
	//	identityEndPoint := os.Getenv("OPENSTACK_IDENTITY_ENDPOINT")
	//获取provider
	opts := sdk.AuthOptions{
		IdentityEndpoint: "http://192.168.66.131/identity",   // idenEndPoint,
		Username:         "admin",                            // username,
		Password:         "secret",                           // password,
		DomainName:       "Default",                          // domain,
		TenantID:         "e851733194d5460c9d3c21b801fe8831", // tenantID,
	}
	return openstack.AuthenticatedClient(opts)
}
