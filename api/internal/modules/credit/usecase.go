package credit

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
	"microservice/pkg/rbac"
	"microservice/pkg/utils"
)

type (
	UsecaseFx struct {
		fx.In
		Locale          locale.ILocale
		Tracer          trace.ITracer
		Logger          logger.ILogger
		Cache           cache.ICache
		Tx              orm.ISqlTx
		CreditRepo      port.ICreditRepository
		TransactionRepo port.ITransactionRepository
	}

	Usecase struct {
		l               locale.ILocale
		trc             trace.ITracer
		lgr             logger.ILogger
		cache           cache.ICache
		tx              orm.ISqlTx
		creditRepo      port.ICreditRepository
		transactionRepo port.ITransactionRepository
	}
)

func NewUsecaseFx(fx UsecaseFx) port.ICreditUsecase {
	return &Usecase{
		l:               fx.Locale,
		trc:             fx.Tracer,
		lgr:             fx.Logger,
		cache:           fx.Cache,
		tx:              fx.Tx,
		creditRepo:      fx.CreditRepo,
		transactionRepo: fx.TransactionRepo,
	}
}

func (uc *Usecase) IncreaseAmount(ctx context.Context, ent domain.Credit) (res domain.Credit, err error) {
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

	//

	tenant, txErr := rbac.CtxTenant(ctx)
	if txErr != nil {
		err = meta.Failed.SetErr(uc.l.Get("invalid_client"))
		return
	}

	//

	transaction := ent.TxAmount()
	transaction.SetID(utils.TransactionIdGen(tenant, ent, nil))

	transactionRes, txErr := uc.transactionRepo.Create(ctx, transaction)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	//

	balance := ent.Balance()
	balance += transactionRes.Amount()
	ent.SetBalance(balance)

	err = uc.creditRepo.Update(ctx, ent)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	//

	ent.SetTxAmount(transactionRes)
	res = ent
	return
}

func (uc *Usecase) GetDetails(ctx context.Context, ent domain.Credit, qp domain.TransactionListReqQryParam) (res domain.Credit, err error) {
	qp.SetRelId(ent.ID())
	txs, txErr := uc.transactionRepo.GetList(ctx, qp)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	ent.SetTransactions(txs.List())
	res = ent
	return
}
