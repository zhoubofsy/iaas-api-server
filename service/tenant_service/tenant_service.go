package tenant

import (
	"iaas/iaas-api-server/proto/tenant"
)

type server struct {
	pb.UnimplementedTenantServiceServer
}

func (s* server) CreateTenant(req* CreateTenantReq) (res* CreateTenantRes, error) {
	
}