package rbac

import (
	"context"
	"github.com/google/uuid"
	"microservice/internal/domain"
	"microservice/pkg/meta"
)

func CtxTenant(ctx context.Context) (t domain.Tenant, err error) {
	userUuid, ok := ctx.Value("X.TENANT.UUID").(string)
	if !ok {
		err = meta.Failed
		return
	}

	uid, err := uuid.Parse(userUuid)

	if err != nil {
		return
	}

	t = *domain.NewTenant()
	t.SetUUID(uid)
	return
}
