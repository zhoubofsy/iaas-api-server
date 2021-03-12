package nasdisksvc

type NasDiskConfigure interface {
	GetMGRConfig(string) (string, string, string, error)
	GetCephFSConfig(string) (string, string, error)
	GetGaneshaConfig(string) (string, string, error)
}

type NasDiskSimpleConfigure struct {
}

func (o *NasDiskSimpleConfigure) GetMGRConfig(r string) (string, string, string, error) {
	return "http://120.92.19.57:60081", "admin", "1qaz2wsx", nil
}

func (o *NasDiskSimpleConfigure) GetCephFSConfig(r string) (string, string, error) {
	return "1", "/nasroot", nil
}

func (o *NasDiskSimpleConfigure) GetGaneshaConfig(r string) (string, string, error) {
	return "ganesha-myfs", "admin", nil // cluster-id, user-id
}

func GetNasDiskConfigure() NasDiskConfigure {
	return new(NasDiskSimpleConfigure)
}
