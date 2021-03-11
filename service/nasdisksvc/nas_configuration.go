package nasdisksvc

type NasDiskConfigure interface {
	GetMGRRestfulAddr(string) (string, error)
	GetMGRUserPasswd(string) (string, string, error)
	GetCephfsID(string) (string, error)
	GetRootPath(string) (string, error)
	GetGaneshaConfig(string) (string, string, error)
}

type NasDiskSimpleConfigure struct {
}

func (o *NasDiskSimpleConfigure) GetMGRRestfulAddr(r string) (string, error) {
	return "http://120.92.19.57:60081", nil
}

func (o *NasDiskSimpleConfigure) GetMGRUserPasswd(r string) (string, string, error) {
	return "admin", "1qaz2wsx", nil
}

func (o *NasDiskSimpleConfigure) GetCephfsID(r string) (string, error) {
	return "1", nil
}

func (o *NasDiskSimpleConfigure) GetRootPath(r string) (string, error) {
	return "/nasroot", nil
}

func (o *NasDiskSimpleConfigure) GetGaneshaConfig(r string) (string, string, error) {
	return "ganesha-myfs", "admin", nil // cluster-id, user-id
}

func GetNasDiskConfigure() NasDiskConfigure {
	return new(NasDiskSimpleConfigure)
}
