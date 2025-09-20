package utils

import (
	"crypto/sha256"
	"fmt"
	"microservice/internal/domain"
)

func TransactionIdGen(t domain.Tenant, c domain.Credit, msg *domain.Message) (id []byte) {
	transaction := c.TxAmount()
	if msg == nil {
		msg = domain.NewMessage()
	}

	data := fmt.Sprintf("%d:%s:%x:%x:%s", t.UUID(), c.UUID(), c.Balance(), transaction.Amount(), msg.MessageText())
	sum := sha256.Sum256([]byte(data))
	id = sum[:]
	return
}
