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

func InitDb() bool {
	driverName := os.Getenv("DRIVER_NAME")
	dbBHostIP := os.Getenv("DB_HOST_IP")
	dbUserName := os.Getenv("DB_USERNAME")
	dbPassWord := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dns := dbUserName + ":" + dbPassWord + "@tcp(" + dbBHostIP + ")/" + dbName
	db, _ = sql.Open(driverName, dns+"?charset=utf8")
	//设置数据库最大连接数
	db.SetConnMaxLifetime(200)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		fmt.Println("open database fail")
		return false
	}
	return true
}

func QueryNasDiskConfigByRegion(region string) (NasDiskConfig, error) {
	var config NasDiskConfig
	if !InitDb() {
		return config, ETTGETMYSQLCLIENT
	}
	defer db.Close()
	sqlStr := "SELECT id, region, mgrendpoint, mgruser, mgrpasswd, cephfsid, cephfsroot, ganeshaendpoint, ganeshaclusterid, ganeshaexportuser FROM nas_disk_config where region = ? "
	err := db.QueryRow(sqlStr, region).Scan(
		&config.ID,
		&config.Region,
		&config.MGREndpoint,
		&config.MGRUser,
		&config.MGRPasswd,
		&config.CephfsID,
		&config.CephfsRoot,
		&config.GaneshaEndpoint,
		&config.GaneshaClusterID,
		&config.GaneshaExportUser)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query NasDiskConfig by region failed.")
		return config, ENASDISKGETCONFIG
	}
	return config, nil
}

func QueryOssConfigByRegion(region string) (OssConfig, error) {
	var ossConfig OssConfig

	if !InitDb() {
		return ossConfig, ETTGETMYSQLCLIENT
	}

	defer db.Close()
	sqlStr := "SELECT id, region, access_key, secret_key, endpoint, description FROM oss_config where region = ? "
	err := db.QueryRow(sqlStr, region).Scan(&ossConfig.ID, &ossConfig.Region, &ossConfig.AccessKey, &ossConfig.SecretKey, &ossConfig.Endpoint, &ossConfig.Description)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query sssConfig by region failed.")
		return ossConfig, EOSSGETCONFIG
	}
	return ossConfig, nil
}

func QueryTenantInfoByTenantIdAndApikey(tenantID string, apiKey string) (TenantInfo, error) {
	sqlStr := "SELECT tenant_id,tenant_name,openstack_domainname,openstack_domainid,openstack_projectname,openstack_projectid,openstack_username,openstack_userid,openstack_password,openstack_rolename,openstack_roleid,apikey FROM tenant_info where tenant_id =? and apikey=?"
	var tenantInfo TenantInfo

	if !InitDb() {
		return tenantInfo, ETTGETMYSQLCLIENT
	}

	defer db.Close()
	err := db.QueryRow(sqlStr, tenantID, apiKey).Scan(&tenantInfo.TenantID, &tenantInfo.TenantName, &tenantInfo.OpenstackDomainname, &tenantInfo.OpenstackDomainid, &tenantInfo.OpenstackProjectname, &tenantInfo.OpenstackProjectid, &tenantInfo.OpenstackUsername, &tenantInfo.OpenstackUserid, &tenantInfo.OpenstackPassword, &tenantInfo.OpenstackRolename, &tenantInfo.OpenstackRoleid, &tenantInfo.ApiKey)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantInfo By tenantName failed.")
		return tenantInfo, ETTGETTENANT
	}
	return tenantInfo, nil
}

func QueryTenantInfoByTenantName(name string) (string, error) {
	sqlStr := "SELECT tenant_id,tenant_name FROM tenant_info where tenant_name =?"
	var tenantInfo TenantInfo

	if !InitDb() {
		return "", ETTGETMYSQLCLIENT
	}

	defer db.Close()
	err := db.QueryRow(sqlStr, name).Scan(&tenantInfo.TenantID, &tenantInfo.TenantName)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantInfo by tenantName failed.")
		return "", ETTGETTENANT
	}
	return tenantInfo.TenantID, nil
}

func GetTenantIDSeq() (string, error) {
	if !InitDb() {
		return "", ETTGETMYSQLCLIENT
	}
	defer db.Close()
	sqlStr := "select nextval(seq)"
	var nextVal int32
	err := db.QueryRow(sqlStr).Scan(&nextVal)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("query tenantId seq failed.")
		return "", ETTGETENATSEQ
	}
	valStr := fmt.Sprintf("%010d", nextVal)
	return valStr, nil
}

func CreateTenantInfo(tenantInfo TenantInfo) (createTenantFlag bool) {
	log.Info("dbutils insert tenantInfo :", tenantInfo)
	log.Info("create tenant flag :", createTenantFlag)
	if !InitDb() {
		return false
	}
	defer db.Close()
	sqlStr := "insert into tenant_info(tenant_id, tenant_name, openstack_domainname, openstack_domainid, openstack_projectname, openstack_projectid, openstack_username, openstack_userid, openstack_password, openstack_rolename, openstack_roleid, apikey) values (?,?,?,?,?,?,?,?,?,?,?,?)"
	ret, err := db.Exec(sqlStr, tenantInfo.TenantID, tenantInfo.TenantName, tenantInfo.OpenstackDomainname, tenantInfo.OpenstackDomainid, tenantInfo.OpenstackProjectname, tenantInfo.OpenstackProjectid, tenantInfo.OpenstackUsername, tenantInfo.OpenstackUserid, tenantInfo.OpenstackPassword, tenantInfo.OpenstackRolename, tenantInfo.OpenstackRoleid, tenantInfo.ApiKey)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("create tenantInfo failed.")
		return false
	}
	n, err := ret.RowsAffected()
	if n > 0 {
		return true
	}
	return false
}

func DeleteTenantInfo(tenantID string) error {
	if !InitDb() {
		return ETTGETMYSQLCLIENT
	}
	defer db.Close()
	sqlStr := "delete from tenant_info where tenant_id = ?"
	_, err := db.Exec(sqlStr, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("delete tenantInfo failed.")
		return ETTDELETETENANT
	}
	return nil
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

type OssConfig struct {
	ID          string
	Region      string
	AccessKey   string
	SecretKey   string
	Endpoint    string
	Description string
}

func (ossConfig OssConfig) IsEmpty() bool {
	return reflect.DeepEqual(ossConfig, OssConfig{})
}

type NasDiskConfig struct {
	ID                string
	Region            string
	MGREndpoint       string
	MGRUser           string
	MGRPasswd         string
	CephfsID          string
	CephfsRoot        string
	GaneshaEndpoint   string
	GaneshaClusterID  string
	GaneshaExportUser string
}

func (nasdiskConfig NasDiskConfig) IsEmpty() bool {
	return reflect.DeepEqual(nasdiskConfig, NasDiskConfig{})
}

// QuerySharedSubnetUsedIP for shared subnet
func QuerySharedSubnetUsedIP(subnetID string) (SharedSubnetIPPool, error) {
	var ip = SharedSubnetIPPool{
		SubnetID: subnetID,
	}
	if !InitDb() {
		return ip, EPLGETIPPOOL
	}

	defer db.Close()
	sqlStr := "SELECT used_ip FROM shared_subnet_ip_pool WHERE subnet_id =?"
	err := db.QueryRow(sqlStr, subnetID).Scan(&ip.UsedIP)
	if err != nil {
		if err == sql.ErrNoRows {
			return ip, EPLGETIPPOOLNONE
		}
		log.WithFields(log.Fields{
			"err":      err,
			"subnetID": subnetID,
			"sql":      sqlStr,
		}).Error("query shared_subnet_ip_pool failed.")
		return ip, EPLGETIPPOOL
	}
	return ip, nil
}

// CreateSharedSubnetUsedIP for shared subnet
func CreateSharedSubnetUsedIP(subnetID string, usedIP string) bool {
	if !InitDb() {
		return false
	}
	defer db.Close()
	sqlStr := "INSERT INTO shared_subnet_ip_pool(subnet_id, used_ip) VALUES (?,?)"
	ret, err := db.Exec(sqlStr, subnetID, usedIP)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"sql":      sqlStr,
			"subnetID": subnetID,
			"usedIP":   usedIP,
		}).Error("create shared subnet ip pool failed.")
		return false
	}
	n, err := ret.RowsAffected()
	if n > 0 {
		return true
	}
	return false
}

// UpdateSharedSubnetUsedIP for shared subnet
func UpdateSharedSubnetUsedIP(subnetID string, usedIP string) bool {
	if !InitDb() {
		return false
	}
	defer db.Close()
	sqlStr := "UPDATE shared_subnet_ip_pool SET used_ip = ? WHERE subnet_id = ?"
	ret, err := db.Exec(sqlStr, usedIP, subnetID)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"sql":      sqlStr,
			"subnetID": subnetID,
			"usedIP":   usedIP,
		}).Error("update shared subnet ip pool failed.")
		return false
	}
	n, err := ret.RowsAffected()
	if n > 0 {
		return true
	}
	return false
}

// SharedSubnetIPPool IP Pool
type SharedSubnetIPPool struct {
	SubnetID string
	UsedIP   string
	IPPool   string
}
