/*================================================================
*
*  文件名称：dbUtils.go
*  创 建 者: tiantingting
*  创建日期：2021年02月09日
*
================================================================*/
package dbutils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"iaas-api-server/common"
)

const (
	DBHostIP   = "localhost:3306"
	DBUserName = "root"
	DBPassWord = "root"
	DBName     = "iaas_api_server"
	DSN        = DBUserName + ":" + DBPassWord + "@tcp(" + DBHostIP + ")/" + DBName
)

//Db info
var db = &sql.DB{}

func InitDb() (bool) {
	var err error
	db, err = sql.Open("mysql", DSN+"?charset=utf8")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get mysql client failed.")
		return false
	}
	return true
}

func QueryTenantInfoByTenantName(name string) (string, *common.Error) {
	sqlStr := "SELECT * FROM tenant_info where tenant_name =?"
	var tenantInfo TenantInfo
	err := db.QueryRow(sqlStr, name).Scan(&tenantInfo.TenantID, &tenantInfo.TenantName)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantInfoByTenantName failed.")
		return "", common.ETTGETTENANT
	}
	return tenantInfo.TenantID, nil
}

func CreateTenantInfo(tenantInfo TenantInfo) (createTenantFlag bool) {
	log.Info("dbutils insert tenantInfo :", tenantInfo)
	log.Info("create tenant flag :", createTenantFlag)
	sqlStr := "insert into tenant_info(tenant_id, tenant_name, openstack_domainname, openstack_domainid, openstack_projectname, openstack_projectid, openstack_username, openstack_userid, openstack_password, openstack_rolename, openstack_roleid, apikey) values (?,?,?,?,?,?,?,?,?,?,?,?)"
	ret, err := db.Exec(sqlStr, tenantInfo.TenantID, tenantInfo.TenantName, tenantInfo.OpenstackDomainname, tenantInfo.OpenstackDomainid, tenantInfo.OpenstackProjectname, tenantInfo.OpenstackProjectid, tenantInfo.OpenstackUsername, tenantInfo.OpenstackUserid, tenantInfo.OpenstackPassword, tenantInfo.OpenstackRolename, tenantInfo.OpenstackRoleid, tenantInfo.ApiKey)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("createTenantInfo failed.")
		return false
	}
	n, err := ret.RowsAffected()
	if n > 0 {
		return true
	}
	return false
}

func DeleteTenant(tenantID string) {
	sqlStr := "delete from tenant_info where tenant_id = ?"
	_, err := db.Exec(sqlStr, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("createTenantInfo failed.")
	}
}
// TenantInfo for tenant
type TenantInfo struct {
	TenantID             string
	TenantName           string
	OpenstackDomainname  string
	OpenstackDomainid    string
	OpenstackProjectname string
	OpenstackProjectid   string
	OpenstackUsername    string
	OpenstackUserid      string
	OpenstackPassword    string
	OpenstackRolename    string
	OpenstackRoleid      string
	ApiKey               string
}

// func (tenantInfo TenantInfo) isEmpty() bool {
// 	return reflect.DeepEqual(tenantInfo, TenantInfo{})
// }
