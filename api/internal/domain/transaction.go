package domain

import (
	"microservice/internal/model"
	"time"
)

type (
	Transaction struct {
		id            []byte
		creditId      uint
		amount        float64
		messageHashId []byte
		createdAt     time.Time
	}

	TransactionList struct {
		BaseList
		list []Transaction
	}
)

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) ID() []byte {
	return t.id
}

func (t *Transaction) SetID(id []byte) {
	t.id = id
}

func (t *Transaction) CreditID() uint {
	return t.creditId
}

func (t *Transaction) SetCreditID(creditId uint) {
	t.creditId = creditId
}

func (t *Transaction) Amount() float64 {
	return t.amount
}

func (t *Transaction) SetAmount(amount float64) {
	t.amount = amount
}

func (t *Transaction) MessageHashID() []byte {
	return t.messageHashId
}

func (t *Transaction) SetMessageHashID(messageHashId []byte) {
	t.messageHashId = messageHashId
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transaction) SetCreatedAt(createdAt time.Time) {
	t.createdAt = createdAt
}

//

func (t *Transaction) FromDB(src model.CreditTransactions) Transaction {
	// base
	t.SetID(src.ID)
	t.SetCreatedAt(src.CreatedAt)
	//fields
	t.SetCreditID(src.CreditID)
	t.SetAmount(src.Amount)

	if src.MessageHashID != nil {
		t.SetMessageHashID(src.MessageHashID)
	}

	return *t
}

func (t *Transaction) ToDB() model.CreditTransactions {
	return model.CreditTransactions{
		ID:            t.ID(),
		CreditID:      t.CreditID(),
		Amount:        t.Amount(),
		MessageHashID: t.MessageHashID(),
		CreatedAt:     t.CreatedAt(),
	}
}

//

func NewTransactionList() *TransactionList { return &TransactionList{} }

func (ul *TransactionList) List() []Transaction { return ul.list }

func (ul *TransactionList) SetList(list []Transaction) { ul.list = list }

func (ul *TransactionList) ListToDB() []model.CreditTransactions {
	transaction := make([]model.CreditTransactions, 0)

	if len(ul.list) == 0 {
		return transaction
	}

	for _, item := range ul.list {
		transaction = append(transaction, item.ToDB())
	}

	return transaction
}

func (ul *TransactionList) ListFromDB(src []model.CreditTransactions) TransactionList {
	ul.list = make([]Transaction, 0)

	total := len(src)
	if ul.total == 0 && total > 0 {
		ul.total = int64(total)
	}

	if ul.total == 0 {
		return *ul
	}

	for _, item := range src {
		ul.list = append(ul.list, NewTransaction().FromDB(item))
	}

	return *ul
}

//

type TransactionListReqQryParam struct {
	ReqBaseQryParam
}

func NewTransactionListReqQryParam() *TransactionListReqQryParam {
	return &TransactionListReqQryParam{}
}
