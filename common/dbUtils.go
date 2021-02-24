/*================================================================
*
*  文件名称：dbUtils.go
*  创 建 者: tiantingting
*  创建日期：2021年02月09日
*
================================================================*/
package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
)
//Db info
var db = &sql.DB{}

func InitDb() (bool) {
	driverName:=os.Getenv("DRIVER_NAME")
	dbBHostIP:=os.Getenv("DB_HOST_IP")
	dbUserName:=os.Getenv("DB_USERNAME")
	dbPassWord:=os.Getenv("DB_PASSWORD")
	dbName:=os.Getenv("DB_NAME")
	dns:= dbUserName + ":" + dbPassWord + "@tcp(" + dbBHostIP + ")/" + dbName
	var err error
	db, err = sql.Open(driverName, dns+"?charset=utf8")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get mysql client failed.")
		return false
	}
	return true
}
func QueryTenantInfoByTenantIdAndApikey(tenantID string,apiKey string) (TenantInfo,error) {
	sqlStr := "SELECT tenant_id,tenant_name,openstack_domainname,openstack_domainid,openstack_projectname,openstack_projectid,openstack_username,openstack_userid,openstack_password,openstack_rolename,openstack_roleid,apikey FROM tenant_info where tenant_id =? and apikey=?"
	var tenantInfo TenantInfo
	err := db.QueryRow(sqlStr,tenantID,apiKey).Scan(&tenantInfo.TenantID, &tenantInfo.TenantName,&tenantInfo.OpenstackDomainname,&tenantInfo.OpenstackDomainid,&tenantInfo.OpenstackProjectname,&tenantInfo.OpenstackProjectid,&tenantInfo.OpenstackUsername,&tenantInfo.OpenstackUserid,&tenantInfo.OpenstackPassword,&tenantInfo.OpenstackRolename,&tenantInfo.OpenstackRoleid,&tenantInfo.ApiKey)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantInfoByTenantName failed.")
		return tenantInfo, ETTGETTENANT
	}
	return tenantInfo, nil
}

func QueryTenantInfoByTenantName(name string) (string, error) {
	sqlStr := "SELECT tenant_id,tenant_name FROM tenant_info where tenant_name =?"
	var tenantInfo TenantInfo
	err := db.QueryRow(sqlStr, name).Scan(&tenantInfo.TenantID, &tenantInfo.TenantName)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantInfoByTenantName failed.")
		return "", ETTGETTENANT
	}
	return tenantInfo.TenantID, nil
}

func GetTenantIDSeq() (string,error) {
	sqlStr :="select nextval(seq)"
	var nextVal int32
	err := db.QueryRow(sqlStr).Scan(&nextVal)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantId seq failed.")
		return "", ETTGETENATSEQ
	}
	valStr:=fmt.Sprintf("%010d",nextVal)
	return valStr, nil
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
		}).Error("delete tenant info failed.")
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

func (tenantInfo TenantInfo) IsEmpty() bool {
	return reflect.DeepEqual(tenantInfo, TenantInfo{})
}
