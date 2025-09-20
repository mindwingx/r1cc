package message

import (
	"math"
	"microservice/internal/domain"
	"microservice/internal/modules/dto"
)

type SendMessageRequest struct {
	Channel string `json:"channel" validate:"required,ascii,oneof=event.prod event.express" example:"event.prod"`
	Mobile  string `json:"mobile" validate:"required,mobile" example:"09123456789"`
	Message string `json:"message"  validate:"required,fa_alphanum" example:"some dummy message"`
}

func (dto *SendMessageRequest) ToDomain() domain.Message {
	d := domain.NewMessage()
	d.SetChannel(dto.Channel)
	d.SetMobile(dto.Mobile)
	d.SetMessageText(dto.Message)
	return *d
}

//

//

type ListQryRequest struct {
	dto.ListQryRequest
}

func (dto *ListQryRequest) ToDomain() domain.MessageListReqQryParam {
	qry := domain.NewMessageListReqQryParam()
	qry.ReqBaseQryParam = dto.EvalBaseQry()

	return *qry
}

type (
	ListItemDetail struct {
		Channel string `json:"channel" example:"event.prod"`
		Mobile  string `json:"mobile" example:"09123456789"`
		Message string `json:"message" example:"Hello R1 Cloud"`
		Status  string `json:"status" example:"sent"`
	}

	ListResponse struct {
		dto.ListBaseResponse
		Messages []ListItemDetail `json:"items"`
	}
)

func ListResp(qry domain.MessageListReqQryParam, src domain.MessageList) ListResponse {
	list := new(ListResponse)
	list.Page = qry.Page()
	list.Limit = qry.Limit()
	list.Pages = int(math.Ceil(float64(src.Total()) / float64(qry.Limit())))
	list.Total = src.Total()
	list.Messages = make([]ListItemDetail, 0)

	if len(src.List()) > 0 {
		for _, message := range src.List() {
			list.Messages = append(list.Messages, ListItemDetail{
				Channel: message.Channel(),
				Mobile:  message.Mobile(),
				Message: message.MessageText(),
				Status:  message.Status(),
			})
		}
	}

	return *list
}
