package model

type Tenants struct {
	BaseSql
	Username   string  `json:"username"`
	TenantName string  `json:"tenant_name"`
	Active     bool    `json:"active"`
	Credit     Credits `json:"credit,omitempty" gorm:"foreignKey:TenantID"`
}

func NewTenant() *Tenants { return &Tenants{} }

func (m *Tenants) TableName() string { return "tenants" }
