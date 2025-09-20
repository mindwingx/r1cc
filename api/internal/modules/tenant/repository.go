package tenant

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

func NewRepositoryFx(fx RepositoryFx) port.ITenantRepository {
	return &Repository{
		l:   fx.Locale,
		trc: fx.Tracer,
		lgr: fx.Logger,
		sql: fx.Sql,
	}
}

func (r *Repository) Create(ctx context.Context, ent domain.Tenant) (res domain.Tenant, err error) {
	m := ent.ToDB()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Tenants{})

	txErr := tx.Omit("uuid", "active").Clauses(clause.Returning{}).Create(&m).Error
	if txErr != nil {
		err = txErr
		r.lgr.Error("tenant.repo.create", zap.Error(err))

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

func (r *Repository) GetDetails(ctx context.Context, ent domain.Tenant) (res domain.Tenant, err error) {
	m := model.NewTenant()

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Tenants{})

	if ent.GetRelations() != nil {
		for _, rel := range ent.GetRelations() {
			tx = tx.Preload(rel, func(db *gorm.DB) *gorm.DB {
				return db.Unscoped()
			})
		}
	}

	u := tx.First(&m, "uuid = ?", ent.UUID())

	if err = u.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.lgr.Error("tenant.repo.detail", zap.Error(err))
		err = meta.Failed
		return
	}

	if u.RowsAffected == 0 {
		err = meta.NotFound
		return
	}

	res = *domain.NewTenant()
	res.FromDB(*m)
	return
}

func (r *Repository) GetList(ctx context.Context, ent domain.TenantListReqQryParam) (res domain.TenantList, err error) {
	defer func() {
		if err != nil {
			r.lgr.Error("tenant.repo.list", zap.Error(err))
		}
	}()

	list := domain.NewTenantList()

	var (
		models []model.Tenants
		total  int64
	)

	db := r.sql.Tx()
	tx := db.WithContext(ctx).Model(&model.Tenants{})

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
		tx.Where("username ILIKE ? OR tenant_name ILIKE ?", val, val)
	}

	//

	if err = tx.Count(&total).Error; err != nil {
		r.lgr.Error("tenant.repo.list.count", zap.Error(err))
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
		r.lgr.Error("tenant.repo.list", zap.Error(err))
		err = meta.Failed
		return
	}

	if items.RowsAffected > 0 {
		list.ListFromDB(models)
	}

	res = *list
	return
}
