package port

import (
	"context"
	"microservice/internal/domain"
)

type (
	IMessageRepository interface {
		Create(ctx context.Context, ent domain.Message) (domain.Message, error)
		Update(ctx context.Context, ent domain.Message) error
		UpdateStatus(ctx context.Context, id uint, status string) error
		GetList(ctx context.Context, ent domain.MessageListReqQryParam) (domain.MessageList, error)
	}

	IMessageUsecase interface {
		Send(ctx context.Context, credit domain.Tenant, ent domain.Message) (domain.Message, error)
		GetList(ctx context.Context, ent domain.MessageListReqQryParam) (domain.MessageList, error)
	}
)
