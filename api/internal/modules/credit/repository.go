package credit

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func NewRepositoryFx(fx RepositoryFx) port.ICreditRepository {
	return &Repository{
		l:   fx.Locale,
		trc: fx.Tracer,
		lgr: fx.Logger,
		sql: fx.Sql,
	}
}

func (r *Repository) Create(ctx context.Context, ent domain.Credit) (res domain.Credit, err error) {
	m := ent.ToDB()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Credits{})

	txErr := tx.Omit("uuid", "created_at", "deleted_at").Clauses(clause.Returning{}).Create(&m).Error
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

func (r *Repository) GetDetails(ctx context.Context, ent domain.Credit) (res domain.Credit, err error) {
	m := model.NewCredit()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Tenants{})

	if ent.GetRelations() != nil {
		for _, rel := range ent.GetRelations() {
			tx = tx.Preload(rel)
		}
	}

	u := tx.First(&m, "uuid = ?", ent.UUID())

	if err = u.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.lgr.Error("credit.repo.detail", zap.Error(err))
		err = meta.Failed
		return
	}

	if u.RowsAffected == 0 {
		err = meta.NotFound
		return
	}

	res = *domain.NewCredit()
	res.FromDB(*m)
	return
}

func (r *Repository) Update(ctx context.Context, ent domain.Credit) (err error) {
	m := ent.ToDB()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Credits{}).Unscoped().
		Omit("created_at", "deleted_at").
		Clauses(clause.Locking{Strength: "UPDATE"}). //locking the row to avoid race condition
		Where("uuid = ?", ent.UUID()).Updates(m)

	if err = tx.Error; err != nil {
		r.lgr.Error("credit.repo.update", zap.Error(err))

		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)

		if pgErr != nil && pgErr.Code == "23505" { // PSQL Unique violation error code
			err = meta.ItemExist.SetErr(pgErr.Detail)
			return
		}

		err = meta.Failed
		return
	}

	return
}
