package tenant_service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"iaas-api-server/proto/tenant"
)

type TenantService struct {
}

func (s *TenantService) CreateTenant(context.Context, *tenant.CreateTenantReq) (*tenant.CreateTenantRes, error) {

}
