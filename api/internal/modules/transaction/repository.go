package transaction

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/trace"
	"microservice/internal/domain"
	"microservice/internal/model"
	"microservice/internal/modules/port"
	"microservice/pkg/meta"
)

type (
	RepositoryFx struct {
		fx.In
		Locale locale.ILocale
		Tracer trace.ITracer
		Logger logger.ILogger
		Sql    orm.ISqlTx
	}

	Repository struct {
		l   locale.ILocale
		trc trace.ITracer
		lgr logger.ILogger
		sql orm.ISqlTx
	}
)

func NewRepositoryFx(fx RepositoryFx) port.ITransactionRepository {
	return &Repository{
		l:   fx.Locale,
		trc: fx.Tracer,
		lgr: fx.Logger,
		sql: fx.Sql,
	}
}

func (r *Repository) Create(ctx context.Context, ent domain.Transaction) (res domain.Transaction, err error) {
	m := ent.ToDB()
	columns := []string{"updated_at", "deleted_at"}

	if m.MessageHashID == nil {
		columns = append(columns, "message_hash_id")
	}

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.CreditTransactions{})

	txErr := tx.Omit(columns...).Clauses(clause.Returning{}).Create(&m).Error
	if txErr != nil {
		err = txErr
		r.lgr.Error("credit.repo.create", zap.Error(err))

		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)

		if pgErr != nil && pgErr.Code == "23505" { // PSQL Unique violation error code
			err = meta.ItemExist.SetErr(pgErr.Detail)
			return
		}

		err = meta.Failed
		return
	}

	res = ent.FromDB(m)
	return
}

func (r *Repository) GetList(ctx context.Context, ent domain.TransactionListReqQryParam) (res domain.TransactionList, err error) {
	defer func() {
		if err != nil {
			r.lgr.Error("transaction.repo.list", zap.Error(err))
		}
	}()

	list := domain.NewTransactionList()

	var (
		models []model.CreditTransactions
		total  int64
	)

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.CreditTransactions{})

	if ent.RelId() != 0 {
		tx.Where("credit_id = ?", ent.RelId()) // get all items
	}

	if ent.Items() != nil && len(ent.Items()) > 0 {
		tx.Where("uuid IN ?", ent.Items()) // get all items
	}

	//

	if err = tx.Count(&total).Error; err != nil {
		r.lgr.Error("transaction.repo.list.count", zap.Error(err))
		err = meta.Failed
		return
	}

	list.SetTotal(total)

	//

	if ent.Items() == nil {
		tx.Offset(ent.Offset()).Limit(ent.Limit())
	}

	items := tx.Order(ent.SortOrder()).Find(&models)
	if err = items.Error; err != nil {
		r.lgr.Error("transaction.repo.list", zap.Error(err))
		err = meta.Failed
		return
	}

	if items.RowsAffected > 0 {
		list.ListFromDB(models)
	}

	res = *list
	return
}
