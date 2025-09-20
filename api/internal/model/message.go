package model

type Messages struct {
	BaseSql
	TenantID    uint     `json:"tenant_id"`
	Mobile      string   `json:"mobile"`
	MessageText string   `json:"message_text"`
	MessageHash string   `json:"message_hash"`
	Status      string   `json:"status"`
	Outbox      Outboxes `json:"outbox,omitempty" gorm:"foreignKey:MessageID"`
}

func NewMessage() *Messages { return &Messages{} }

func (m *Messages) TableName() string { return "messages" }
