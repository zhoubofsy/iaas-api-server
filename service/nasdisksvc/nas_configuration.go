package nasdisksvc

type NasDiskConfigure interface {
	GetMGRRestfulAddr(string) (string, error)
	GetMGRUserPasswd(string) (string, string, error)
	GetCephfsID(string) (string, error)
}

type NasDiskSimpleConfigure struct {
}

func (o *NasDiskSimpleConfigure) GetMGRRestfulAddr(r string) (string, error) {
	return "http://120", nil
}

func (o *NasDiskSimpleConfigure) GetMGRUserPasswd(r string) (string, string, error) {
	return "admin", "1qaz2wsx", nil
}

func (o *NasDiskSimpleConfigure) GetCephfsID(r string) (string, error) {
	return "1", nil
}

func GetNasDiskConfigure() NasDiskConfigure {
	return new(NasDiskSimpleConfigure)
}
