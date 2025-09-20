package domain

import "encoding/json"

type OutboxMessage struct {
	TenantId    uint   `json:"tenantId"`
	MessageId   uint   `json:"messageId"`
	OutboxId    uint   `json:"outboxId"`
	Channel     string `json:"channel"`
	Mobile      string `json:"mobile"`
	MessageText string `json:"messageText"`
	MessageHash string `json:"messageHash"`
	Status      string `json:"status"`
}

func NewOutboxMessage() *OutboxMessage {
	return &OutboxMessage{}
}

func (om *OutboxMessage) FromMessage(msg Message) {
	om.TenantId = msg.TenantID()
	om.MessageId = msg.ID()
	om.Channel = msg.Channel()
	om.Mobile = msg.Mobile()
	om.MessageText = msg.MessageText()
	om.MessageHash = msg.MessageHash()
	om.Status = msg.Status()
}

func (om *OutboxMessage) SetOutboxID(id uint) {
	om.OutboxId = id
}
func (om *OutboxMessage) Json() []byte {
	payload, _ := json.Marshal(om)
	return payload
}
