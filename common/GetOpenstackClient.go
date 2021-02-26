package common

import (
	sdk "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	"os"
)

// GetOpenstackClient is for creating an openstack provider client
func GetOpenstackClient(apikey string, tenantID string, platformUserID string,
	resourceID ...string) (*sdk.ProviderClient, error) {
	authResult := APIAuth(apikey, tenantID, platformUserID, resourceID...)
	if !authResult {
		return nil, EAPIAUTH
	}
	tenantInfo, err := QueryTenantInfoByTenantIdAndApikey(tenantID, apikey)
	if err != nil {
		return nil, err
	}
	identityEndPoint := os.Getenv("OPENSTACK_IDENTITY_ENDPOINT")
	//获取provider
	opts := sdk.AuthOptions{
		IdentityEndpoint: identityEndPoint,               // idenEndPoint,
		Username:         tenantInfo.OpenstackUsername,   // username,
		Password:         tenantInfo.OpenstackPassword,   // password,
		DomainName:       tenantInfo.OpenstackDomainname, // domain,
		TenantID:         tenantInfo.OpenstackProjectid,  // tenantID,
	}
	return openstack.AuthenticatedClient(opts)
}
