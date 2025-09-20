package message

import (
	"context"
	"errors"
	"fmt"
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

func NewRepositoryFx(fx RepositoryFx) port.IMessageRepository {
	return &Repository{
		l:   fx.Locale,
		trc: fx.Tracer,
		lgr: fx.Logger,
		sql: fx.Sql,
	}
}

func (r *Repository) Create(ctx context.Context, ent domain.Message) (res domain.Message, err error) {
	m := ent.ToDB()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Messages{})

	txErr := tx.Omit("uuid", "status").Clauses(clause.Returning{}).Create(&m).Error
	if txErr != nil {
		err = txErr
		r.lgr.Error("message.repo.create", zap.Error(err))

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

func (r *Repository) GetDetails(ctx context.Context, ent domain.Message) (res domain.Message, err error) {
	m := model.NewMessage()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Messages{})

	if ent.GetRelations() != nil {
		for _, rel := range ent.GetRelations() {
			tx = tx.Preload(rel)
		}
	}

	u := tx.First(&m, "uuid = ?", ent.UUID())

	if err = u.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.lgr.Error("message.repo.detail", zap.Error(err))
		err = meta.Failed
		return
	}

	if u.RowsAffected == 0 {
		err = meta.NotFound
		return
	}

	res = *domain.NewMessage()
	res.FromDB(*m)
	return
}

func (r *Repository) Update(ctx context.Context, ent domain.Message) (err error) {
	m := ent.ToDB()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Messages{}).Where("uuid = ?", ent.UUID()).Updates(m)
	if err = tx.Error; err != nil {
		r.lgr.Error("message.repo.update", zap.Error(err))

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

func (r *Repository) UpdateStatus(ctx context.Context, id uint, status string) (err error) {
	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Messages{}).Where("id = ?", id).Update("status", status)
	if err = tx.Error; err != nil {
		r.lgr.Error("message.repo.update", zap.Error(err))

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

func (r *Repository) GetList(ctx context.Context, ent domain.MessageListReqQryParam) (res domain.MessageList, err error) {
	defer func() {
		if err != nil {
			r.lgr.Error("message.repo.list", zap.Error(err))
		}
	}()

	list := domain.NewMessageList()

	var (
		models []model.Messages
		total  int64
	)

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Messages{})

	tx.Where("tenant_id = ?", ent.TenantId())

	if ent.GetRelations() != nil {
		for _, rel := range ent.GetRelations() {
			tx = tx.Preload(rel)
		}
	}

	if ent.Items() != nil && len(ent.Items()) > 0 {
		tx.Where("uuid IN ?", ent.Items()) // get all items
	}

	if len(ent.Search()) > 0 {
		val := fmt.Sprintf("%%%s%%", ent.Search()) // this returns %search_value%
		tx.Where("message_text ILIKE ? ", val)
	}

	//

	if err = tx.Count(&total).Error; err != nil {
		r.lgr.Error("message.repo.list.count", zap.Error(err))
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
		r.lgr.Error("message.repo.list", zap.Error(err))
		err = meta.Failed
		return
	}

	if items.RowsAffected > 0 {
		list.ListFromDB(models)
	}

	res = *list
	return
}
