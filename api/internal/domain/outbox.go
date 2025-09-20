package domain

import (
	"database/sql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"microservice/internal/model"
	"time"
)

type (
	OutboxStatus string

	Outbox struct {
		Base
		eventType string
		messageId uint
		payload   datatypes.JSON
		status    string
		retries   int
		createdAt time.Time
		retryAt   time.Time
		deletedAt time.Time
	}

	OutboxList struct {
		BaseList
		list []Outbox
	}
)

const (
	OutboxPending    OutboxStatus = "pending"
	OutboxPublishing OutboxStatus = "publishing"
	OutboxPublished  OutboxStatus = "published"
	OutboxFailed     OutboxStatus = "failed"
)

func NewOutbox() *Outbox {
	return &Outbox{}
}

func (o *Outbox) EventType() string {
	return o.eventType
}

func (o *Outbox) SetEventType(event string) {
	o.eventType = event
}

func (o *Outbox) MessageId() uint {
	return o.messageId
}

func (o *Outbox) SetMessageId(messageId uint) {
	o.messageId = messageId
}

func (o *Outbox) Payload() datatypes.JSON {
	return o.payload
}

func (o *Outbox) SetPayload(payload datatypes.JSON) {
	o.payload = payload
}

func (o *Outbox) Status() string {
	return o.status
}

func (o *Outbox) SetStatus(status OutboxStatus) {
	o.status = string(status)
}

func (o *Outbox) Retries() int {
	return o.retries
}

func (o *Outbox) SetRetries(retries int) {
	o.retries = retries
}

func (o *Outbox) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Outbox) SetCreatedAt(createdAt time.Time) {
	o.createdAt = createdAt
}

func (o *Outbox) RetryAt() time.Time {
	return o.retryAt
}

func (o *Outbox) SetRetryAt(retryAt time.Time) {
	o.retryAt = retryAt
}

func (o *Outbox) DeletedAt() time.Time {
	return o.deletedAt
}

func (o *Outbox) SetDeletedAt(deletedAt time.Time) {
	o.deletedAt = deletedAt
}

//

func (o *Outbox) FromDB(src model.Outboxes) Outbox {
	// base
	o.SetID(src.ID)
	o.SetUUID(src.Uuid)
	o.SetCreatedAt(src.CreatedAt)
	o.SetDeletedAt(src.DeletedAt)
	//fields
	o.SetEventType(src.EventType)
	o.SetMessageId(src.MessageID)
	o.SetPayload(src.Payload)
	o.SetStatus(OutboxStatus(src.Status))
	o.SetRetries(src.Retries)
	o.SetRetryAt(src.RetryAt)

	return *o
}

func (o *Outbox) ToDB() model.Outboxes {
	return model.Outboxes{
		BaseSql: model.BaseSql{
			Model: gorm.Model{
				ID:        o.ID(),
				CreatedAt: o.CreatedAt(),
				UpdatedAt: o.UpdatedAt(),
				DeletedAt: gorm.DeletedAt(
					sql.NullTime{
						Time: o.DeletedAt(),
						Valid: func() bool {
							if o.DeletedAt().IsZero() {
								return false
							}
							return true
						}(),
					},
				),
			},
			Uuid: o.UUID(),
		},
		EventType: o.EventType(),
		MessageID: o.MessageId(),
		Payload:   o.Payload(),
		Status:    o.Status(),
		Retries:   o.Retries(),
		CreatedAt: o.CreatedAt(),
		RetryAt:   o.RetryAt(),
		DeletedAt: o.DeletedAt(),
	}
}

//

func NewOutboxList() *OutboxList { return &OutboxList{} }

func (ul *OutboxList) List() []Outbox { return ul.list }

func (ul *OutboxList) SetList(list []Outbox) { ul.list = list }

func (ul *OutboxList) ListToDB() []model.Outboxes {
	outbox := make([]model.Outboxes, 0)

	if len(ul.list) == 0 {
		return outbox
	}

	for _, item := range ul.list {
		outbox = append(outbox, item.ToDB())
	}

	return outbox
}

func (ul *OutboxList) ListFromDB(src []model.Outboxes) OutboxList {
	ul.list = make([]Outbox, 0)

	total := len(src)
	if ul.total == 0 && total > 0 {
		ul.total = int64(total)
	}

	if ul.total == 0 {
		return *ul
	}

	for _, item := range src {
		ul.list = append(ul.list, NewOutbox().FromDB(item))
	}

	return *ul
}

//

type OutboxListReqQryParam struct {
	ReqBaseQryParam
}

func NewOutboxListReqQryParam() *OutboxListReqQryParam {
	return &OutboxListReqQryParam{}
}
