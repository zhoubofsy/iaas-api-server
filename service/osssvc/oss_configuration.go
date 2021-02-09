package osssvc

type OSSConfigure interface {
	GetEndpointByRegion(string) string
	GetRGWAdminAccessSecretKeys(string) (string, string)
}

type OSSSimpleConfigure struct {
}

func (o *OSSSimpleConfigure) GetEndpointByRegion(r string) string {
	return "http://120.48.27.190:80"
}

func (o *OSSSimpleConfigure) GetRGWAdminAccessSecretKeys(r string) (string, string) {
	return "C3ZBITE3VS5AD4Y3YEZB", "3cZZ8D7mP0hNCiqUIYnxKmhEPmbzcCFBkr7Bz4ey"
}

func GetOSSConfigure() OSSConfigure {
	return new(OSSSimpleConfigure)
}
