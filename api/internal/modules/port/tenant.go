package port

import (
	"context"
	"microservice/internal/domain"
)

type (
	ITenantRepository interface {
		Create(ctx context.Context, ent domain.Tenant) (domain.Tenant, error)
		GetDetails(ctx context.Context, ent domain.Tenant) (domain.Tenant, error)
		GetList(ctx context.Context, ent domain.TenantListReqQryParam) (domain.TenantList, error)
	}

	ITenantUsecase interface {
		Create(ctx context.Context, ent domain.Tenant) (domain.Tenant, error)
		GetDetails(ctx context.Context, ent domain.Tenant) (domain.Tenant, error)
		GetList(ctx context.Context, ent domain.TenantListReqQryParam) (domain.TenantList, error)
	}
)
