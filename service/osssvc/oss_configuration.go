package osssvc

type OSSConfigure interface {
	GetEndpointByRegion(string) string
	GetRGWAdminAccessSecretKeys(string) (string, string)
}

type OSSSimpleConfigure struct {
}

func (o *OSSSimpleConfigure) GetEndpointByRegion(r string) string {
	//TODO 从配置文件中读取，r相当于Key，可能存在多个key需要比对

	return "http://120.48.27.190:80"
}

func (o *OSSSimpleConfigure) GetRGWAdminAccessSecretKeys(r string) (string, string) {
	//TODO 从数据库中获取，根据region获取AccessKey和SecretKey
	return "C3ZBITE3VS5AD4Y3YEZB", "3cZZ8D7mP0hNCiqUIYnxKmhEPmbzcCFBkr7Bz4ey"
}

func GetOSSConfigure() OSSConfigure {
	return new(OSSSimpleConfigure)
}
