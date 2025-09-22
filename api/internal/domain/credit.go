package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"microservice/internal/model"
)

type (
	Credit struct {
		Base
		tenantId     uint
		balance      float64
		txAmount     Transaction
		transactions TransactionList
	}

	CreditList struct {
		BaseList
		list []Credit
	}
)

func NewCredit() *Credit {
	return &Credit{}
}

func (c *Credit) TenantID() uint {
	return c.tenantId
}

func (c *Credit) SetTenantID(tenantId uint) {
	c.tenantId = tenantId
}

func (c *Credit) Balance() float64 {
	return c.balance
}

func (c *Credit) SetBalance(balance float64) {
	c.balance = balance
}

//

func (c *Credit) TxAmount() Transaction {
	return c.txAmount
}

func (c *Credit) SetTxAmount(src Transaction) {
	c.txAmount = src
}

//

func (c *Credit) Transactions() TransactionList {
	return c.transactions
}

func (c *Credit) SetTransactions(src []Transaction) {
	if len(src) > 0 {
		c.transactions.SetList(src)
		c.transactions.SetTotal(int64(len(src)))
	}
}

func (c *Credit) SetTransactionList(src TransactionList) {
	if src.Total() > 0 {
		c.transactions.SetList(src.List())
		c.transactions.SetTotal(src.Total())
	}
}

//

func (c *Credit) FromDB(src model.Credits) Credit {
	// base
	c.SetID(src.ID)
	c.SetUUID(src.Uuid)
	c.SetCreatedAt(src.CreatedAt)
	c.SetUpdatedAt(src.UpdatedAt)
	c.SetDeletedAt(src.DeletedAt.Time)
	//fields
	c.SetTenantID(src.TenantID)
	c.SetBalance(src.Balance)

	if src.Transactions != nil {
		txs := NewTransactionList()
		txs.ListFromDB(src.Transactions)
		c.SetTransactions(txs.List())
	}

	return *c
}

func (c *Credit) ToDB() model.Credits {
	return model.Credits{
		BaseSql: model.BaseSql{
			Model: gorm.Model{
				ID:        c.ID(),
				CreatedAt: c.CreatedAt(),
				UpdatedAt: c.UpdatedAt(),
				DeletedAt: gorm.DeletedAt(
					sql.NullTime{
						Time: c.DeletedAt(),
						Valid: func() bool {
							if c.DeletedAt().IsZero() {
								return false
							}
							return true
						}(),
					},
				),
			},
			Uuid: c.UUID(),
		},
		TenantID: c.TenantID(),
		Balance:  c.Balance(),
	}
}

//

func NewCreditList() *CreditList { return &CreditList{} }

func (ul *CreditList) List() []Credit { return ul.list }

func (ul *CreditList) SetList(list []Credit) { ul.list = list }

func (ul *CreditList) ListToDB() []model.Credits {
	tenant := make([]model.Credits, 0)

	if len(ul.list) == 0 {
		return tenant
	}

	for _, item := range ul.list {
		tenant = append(tenant, item.ToDB())
	}

	return tenant
}

func (ul *CreditList) ListFromDB(src []model.Credits) CreditList {
	ul.list = make([]Credit, 0)

	total := len(src)
	if ul.total == 0 && total > 0 {
		ul.total = int64(total)
	}

	if ul.total == 0 {
		return *ul
	}

	for _, item := range src {
		ul.list = append(ul.list, NewCredit().FromDB(item))
	}

	return *ul
}

//

type CreditListReqQryParam struct {
	ReqBaseQryParam
}

func NewCreditListReqQryParam() *CreditListReqQryParam {
	return &CreditListReqQryParam{}
}
