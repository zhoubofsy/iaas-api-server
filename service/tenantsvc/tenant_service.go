/*================================================================
*
*  文件名称：tenant_service.go
*  创 建 者: tiantingting
*  创建日期：2021年01月29日
*
================================================================*/

package tenantsvc

import (
	"iaas-api-server/common"
	"iaas-api-server/proto/tenant"
	"iaas-api-server/randpass"
	"reflect"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users" //导入连接MySQL数据库的驱动包
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// TenantService  for tenant
type TenantService struct {
	tenant.UnimplementedTenantServiceServer
}

// CreateTenant create tenant
func (s *TenantService) CreateTenant(cxt context.Context, tenant *tenant.CreateTenantReq) (res *tenant.CreateTenantRes, Err error) {
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
		log.Error("call openstack, get openstack admin client error")
		return
	}

	//查询表里是否有该租户（校验租户名称的唯一性）
	tenantInfo := queryTenantInfoByTenantName(tenant.TenantName)
	if !tenantInfo.isEmpty() {
		return nil, common.ETTISEMPTYTENANT
	}
	//校验业务数据
	for true {
		switch FLAG {
		case "createDomain":
			//调用openStack添加域信息
			domainResult, err = createDomain(provider, tenantID)
			if err == nil {
				FLAG = "createProject"
			} else {
				termianator = true
				break
			}
		case "createProject":
			//调用openStack创建租户到指定domain内，得到租户ID
			projectResult, err = createProject(provider, tenant.TenantName, domainResult.ID)
			if err == nil {
				FLAG = "createUser"
			} else {
				termianator = true
				break
			}
		case "createUser":
			//创建project下的admin用户，指定用户角色
			userResult, err = createUser(provider, tenantID, domainResult.ID, projectResult.ID)
			if err == nil {
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
			createTenantFlag = createTenantInfo(tenantInfoCreate)
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
		if createTenantFlag {
			//删除租户信息
			deleteTenant(tenantID)
		}
		if userFlag {
			//删除用户信息
			deleteUser(provider, userResult.ID)
		}
		if projectFlag {
			// 删除项目信息
			deleteProject(provider, projectResult.ID)
		}
		if domainFlag {
			// 删除域信息
			deleteDomain(provider, domainResult.ID)
		}

		res.TenantId = ""
		res.Apikey = ""
		res.Code = 500
		res.Msg = "创建租户失败"
		return res, common.ETTDELETETENANT
	} else {
		//返回租户ID和appKey
		res.TenantId = projectResult.ID
		res.Apikey = apiKey
		res.Code = 200
		res.Msg = "创建租户成功"
		return res, common.ETTCREATETENANT
	}
}

//DB的信息
const (
	DBHostIP   = "localhost:3306"
	DBUserName = "root"
	DBPassWord = "root"
	DBName     = "iaas_api_server"
	PASSWORD   = "password"
	DSN        = DBUserName + ":" + DBPassWord + "@tcp(" + DBHostIP + ")/" + DBName
)

// TenantInfo for tenant
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

func (tenantInfo TenantInfo) isEmpty() bool {
	return reflect.DeepEqual(tenantInfo, TenantInfo{})
}

func getDBClient() (db *gorm.DB) {
	db, err := gorm.Open("mysql", DSN+"?charset=utf8")
	if nil != err {
		log.Error("call db client, get mysql client error")
		return
	}
	defer db.Close()
	return db
}

func getOpenstackClient(provider *gophercloud.ProviderClient) (*gophercloud.ServiceClient, *common.Error) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr != nil {
		log.WithFields(log.Fields{
			"err": serviceErr,
		}).Error("get identity failed.")
		return nil, common.ETTGETIDENTITYCLIENT
	}
	return sc, nil
}

func queryTenantInfoByTenantName(name string) (tenantInfo TenantInfo) {
	db := getDBClient()
	db.Find(&tenantInfo, "tenant_name=?", name)
	return tenantInfo
}

func createTenantInfo(tenantInfo TenantInfo) (createTenantFlag bool) {
	db := getDBClient()
	return !db.NewRecord(&tenantInfo)
}

func deleteTenant(tenantID string) {
	db := getDBClient()
	var tinfo = TenantInfo{tenantID: tenantID}
	db.Delete(&tinfo)
}

func deleteDomain(provider *gophercloud.ProviderClient, domainID string) *common.Error {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr == nil {
		return common.ETTGETIDENTITYCLIENT
	}
	err := domains.Delete(sc, domainID).ExtractErr()
	if err != nil {
		return common.ETTDELETEDOMAIN
	}
	return nil
}

func deleteProject(provider *gophercloud.ProviderClient, projectID string) *common.Error {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr != nil {
		return common.ETTGETIDENTITYCLIENT
	}
	err := projects.Delete(sc, projectID).ExtractErr()
	if err != nil {
		return common.ETTDELETEPROJECT
	}
	return nil

}

func deleteUser(provider *gophercloud.ProviderClient, userID string) *common.Error {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr == nil {
		return common.ETTGETIDENTITYCLIENT
	}
	err := users.Delete(sc, userID).ExtractErr()
	if err != nil {
		return common.ETTDELETEUSER
	}
	return nil

}

func createDomain(provider *gophercloud.ProviderClient, name string) (*domains.Domain, *common.Error) {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr == nil {
		createOpts := domains.CreateOpts{
			Name:        name,
			Description: name,
		}

		domain, err := domains.Create(sc, createOpts).Extract()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("create project failed.")
			return domain, common.ETTCREATEDOMAIN
		}
		return domain, nil
	}
	return nil, common.ETTGETIDENTITYCLIENT

}

func createProject(provider *gophercloud.ProviderClient, name string, domainID string) (*projects.Project, *common.Error) {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr == nil {
		createOpts := projects.CreateOpts{
			Name:        name,
			DomainID:    domainID,
			Description: name,
			Enabled:     gophercloud.Enabled,
		}

		project, err := projects.Create(sc, createOpts).Extract()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("create project failed.")
			return project, common.ETTCREATEPROJECT
		}
		return project, nil
	}
	return nil, common.ETTGETIDENTITYCLIENT
}

func createUser(provider *gophercloud.ProviderClient, name string, domainID string, projectID string) (*users.User, *common.Error) {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr == nil {
		createOpts := users.CreateOpts{
			Name:             name,
			DomainID:         domainID,
			DefaultProjectID: projectID,
			Enabled:          gophercloud.Enabled,
			Password:         PASSWORD,
		}
		user, err := users.Create(sc, createOpts).Extract()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("create user failed.")
			return nil, common.ETTCREATEUSER
		}
		return user, nil
	}
	return nil, common.ETTGETIDENTITYCLIENT
}
