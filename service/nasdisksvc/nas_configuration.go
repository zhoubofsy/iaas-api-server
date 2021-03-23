package nasdisksvc

import (
	log "github.com/sirupsen/logrus"
	"iaas-api-server/common"
)

type NasDiskConfigure interface {
	GetMGRConfig(string) (string, string, string, error)
	GetCephFSConfig(string) (string, string, error)
	GetGaneshaConfig(string) (string, string, string, error)
}

type NasDiskSimpleConfigure struct {
}

func (o *NasDiskSimpleConfigure) GetMGRConfig(r string) (string, string, string, error) {
	return "http://120.92.19.57:60081", "admin", "1qaz2wsx", nil
}

func (o *NasDiskSimpleConfigure) GetCephFSConfig(r string) (string, string, error) {
	return "1", "/nasroot", nil
}

func (o *NasDiskSimpleConfigure) GetGaneshaConfig(r string) (string, string, string, error) {
	return "192.168.122.6", "ganesha-myfs", "admin", nil // cluster-id, user-id
}

type NasDiskMariadbConfig struct {
}

func (o *NasDiskMariadbConfig) GetMGRConfig(r string) (string, string, string, error) {
	config, err := common.QueryNasDiskConfigByRegion(r)
	if nil != err {
		log.Error("[NasDiskService] NasDiskMariadbConfig GetMGRConfig Failure. ", err)
		return "", "", "", err
	}
	return config.MGREndpoint, config.MGRUser, config.MGRPasswd, nil
}

func (o *NasDiskMariadbConfig) GetCephFSConfig(r string) (string, string, error) {
	config, err := common.QueryNasDiskConfigByRegion(r)
	if nil != err {
		log.Error("[NasDiskService] NasDiskMariadbConfig GetCephFSConfig Failure. ", err)
		return "", "", err
	}
	return config.CephfsID, config.CephfsRoot, nil
}

func (o *NasDiskMariadbConfig) GetGaneshaConfig(r string) (string, string, string, error) {
	config, err := common.QueryNasDiskConfigByRegion(r)
	if nil != err {
		log.Error("[NasDiskService] NasDiskMariadbConfig GetGaneshaConfig Failure. ", err)
		return "", "", "", err
	}
	return config.GaneshaEndpoint, config.GaneshaClusterID, config.GaneshaExportUser, nil
}

func GetNasDiskConfigure() NasDiskConfigure {
	//return new(NasDiskSimpleConfigure)
	return new(NasDiskMariadbConfig)
}
