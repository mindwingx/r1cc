package credit

import (
	"encoding/hex"
	"math"
	"microservice/internal/domain"
	"microservice/internal/modules/dto"
	"microservice/pkg/utils"
)

type IncreaseCreditRequest struct {
	Amount float64 `json:"amount" validate:"required,numeric" example:"10.50"`
}

func (dto *IncreaseCreditRequest) ToDomain() domain.Transaction {
	d := domain.NewTransaction()
	d.SetAmount(dto.Amount)
	return *d
}

type IncreaseCreditResponse struct {
	ID      string  `json:"id" example:"d13752d98dd22ab094b947f4346f15134819c93f7b1ff658c832ae466f6ebb36"`
	Amount  float64 `json:"amount" example:"10.50"`
	Balance float64 `json:"balance" example:"82.70"`
}

func IncreaseCreditResp(src domain.Credit) IncreaseCreditResponse {
	transaction := src.TxAmount()

	amount := utils.RoundToPrecision(transaction.Amount(), 4)
	return IncreaseCreditResponse{
		ID:      hex.EncodeToString(transaction.ID()),
		Amount:  amount,
		Balance: src.Balance(),
	}
}

//

type ListQryRequest struct {
	dto.ListQryRequest
}

func (dto *ListQryRequest) ToDomain() domain.TransactionListReqQryParam {
	qry := domain.NewTransactionListReqQryParam()
	qry.ReqBaseQryParam = dto.EvalBaseQry()

	return *qry
}

type (
	ListItemDetail struct {
		ID          string  `json:"uuid" example:"67f5627c-2d71-48f0-8afc-b7bed370bb45"`
		Amount      float64 `json:"amount" example:"10.23"`
		Incremented bool    `json:"incremented" example:"true"`
		CreatedAt   string  `json:"createdAt" example:"2025-01-01 12:13:14"`
	}

	ListResponse struct {
		dto.ListBaseResponse
		Balance      float64          `json:"balance"`
		Transactions []ListItemDetail `json:"transactions"`
	}
)

func ListResp(qry domain.TransactionListReqQryParam, src domain.Credit) ListResponse {
	transactions := src.Transactions()

	list := new(ListResponse)
	list.Page = qry.Page()
	list.Limit = qry.Limit()
	list.Pages = int(math.Ceil(float64(transactions.Total()) / float64(qry.Limit())))
	list.Total = transactions.Total()
	list.Balance = src.Balance()
	list.Transactions = make([]ListItemDetail, 0)

	if len(transactions.List()) > 0 {
		for _, transaction := range transactions.List() {
			amount := utils.RoundToPrecision(transaction.Amount(), 4)

			list.Transactions = append(list.Transactions, ListItemDetail{
				ID:          hex.EncodeToString(transaction.ID()),
				Amount:      amount,
				Incremented: transaction.MessageHashID() == nil,
				CreatedAt:   transaction.CreatedAt().Format("2006-01-02 15:04:05"),
			})
		}
	}

	return *list
}
