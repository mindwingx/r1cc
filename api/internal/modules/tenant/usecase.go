package tenant

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/internal/adapter/cache"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/trace"
	"microservice/internal/domain"
	"microservice/internal/modules/port"
	"microservice/pkg/meta"
)

type (
	UsecaseFx struct {
		fx.In
		Locale     locale.ILocale
		Tracer     trace.ITracer
		Logger     logger.ILogger
		Cache      cache.ICache
		Tx         orm.ISqlTx
		TenantRepo port.ITenantRepository
		CreditRepo port.ICreditRepository
	}

	Usecase struct {
		l          locale.ILocale
		trc        trace.ITracer
		lgr        logger.ILogger
		cache      cache.ICache
		tx         orm.ISqlTx
		tenantRepo port.ITenantRepository
		creditRepo port.ICreditRepository
	}
)

func NewUsecaseFx(fx UsecaseFx) port.ITenantUsecase {
	return &Usecase{
		l:          fx.Locale,
		trc:        fx.Tracer,
		lgr:        fx.Logger,
		cache:      fx.Cache,
		tx:         fx.Tx,
		tenantRepo: fx.TenantRepo,
		creditRepo: fx.CreditRepo,
	}
}

func (uc *Usecase) Create(ctx context.Context, ent domain.Tenant) (res domain.Tenant, err error) {
	var txErr error

	uc.tx.Begin()
	defer func() {
		if r := recover(); r != nil {
			txErr = r.(error)
			uc.lgr.Error("tenant.create.recover", zap.Error(txErr))
			err = meta.Failed
		}

		if txResErr := uc.tx.Resolve(txErr); txResErr != nil {
			uc.lgr.Error("tenant.create.tx.resolve", zap.Error(txErr))
		}
	}()

	tenant, txErr := uc.tenantRepo.Create(ctx, ent)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	{
		credit := domain.NewCredit()
		credit.SetTenantID(tenant.ID())
		credit.SetBalance(0)

		_, txErr = uc.creditRepo.Create(ctx, *credit)
		if txErr != nil {
			err = meta.EvalTxErr(txErr)
			return
		}
	}

	res = tenant
	return
}

func (uc *Usecase) GetDetails(ctx context.Context, ent domain.Tenant) (res domain.Tenant, err error) {
	ent.SetRelations("Credit")

	res, txErr := uc.tenantRepo.GetDetails(ctx, ent)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	return
}

func (uc *Usecase) GetList(ctx context.Context, ent domain.TenantListReqQryParam) (res domain.TenantList, err error) {
	res, txErr := uc.tenantRepo.GetList(ctx, ent)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	return
}
