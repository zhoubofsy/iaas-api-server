package tenant_service

import (
	"fmt"
	"iaas-api-server/proto/tenant"
	"randpass"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"golang.org/x/net/context"
)

type TenantService struct {
}

func (s *TenantService) CreateTenant(cxt context.Context, tenant *tenant.CreateTenantReq) (*tenant.CreateTenantRes, error) {
	//获取provider
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.1.211/identity",
		Username:         "admin",
		Password:         "password",
		DomainName:       "default",
		TenantID:         "6e60f8565fd54caa8b18f4e4cb501fb4",
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if nil != err {
		fmt.Printf("openstack create auth client failed. %s", err)
		return nil, err
	}

	//调用openStack添加域信息
	domainResult, domainFlag := CreateDomain(provider, tenant.TenantName)
	if !domainFlag {

	}

	//调用openStack创建租户到指定domain内，得到租户ID
	projectResult, tenantFlag := CreateProject(provider, tenant.TenantName, domainResult.ID)
	if !tenantFlag {

	}

	//创建project下的admin用户，指定用户角色
	_, userFlag := CreateUser(provider, tenant.TenantName, domainResult.ID, projectResult.ID)
	if !userFlag {

	}
	//生成app_key，创建指定租户和appKey间的关系
	apiKey := randpass.GetRandomString(10)
	//向数据库添加数据
	//返回租户ID和appKey
	return &tenant.CreateTenantRes{TenantId: projectResult.ID, Apikey: apiKey}, nil
}

func CreateDomain(provider *gophercloud.ProviderClient, name string) (result *domains.Domain, domainFlag bool) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init compute service client error:", serviceErr)
		panic(serviceErr)
	}
	createOpts := domains.CreateOpts{
		Name:        name,
		Description: name,
	}

	result, err := domains.Create(sc, createOpts).Extract()
	if err != nil {
		panic(err)
	} else {
		domainFlag = true
		return result, domainFlag
	}
}

func CreateProject(provider *gophercloud.ProviderClient, name string, domainID string) (tentant *projects.Project, tenantFlag bool) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init compute service client error:", serviceErr)
		panic(serviceErr)
	}

	createOpts := projects.CreateOpts{
		Name:        name,
		DomainID:    domainID,
		Description: name,
		Enabled:     gophercloud.Enabled,
	}

	tenant, err := projects.Create(sc, createOpts).Extract()
	if err != nil {
		panic(err)
	} else {
		tenantFlag = true
		return tenant, tenantFlag
	}
}

func CreateUser(provider *gophercloud.ProviderClient, name string, domainID string, projectID string) (user *users.User, userFlag bool) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init compute service client error:", serviceErr)
		panic(serviceErr)
	}
	createOpts := users.CreateOpts{
		Name:             name,
		DomainID:         domainID,
		DefaultProjectID: projectID,
		Enabled:          gophercloud.Enabled,
		Password:         "password",
	}
	user, err := users.Create(sc, createOpts).Extract()
	if err != nil {
		panic(err)
	} else {
		userFlag = true
		return user, userFlag
	}
}
