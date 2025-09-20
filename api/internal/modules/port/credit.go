package port

import (
	"context"
	"microservice/internal/domain"
)

type (
	ICreditRepository interface {
		Create(ctx context.Context, ent domain.Credit) (domain.Credit, error)
		GetDetails(ctx context.Context, ent domain.Credit) (domain.Credit, error)
		Update(ctx context.Context, ent domain.Credit) error
	}

	ICreditUsecase interface {
		IncreaseAmount(ctx context.Context, ent domain.Credit) (domain.Credit, error)
		GetDetails(ctx context.Context, ent domain.Credit, qp domain.TransactionListReqQryParam) (domain.Credit, error)
	}
)
