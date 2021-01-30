package tenantsvc

import (
	"fmt"
	"iaas-api-server/proto/tenant"
	"randpass"
	"reflect"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users" //导入连接MySQL数据库的驱动包
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/context"
)

const (
	DBHostsIp  = "localhost:3306"
	DBUserName = "root"
	DBPassWord = "root"
	DBName     = "iaas_api_server"
	PASSWORD   = "password"
	DSN        = DBUserName + ":" + DBPassWord + "@tcp(" + DBHostsIp + ")/" + DBName
)

type TenantService struct {
}

type TenantInfo struct {
	tenantID             string `db:"tenant_id"`
	tenantName           string `db:"tenant_name"`
	openstackDomainname  string `db:"openstack_domainname"`
	openstackDomainid    string `db:"openstack_domainid"`
	openstackProjectname string `db:"openstack_projectname"`
	openstackProjectid   string `db:"openstack_projectid"`
	openstackUsername    string `db:"openstack_username"`
	openstackUserid      string `db:"openstack_userid"`
	openstackPassword    string `db:"openstack_password"`
	openstackRolename    string `db:"openstack_rolename"`
	openstackRoleid      string `db:"openstack_roleid"`
	apiKey               string `db:"apikey"`
}

func (tenantInfo TenantInfo) IsEmpty() bool {
	return reflect.DeepEqual(tenantInfo, TenantInfo{})
}

func QueryTenantInfoByTenantName(name string) (tenantInfo TenantInfo) {
	db, err := gorm.Open("mysql", DSN+"?charset=utf8")
	check(err)
	defer db.Close()
	db.Find(&tenantInfo, "tenant_name=?", name)
	return tenantInfo
}

func CreateTenantInfo(tenantInfo TenantInfo) (createTenantFlag bool) {
	db, err := gorm.Open("mysql", DSN+"?charset=utf8")
	check(err)
	defer db.Close()
	return !db.NewRecord(&tenantInfo)
}

func DeleteTenant(tenantID string) {
	db, err := gorm.Open("mysql", DSN+"?charset=utf8")
	check(err)
	defer db.Close()
	var tinfo = TenantInfo{tenantID: tenantID}
	db.Delete(&tinfo)
}

func check(err error) {
	if nil != err {
		panic(err)
	}
}

func (s *TenantService) CreateTenant(cxt context.Context, tenant *tenant.CreateTenantReq) (*tenant.CreateTenantRes, error) {
	//生成租户ID
	tenantID := "t-" + randpass.GetRandomString(10)
	//生成app_key，创建指定租户和appKey间的关系
	apiKey := randpass.GetRandomString(10)
	var domainFlag, projectFlag, userFlag, createTenantFlag, termianator bool
	var FLAG = "createDomain"
	var domainResult *domains.Domain
	var projectResult *projects.Project
	var userResult *users.User
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

	//查询表里是否有该租户（校验租户名称的唯一性）
	tenantInfo := QueryTenantInfoByTenantName(tenant.TenantName)
	if !tenantInfo.IsEmpty() {
		return nil, nil
	}
	//校验业务数据
	for true {
		switch FLAG {
		case "createDomain":
			//调用openStack添加域信息
			domainResult, domainFlag = CreateDomain(provider, tenantID)
			if domainFlag {
				FLAG = "createProject"
			} else {
				termianator = true
				break
			}
		case "createProject":
			//调用openStack创建租户到指定domain内，得到租户ID
			projectResult, projectFlag = CreateProject(provider, tenant.TenantName, domainResult.ID)
			if projectFlag {
				FLAG = "createUser"
			} else {
				termianator = true
				break
			}
		case "createUser":
			//创建project下的admin用户，指定用户角色
			userResult, userFlag = CreateUser(provider, tenantID, domainResult.ID, projectResult.ID)
			if userFlag {
				FLAG = "createTenant"
			} else {
				termianator = true
				break
			}
		case "createTenant":
			tenantInfoCreate := TenantInfo{tenantID: tenantID, tenantName: tenant.TenantName, openstackDomainname: domainResult.Name,
				openstackDomainid: domainResult.ID, openstackProjectname: tenant.TenantName, openstackProjectid: projectResult.ID,
				openstackUsername: userResult.Name, openstackUserid: userResult.ID, openstackPassword: PASSWORD,
				openstackRolename: tenant.TenantName, openstackRoleid: tenantID, apiKey: apiKey}
			//向数据库添加数据
			createTenantFlag = CreateTenantInfo(tenantInfoCreate)
			if createTenantFlag {
				break
			} else {
				termianator = true
				break
			}
		}
		break
	}
	if termianator {
		if domainFlag {
			// 删除域信息
			DeleteDomain(provider, domainResult.ID)
		}
		if projectFlag {
			// 删除项目信息
			DeleteProject(provider, projectResult.ID)
		}
		if userFlag {
			//删除用户信息
			DeleteUser(provider, userResult.ID)
		}
		if createTenantFlag {
			//删除租户信息
			DeleteTenant(tenantID)
		}
		return &tenant.CreateTenantRes{Code: 500, Msg: "创建租户失败", TenantId: "", Apikey: ""}, nil
	} else {
		//返回租户ID和appKey
		return &tenant.CreateTenantRes{Code: 200, Msg: "创建租户成功", TenantId: projectResult.ID, Apikey: apiKey}, nil
	}
}

func DeleteDomain(provider *gophercloud.ProviderClient, domainID string) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init identity service client error:", serviceErr)
		panic(serviceErr)
	}
	err := domains.Delete(sc, domainID).ExtractErr()
	if err != nil {
		panic(err)
	}
}

func DeleteProject(provider *gophercloud.ProviderClient, projectID string) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init identity service client error:", serviceErr)
		panic(serviceErr)
	}
	err := projects.Delete(sc, projectID).ExtractErr()
	if err != nil {
		panic(err)
	}
}

func DeleteUser(provider *gophercloud.ProviderClient, userID string) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init identity service client error:", serviceErr)
		panic(serviceErr)
	}
	err := users.Delete(sc, userID).ExtractErr()
	if err != nil {
		panic(err)
	}
}

func CreateDomain(provider *gophercloud.ProviderClient, name string) (result *domains.Domain, domainFlag bool) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		fmt.Println("init identity service client error:", serviceErr)
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
		fmt.Println("init identity service client error:", serviceErr)
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
		fmt.Println("init identity service client error:", serviceErr)
		panic(serviceErr)
	}
	createOpts := users.CreateOpts{
		Name:             name,
		DomainID:         domainID,
		DefaultProjectID: projectID,
		Enabled:          gophercloud.Enabled,
		Password:         PASSWORD,
	}
	user, err := users.Create(sc, createOpts).Extract()
	if err != nil {
		panic(err)
	} else {
		userFlag = true
		return user, userFlag
	}
}
