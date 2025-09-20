package outbox

import (
	"context"
	"go.uber.org/fx"
	"microservice/internal/adapter/cache"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/trace"
	"microservice/internal/domain"
	"microservice/internal/modules/port"
)

type (
	UsecaseFx struct {
		fx.In
		Locale     locale.ILocale
		Tracer     trace.ITracer
		Logger     logger.ILogger
		Cache      cache.ICache
		Tx         orm.ISqlTx
		OutboxRepo port.IOutboxRepository
	}

	Usecase struct {
		l          locale.ILocale
		trc        trace.ITracer
		lgr        logger.ILogger
		cache      cache.ICache
		tx         orm.ISqlTx
		outboxRepo port.IOutboxRepository
	}
)

func NewUsecaseFx(fx UsecaseFx) port.IOutboxUsecase {
	return &Usecase{
		l:          fx.Locale,
		trc:        fx.Tracer,
		lgr:        fx.Logger,
		cache:      fx.Cache,
		tx:         fx.Tx,
		outboxRepo: fx.OutboxRepo,
	}
}

func (u *Usecase) Create(ctx context.Context, ent domain.Outbox) (domain.Outbox, error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) GetDetails(ctx context.Context, ent domain.Outbox) (domain.Outbox, error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) Update(ctx context.Context, ent domain.Outbox) error {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) Delete(ctx context.Context, ent domain.Outbox) error {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) GetList(ctx context.Context, ent domain.OutboxListReqQryParam) (domain.OutboxList, error) {
	//TODO implement me
	panic("implement me")
}
