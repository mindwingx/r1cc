package model

type Credits struct {
	BaseSql
	TenantID     uint                 `json:"tenant_id"`
	Balance      float64              `json:"balance"`
	Transactions []CreditTransactions `json:"transactions" gorm:"foreignKey:CreditID"`
}

func NewCredit() *Credits { return &Credits{} }

func (m *Credits) TableName() string { return "credits" }
