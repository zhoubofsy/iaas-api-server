package tenant_service

import (
	"google.golang.org/grpc"
	"iaas-api-server/proto/tenant"
	"golang.org/x/net/context"
)

type TenantService struct {

}

func (s *TenantService) CreateTenant(context.Context, *tenant.CreateTenantReq) (*tenant.CreateTenantRes, error) {

}