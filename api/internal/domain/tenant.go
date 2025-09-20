package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"microservice/internal/model"
)

type (
	Tenant struct {
		Base
		username   string
		tenantName string
		active     bool
		credit     Credit
	}

	TenantList struct {
		BaseList
		list []Tenant
	}
)

func NewTenant() *Tenant {
	return &Tenant{}
}

func (t *Tenant) Username() string {
	return t.username
}

func (t *Tenant) SetUsername(username string) {
	t.username = username
}

func (t *Tenant) TenantName() string {
	return t.tenantName
}

func (t *Tenant) SetTenantName(name string) {
	t.tenantName = name
}

func (t *Tenant) Active() bool {
	return t.active
}

func (t *Tenant) SetActive(active bool) {
	t.active = active
}

//

func (t *Tenant) Credit() Credit {
	return t.credit
}

func (t *Tenant) SetCredit(credit Credit) {
	t.credit = credit
}

//

func (t *Tenant) FromDB(src model.Tenants) Tenant {
	// base
	t.SetID(src.ID)
	t.SetUUID(src.Uuid)
	t.SetCreatedAt(src.CreatedAt)
	t.SetUpdatedAt(src.UpdatedAt)
	t.SetDeletedAt(src.DeletedAt.Time)
	//fields
	t.SetUsername(src.Username)
	t.SetTenantName(src.TenantName)
	t.SetActive(src.Active)
	// relations
	if src.Credit.ID != 0 {
		c := NewCredit().FromDB(src.Credit)
		t.SetCredit(c)
	}

	return *t
}

func (t *Tenant) ToDB() model.Tenants {
	return model.Tenants{
		BaseSql: model.BaseSql{
			Model: gorm.Model{
				ID:        t.ID(),
				CreatedAt: t.CreatedAt(),
				UpdatedAt: t.UpdatedAt(),
				DeletedAt: gorm.DeletedAt(
					sql.NullTime{
						Time: t.DeletedAt(),
						Valid: func() bool {
							if t.DeletedAt().IsZero() {
								return false
							}
							return true
						}(),
					},
				),
			},
			Uuid: t.UUID(),
		},
		Username:   t.Username(),
		TenantName: t.TenantName(),
		Active:     t.Active(),
	}
}

//

func NewTenantList() *TenantList { return &TenantList{} }

func (ul *TenantList) List() []Tenant { return ul.list }

func (ul *TenantList) SetList(list []Tenant) { ul.list = list }

func (ul *TenantList) ListToDB() []model.Tenants {
	tenant := make([]model.Tenants, 0)

	if len(ul.list) == 0 {
		return tenant
	}

	for _, item := range ul.list {
		tenant = append(tenant, item.ToDB())
	}

	return tenant
}

func (ul *TenantList) ListFromDB(src []model.Tenants) TenantList {
	ul.list = make([]Tenant, 0)

	total := len(src)
	if ul.total == 0 && total > 0 {
		ul.total = int64(total)
	}

	if ul.total == 0 {
		return *ul
	}

	for _, item := range src {
		ul.list = append(ul.list, NewTenant().FromDB(item))
	}

	return *ul
}

//

type TenantListReqQryParam struct {
	ReqBaseQryParam
}

func NewTenantListReqQryParam() *TenantListReqQryParam {
	return &TenantListReqQryParam{}
}
