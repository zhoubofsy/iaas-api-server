package common

import (
	sdk "github.com/gophercloud/gophercloud"
	openstack "github.com/gophercloud/gophercloud/openstack"
	"iaas-api-server/common/config"
)

// GetOpenstackClient is for creating an openstack provider client
func GetOpenstackClient(apikey string, tenantID string, platformUserID string,
	resourceID ...string) (*sdk.ProviderClient, error) {
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
	}

	return openstack.AuthenticatedClient(opts)
}
