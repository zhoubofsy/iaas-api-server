package osssvc

import (
	"iaas-api-server/common"

	log "github.com/sirupsen/logrus"
)

type OSSConfigure interface {
	GetEndpointByRegion(string) (string, error)
	GetRGWAdminAccessSecretKeys(string) (string, string, error)
}

type OSSSimpleConfigure struct {
}

func (o *OSSSimpleConfigure) GetEndpointByRegion(r string) (string, error) {
	return "http://120.48.27.190:80", nil
}
func (o *OSSSimpleConfigure) GetRGWAdminAccessSecretKeys(r string) (string, string, error) {
	return "C3ZBITE3VS5AD4Y3YEZB", "3cZZ8D7mP0hNCiqUIYnxKmhEPmbzcCFBkr7Bz4ey", nil
}

type OSSSimpleConfigureNormal struct {
}

func (o *OSSSimpleConfigureNormal) GetEndpointByRegion(r string) (string, error) {
	ossConfig, err := common.QueryOssConfigByRegion(r)
	if nil != err {
		log.Error("[OSSService] OSSSimpleConfigureNormal GetEndpointByRegion Failure. ", err)
		return "", err
	}
	return ossConfig.Endpoint, nil
}
func (o *OSSSimpleConfigureNormal) GetRGWAdminAccessSecretKeys(r string) (string, string, error) {
	//从数据库中获取，根据region获取AccessKey和SecretKey
	ossConfig, err := common.QueryOssConfigByRegion(r)
	if nil != err {
		log.Error("[OSSService] OSSSimpleConfigureNormal GetRGWAdminAccessSecretKeys Failure. ", err)
		return "", "", err
	}
	return ossConfig.AccessKey, ossConfig.SecretKey, nil
}

func GetOSSConfigure() OSSConfigure {
	return new(OSSSimpleConfigureNormal)
}
