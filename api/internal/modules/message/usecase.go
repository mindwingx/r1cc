package message

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/internal/adapter/cache"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/queue"
	"microservice/internal/adapter/trace"
	"microservice/internal/domain"
	"microservice/internal/modules/port"
	"microservice/pkg/meta"
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
		MessageRepo     port.IMessageRepository
		OutboxRepo      port.IOutboxRepository
		CreditRepo      port.ICreditRepository
		TransactionRepo port.ITransactionRepository
		Queue           queue.IQueue
	}

	Usecase struct {
		l               locale.ILocale
		trc             trace.ITracer
		lgr             logger.ILogger
		cache           cache.ICache
		tx              orm.ISqlTx
		messageRepo     port.IMessageRepository
		outboxRepo      port.IOutboxRepository
		creditRepo      port.ICreditRepository
		transactionRepo port.ITransactionRepository
		queue           queue.IQueue
	}
)

func NewUsecaseFx(fx UsecaseFx) port.IMessageUsecase {
	return &Usecase{
		l:               fx.Locale,
		trc:             fx.Tracer,
		lgr:             fx.Logger,
		cache:           fx.Cache,
		tx:              fx.Tx,
		messageRepo:     fx.MessageRepo,
		outboxRepo:      fx.OutboxRepo,
		creditRepo:      fx.CreditRepo,
		transactionRepo: fx.TransactionRepo,
		queue:           fx.Queue,
	}
}

const MciMessagePrice float64 = 8.9

func (uc *Usecase) Send(ctx context.Context, tenant domain.Tenant, ent domain.Message) (res domain.Message, err error) {
	var txErr error

	if len(ent.MessageText()) > 160 {
		err = meta.Conflict.SetErr(uc.l.Get("sms_char_exceed"))
		return
	}

	credit := tenant.Credit()

	if credit.Balance() < MciMessagePrice {
		err = meta.Conflict.SetErr(uc.l.Get("sms_balance_err"))
		return
	}

	//

	uc.tx.Begin()
	defer func() {
		if r := recover(); r != nil {
			txErr = r.(error)
			uc.lgr.Error("message.create.recover", zap.Error(txErr))
			err = meta.Failed
		}

		if txResErr := uc.tx.Resolve(txErr); txResErr != nil {
			uc.lgr.Error("message.create.tx.rollback", zap.Error(txErr))
		}
	}()

	//

	hashedMessage := messageHashedIdGen(ent)
	ent.SetMessageHash(hashedMessage)

	message, txErr := uc.messageRepo.Create(ctx, ent)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	//

	om := domain.NewOutboxMessage()
	om.FromMessage(message)

	outboxEnt := domain.NewOutbox()
	outboxEnt.SetEventType(message.Channel())
	outboxEnt.SetMessageId(message.ID())
	outboxEnt.SetPayload(om.Json())

	outbox, txErr := uc.outboxRepo.Create(ctx, *outboxEnt)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	om.SetOutboxID(outbox.ID())

	//

	decreasedCredit := credit.Balance() - MciMessagePrice
	credit.SetBalance(decreasedCredit)

	txErr = uc.creditRepo.Update(ctx, credit)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	//

	transaction := domain.NewTransaction()
	transaction.SetID(utils.TransactionIdGen(tenant, credit, &ent))
	transaction.SetCreditID(credit.ID())
	transaction.SetAmount(MciMessagePrice)
	transaction.SetMessageHashID([]byte(hashedMessage))

	_, txErr = uc.transactionRepo.Create(ctx, *transaction)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	//

	if txErr = uc.tx.Commit(); txErr != nil {
		uc.lgr.Error("message.create.tx.commit", zap.Error(txErr))
		err = meta.Failed
		return
	}

	//

	txErr = uc.queue.Produce(ctx, ent.Channel(), hashedMessage, om.Json())
	if txErr != nil {
		uc.lgr.Error("message.create.queue.produce", zap.Error(txErr))
		return
	}

	res = message
	return
}

func (uc *Usecase) GetList(ctx context.Context, ent domain.MessageListReqQryParam) (res domain.MessageList, err error) {
	ent.SetRelations("Outbox")
	res, txErr := uc.messageRepo.GetList(ctx, ent)
	if txErr != nil {
		err = meta.EvalTxErr(txErr)
		return
	}

	return
}

// HELPERS

func messageHashedIdGen(msg domain.Message) string {
	id := fmt.Sprintf("%d:%s:%s", msg.TenantID(), msg.Mobile(), msg.MessageText())
	h := sha256.Sum256([]byte(id))
	return hex.EncodeToString(h[:])
}
