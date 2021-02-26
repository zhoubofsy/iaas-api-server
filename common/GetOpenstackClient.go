package common

import (
	sdk "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	"iaas-api-server/common/config"
)

// GetOpenstackClient is for creating an openstack provider client
func GetOpenstackClient(apikey string, tenantID string, platformUserID string,
	resourceID ...string) (*sdk.ProviderClient, error) {
<<<<<<< HEAD
	// TODO:
	//   1. auth
	//   2. get tenant info

	identityEndpoint, _ := config.GetString("identity_endpoint")
	username, _ := config.GetString("test_username")
	passwd, _ := config.GetString("test_passwd")
	prjID, _ := config.GetString("test_project_id")

	opts := sdk.AuthOptions{
		//IdentityEndpoint: "http://192.168.247.130/identity", // idenEndPoint,
		IdentityEndpoint: identityEndpoint, // idenEndPoint,
		Username:         username,         // username,
		Password:         passwd,           // password,
		DomainName:       "default",        // domain,
		TenantID:         prjID,            // projectID,
=======
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
>>>>>>> master
	}

	return openstack.AuthenticatedClient(opts)
}
