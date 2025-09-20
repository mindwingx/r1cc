package port

import (
	"context"
	"microservice/internal/domain"
)

type (
	ITransactionRepository interface {
		Create(ctx context.Context, ent domain.Transaction) (domain.Transaction, error)
		GetList(ctx context.Context, ent domain.TransactionListReqQryParam) (domain.TransactionList, error)
	}

	ITransactionUsecase interface {
		//Create(ctx context.Context, ent domain.Transaction) (domain.Transaction, error)
		//GetList(ctx context.Context, ent domain.TransactionListReqQryParam) (domain.TransactionList, error)
	}
)
