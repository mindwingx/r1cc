package port

import (
	"context"
	"microservice/internal/domain"
)

type (
	IOutboxRepository interface {
		Create(ctx context.Context, ent domain.Outbox) (domain.Outbox, error)
		GetDetails(ctx context.Context, ent domain.Outbox) (domain.Outbox, error)
		Update(ctx context.Context, ent domain.Outbox) error
		UpdateStatus(ctx context.Context, id uint, status string) error
		UpdateTryCount(ctx context.Context, id uint, count int) error
		Delete(ctx context.Context, ent domain.Outbox) error
		GetList(ctx context.Context, ent domain.OutboxListReqQryParam) (domain.OutboxList, error)
	}

	IOutboxUsecase interface {
		Create(ctx context.Context, ent domain.Outbox) (domain.Outbox, error)
		GetDetails(ctx context.Context, ent domain.Outbox) (domain.Outbox, error)
		Update(ctx context.Context, ent domain.Outbox) error
		Delete(ctx context.Context, ent domain.Outbox) error
		GetList(ctx context.Context, ent domain.OutboxListReqQryParam) (domain.OutboxList, error)
	}
)
