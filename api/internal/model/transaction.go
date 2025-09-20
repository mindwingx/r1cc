package model

import "time"

type CreditTransactions struct {
	ID            []byte    `json:"id" gorm:"type:VARCHAR(64);primaryKey;default:null"`
	CreditID      uint      `json:"credit_id"`
	Amount        float64   `json:"amount"`
	MessageHashID []byte    `json:"message_hash_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewTransaction() *CreditTransactions { return &CreditTransactions{} }

func (m *CreditTransactions) TableName() string { return "credit_transactions" }
