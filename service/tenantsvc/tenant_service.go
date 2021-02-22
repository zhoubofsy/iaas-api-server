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
	"os"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users" //导入连接MySQL数据库的驱动包
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// TenantService  for tenant
type TenantService struct {
	tenant.UnimplementedTenantServiceServer
}

var identityEndPoint,adminUserName,adminPassword,adminProjectId,defaultPwd,adminRoleId string

func getEnvValue()  {
	identityEndPoint=os.Getenv("OPENSTACK_IDENTITY_ENDPOINT")
	adminUserName =os.Getenv("OPENSTACK_ADMIN")
	adminPassword=os.Getenv("OPENSTACK_ADMIN_PWD")
	adminRoleId=os.Getenv("OPENSTACK_ADMIN_ROLE_ID")
	adminProjectId=os.Getenv("OPENSTACK_ADMIN_PROJECT_ID")
}

// CreateTenant create tenant
func (s *TenantService) CreateTenant(cxt context.Context, tenantReq *tenant.CreateTenantReq) (*tenant.CreateTenantRes, error) {
	res := &tenant.CreateTenantRes{}
	//获取seq序列
	tenantIDSeq,seqErr:=common.GetTenantIDSeq()
	if seqErr!=nil{
		log.Error("call mysql, get mysql sequence error", seqErr)
		res.Apikey = ""
		res.TenantId = ""
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		return res, common.ETTGETENATSEQ
	}
	//生成租户ID
	tenantID := "t-" + tenantIDSeq
	//生成app_key，创建指定租户和appKey间的关系
	apiKey := randpass.GetRandomString(10)
	defaultPwd=  randpass.GetRandomString(10)
	var domainFlag, projectFlag, userFlag, createTenantFlag, termianator bool
	var FLAG = "createDomain"
	var domainResult *domains.Domain
	var projectResult *projects.Project
	var userResult *users.User
	//获取环境变量
	getEnvValue()
	//获取provider
	opts := gophercloud.AuthOptions{
		IdentityEndpoint:identityEndPoint,
		Username:        adminUserName,
		Password:        adminPassword,
		DomainName:       "default",
		TenantID:        adminProjectId,
	}
	log.Info("test:--------------------------", opts)
	provider, err := openstack.AuthenticatedClient(opts)
	if nil != err {
		log.Error("call openstack, get openstack admin client error", err)
		res.Apikey = ""
		res.TenantId = ""
		res.Code = common.EGETOPSTACKCLIENT.Code
		res.Msg = common.EGETOPSTACKCLIENT.Msg
		return res, common.EGETOPSTACKCLIENT
	}

	//查询表里是否有该租户（校验租户名称的唯一性）
	tid, _ := common.QueryTenantInfoByTenantName(tenantReq.TenantName)
	if tid != "" {
		log.Info("tenant info :", tid)
		res.Apikey = ""
		res.TenantId = ""
		res.Code = common.ETTGETTENANTNOTNULL.Code
		res.Msg = common.ETTGETTENANTNOTNULL.Msg
		return res, common.ETTGETTENANTNOTNULL
	}
	//校验业务数据
	for true {
		switch FLAG {
		case "createDomain":
			//调用openStack添加域信息
			domainResult, err = createDomain(provider, tenantID)
			if err == common.EOK {
				FLAG = "createProject"
				domainFlag = true
				log.Info("create yewu flag:", FLAG)
			} else {
				log.Info("domainResult err:", err)
				termianator = true
				break
			}
		case "createProject":
			//调用openStack创建租户到指定domain内，得到租户ID
			projectResult, err = createProject(provider, tenantID, tenantReq.TenantName, domainResult.ID)
			log.Info("projectResult err:", err)
			if err == common.EOK {
				FLAG = "createUser"
				projectFlag = true
			} else {
				termianator = true
				break
			}
		case "createUser":
			//创建project下的admin用户，指定用户角色
			userResult, err = createUser(provider, tenantID, domainResult.ID, projectResult.ID)
			log.Info("userResult err:", err)
			if err == common.EOK {
				FLAG = "createUserAndRoleR"
				userFlag = true
			} else {
				termianator = true
				break
			}
		case "createUserAndRoleR":
			//建立用户和角色间的关系
			err := createUserAndRoleRelation(provider, projectResult.ID, userResult.ID)
			if err == common.EOK {
				FLAG = "createTenant"
				userFlag = true
			} else {
				termianator = true
				break
			}
		case "createTenant":
			var tenantInfoCreate = common.TenantInfo{
				TenantID: tenantID, TenantName: tenantReq.TenantName, OpenstackDomainname: domainResult.Name,
				OpenstackDomainid: domainResult.ID, OpenstackProjectname:projectResult.Name, OpenstackProjectid: projectResult.ID,
				OpenstackUsername: userResult.Name, OpenstackUserid: userResult.ID, OpenstackPassword: defaultPwd,
				OpenstackRolename: "admin", OpenstackRoleid: adminRoleId, ApiKey: apiKey,
			} //向数据库添加数据
			createTenantFlag = common.CreateTenantInfo(tenantInfoCreate)
			log.Info("createTenantFlag:", createTenantFlag)
			if createTenantFlag {
				break
			} else {
				termianator = true
				break
			}
		}
		if termianator {
			break
		}
		if domainFlag && projectFlag && userFlag && createTenantFlag {
			break
		}
	}
	if !createTenantFlag || termianator {
		//修改域的enabled属性
		_, editErr := editDomain(provider, domainResult.ID)
		log.Info("edit domain err:", editErr)
		if editErr == common.EOK {
			// 删除域信息
			err := deleteDomain(provider, domainResult.ID)
			log.Info("delete domain err:", err)
		}
		res.TenantId = ""
		res.Apikey = ""
		res.Code = common.ETTCREATETENANT.Code
		res.Msg = common.ETTCREATETENANT.Msg
		return res, common.ETTCREATETENANT
	}
	//返回租户ID和appKey
	log.Info("project result id:", projectResult.ID)
	res.TenantId = tenantID
	res.Apikey = apiKey
	res.Code = common.EOK.Code
	res.Msg = common.EOK.Msg
	return res, nil
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

func editDomain(provider *gophercloud.ProviderClient, domainID string) (*domains.Domain, *common.Error) {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr != nil {
		return nil, common.ETTGETIDENTITYCLIENT
	}
	var iFalse bool = false
	updateOpts := domains.UpdateOpts{
		Enabled: &iFalse,
	}
	log.Info("updateOpts:", updateOpts)
	domain, err := domains.Update(sc, domainID, updateOpts).Extract()
	log.Info("update domain err:", err)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("edit project failed.")
		return domain, common.ETTEDITDOMAIN
	}
	return domain, common.EOK
}

func deleteDomain(provider *gophercloud.ProviderClient, domainID string) *common.Error {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr != nil {
		return common.ETTGETIDENTITYCLIENT
	}
	err := domains.Delete(sc, domainID).ExtractErr()
	if err != nil {
		return common.ETTDELETEDOMAIN
	}
	return common.EOK
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
	return common.EOK

}

func deleteUser(provider *gophercloud.ProviderClient, userID string) *common.Error {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr != nil {
		return common.ETTGETIDENTITYCLIENT
	}
	err := users.Delete(sc, userID).ExtractErr()
	if err != nil {
		return common.ETTDELETEUSER
	}
	return common.EOK

}

func createDomain(provider *gophercloud.ProviderClient, name string) (*domains.Domain, *common.Error) {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr == nil {
		createOpts := domains.CreateOpts{
			Name:        name,
			Description: name,
		}
		log.Info("createOpts:", createOpts)
		domain, err := domains.Create(sc, createOpts).Extract()
		log.Info("domain err:", err)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("create project failed.")
			return domain, common.ETTCREATEDOMAIN
		}
		return domain, common.EOK
	}
	return nil, common.ETTGETIDENTITYCLIENT

}

func createProject(provider *gophercloud.ProviderClient, tenantID string, name string, domainID string) (*projects.Project, *common.Error) {
	sc, serviceErr := getOpenstackClient(provider)
	if serviceErr == nil {
		createOpts := projects.CreateOpts{
			Name:        tenantID,
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
		return project, common.EOK
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
			Password:         defaultPwd,
		}
		user, err := users.Create(sc, createOpts).Extract()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("create user failed.")
			return nil, common.ETTCREATEUSER
		}
		return user, common.EOK
	}
	return nil, common.ETTGETIDENTITYCLIENT
}

func createUserAndRoleRelation(provider *gophercloud.ProviderClient, projectId string, userId string) *common.Error {
	sc, serviceErr := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if serviceErr == nil {
		err := roles.Assign(sc, adminRoleId, roles.AssignOpts{
			UserID:    userId,
			ProjectID: projectId,
		}).ExtractErr()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("create user and role relation failed.")
			return common.ETTCREATEUSERANDROLER
		}
		return common.EOK
	}
	return common.ETTGETIDENTITYCLIENT
}
