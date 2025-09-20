package model

import (
	"gorm.io/datatypes"
	"time"
)

type Outboxes struct {
	BaseSql
	EventType string         `json:"event_type"`
	MessageID uint           `json:"message_id"`
	Payload   datatypes.JSON `json:"payload"`
	Status    string         `json:"status"`
	Retries   int            `json:"retries"`
	CreatedAt time.Time      `json:"created_at"`
	RetryAt   time.Time      `json:"retry_at"`
	DeletedAt time.Time      `json:"deleted_at"`
}

func NewOutbox() *Outboxes { return &Outboxes{} }

func (m *Outboxes) TableName() string { return "outboxes" }
