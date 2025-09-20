package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"microservice/internal/model"
)

type (
	MessageStatus string
	Message       struct {
		Base
		channel     string
		tenantId    uint
		mobile      string
		messageText string
		messageHash string
		status      string
		outbox      Outbox
	}

	MessageList struct {
		BaseList
		list []Message
	}
)

const (
	MsgQueued    MessageStatus = "queued"
	MsgSending   MessageStatus = "sending"
	MsgSent      MessageStatus = "sent"
	MsgDelivered MessageStatus = "delivered"
	MsgFailed    MessageStatus = "failed"
)

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) Channel() string {
	return m.channel
}

func (m *Message) SetChannel(channel string) {
	m.channel = channel
}

func (m *Message) TenantID() uint {
	return m.tenantId
}

func (m *Message) SetTenantID(tenantId uint) {
	m.tenantId = tenantId
}

func (m *Message) Mobile() string {
	return m.mobile
}

func (m *Message) SetMobile(mobile string) {
	m.mobile = mobile
}

func (m *Message) MessageText() string {
	return m.messageText
}

func (m *Message) SetMessageText(messageText string) {
	m.messageText = messageText
}

func (m *Message) MessageHash() string {
	return m.messageHash
}

func (m *Message) SetMessageHash(messageHash string) {
	m.messageHash = messageHash
}

func (m *Message) Status() string {
	return m.status
}

func (m *Message) SetStatus(status string) {
	m.status = status
}

//

func (m *Message) FromDB(src model.Messages) Message {
	// base
	m.SetID(src.ID)
	m.SetUUID(src.Uuid)
	m.SetCreatedAt(src.CreatedAt)
	m.SetUpdatedAt(src.UpdatedAt)
	m.SetDeletedAt(src.DeletedAt.Time)
	//fields
	m.SetTenantID(src.TenantID)
	m.SetMobile(src.Mobile)
	m.SetMessageText(src.MessageText)
	m.SetMessageHash(src.MessageHash)
	m.SetStatus(src.Status)

	if src.Outbox.ID != 0 {
		m.SetChannel(src.Outbox.EventType)
	}

	return *m
}

func (m *Message) ToDB() model.Messages {
	return model.Messages{
		BaseSql: model.BaseSql{
			Model: gorm.Model{
				ID:        m.ID(),
				CreatedAt: m.CreatedAt(),
				UpdatedAt: m.UpdatedAt(),
				DeletedAt: gorm.DeletedAt(
					sql.NullTime{
						Time: m.DeletedAt(),
						Valid: func() bool {
							if m.DeletedAt().IsZero() {
								return false
							}
							return true
						}(),
					},
				),
			},
			Uuid: m.UUID(),
		},
		TenantID:    m.TenantID(),
		Mobile:      m.Mobile(),
		MessageText: m.MessageText(),
		MessageHash: m.MessageHash(),
		Status:      m.Status(),
	}
}

//

func NewMessageList() *MessageList { return &MessageList{} }

func (ul *MessageList) List() []Message { return ul.list }

func (ul *MessageList) SetList(list []Message) { ul.list = list }

func (ul *MessageList) ListToDB() []model.Messages {
	message := make([]model.Messages, 0)

	if len(ul.list) == 0 {
		return message
	}

	for _, item := range ul.list {
		message = append(message, item.ToDB())
	}

	return message
}

func (ul *MessageList) ListFromDB(src []model.Messages) MessageList {
	ul.list = make([]Message, 0)

	total := len(src)
	if ul.total == 0 && total > 0 {
		ul.total = int64(total)
	}

	if ul.total == 0 {
		return *ul
	}

	for _, item := range src {
		ul.list = append(ul.list, NewMessage().FromDB(item))
	}

	return *ul
}

//

type MessageListReqQryParam struct {
	ReqBaseQryParam
	tenantId uint
}

func NewMessageListReqQryParam() *MessageListReqQryParam {
	return &MessageListReqQryParam{}
}

func (m *MessageListReqQryParam) TenantId() uint {
	return m.tenantId
}

func (m *MessageListReqQryParam) SetTenantId(tenantId uint) {
	m.tenantId = tenantId
}
