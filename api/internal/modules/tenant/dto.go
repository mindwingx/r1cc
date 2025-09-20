package tenant

import (
	"github.com/google/uuid"
	"math"
	"microservice/internal/domain"
	"microservice/internal/modules/dto"
)

type CreateRequest struct {
	Username   string `json:"username" validate:"required,alphanum" example:"someco"`
	TenantName string `json:"tenantName" validate:"required,fa_alphanum" example:"Jack"`
}

func (dto *CreateRequest) ToDomain() domain.Tenant {
	d := domain.NewTenant()
	d.SetUsername(dto.Username)
	d.SetTenantName(dto.TenantName)

	return *d
}

type CreateResponse struct {
	Uuid   string `json:"uuid" example:"e48c48a3-cb72-4d64-b035-5c30fc900ef6"`
	Active bool   `json:"active" example:"false"`
}

func CreateResp(src domain.Tenant) CreateResponse {
	return CreateResponse{
		Uuid: func() string {
			if src.UUID() == uuid.Nil {
				return ""
			}

			return src.UUID().String()
		}(),
		Active: src.Active(),
	}
}

//

type DetailsRequest struct {
	Uuid string `json:"uuid" param:"uuid" validate:"required,uuid" example:"bf56c6b6-dd02-47ba-8dc4-bd7d2843a77a"`
}

func (dto *DetailsRequest) ToDomain() domain.Tenant {
	id := uuid.MustParse(dto.Uuid)
	d := domain.NewTenant()
	d.SetUUID(id)
	return *d
}

type (
	Credit struct {
		Balance float64 `json:"balance" example:"10.0000"`
	}
	DetailsResponse struct {
		Uuid       string `json:"uuid"  example:"bf56c6b6-dd02-47ba-8dc4-bd7d2843a77a"`
		Username   string `json:"username"  example:"dummyUsername"`
		TenantName string `json:"tenantName"  example:"Jack"`
		Active     bool   `json:"active"  example:"true"`
		Credit     Credit `json:"credit"`
	}
)

func DetailsResp(src domain.Tenant) DetailsResponse {
	detail := DetailsResponse{
		Uuid:       src.UUID().String(),
		Username:   src.Username(),
		TenantName: src.TenantName(),
		Active:     src.Active(),
	}

	if credit := src.Credit(); credit.ID() != 0 {
		detail.Credit = Credit{Balance: credit.Balance()}
	}

	return detail
}

//

type ListQryRequest struct {
	dto.ListQryRequest
}

func (dto *ListQryRequest) ToDomain() domain.TenantListReqQryParam {
	qry := domain.NewTenantListReqQryParam()
	qry.ReqBaseQryParam = dto.EvalBaseQry()

	return *qry
}

type (
	ListItemDetail struct {
		Uuid       string `json:"uuid" example:"67f5627c-2d71-48f0-8afc-b7bed370bb45"`
		TenantName string `json:"tenantName" example:"R1 Cloud"`
		Active     bool   `json:"active" example:"true"`
	}

	ListResponse struct {
		dto.ListBaseResponse
		Tenants []ListItemDetail `json:"items"`
	}
)

func ListResp(qry domain.TenantListReqQryParam, src domain.TenantList) ListResponse {
	list := new(ListResponse)
	list.Page = qry.Page()
	list.Limit = qry.Limit()
	list.Pages = int(math.Ceil(float64(src.Total()) / float64(qry.Limit())))
	list.Total = src.Total()
	list.Tenants = make([]ListItemDetail, 0)

	if len(src.List()) > 0 {
		for _, tenant := range src.List() {
			list.Tenants = append(list.Tenants, ListItemDetail{
				Uuid:       tenant.UUID().String(),
				TenantName: tenant.TenantName(),
				Active:     tenant.Active(),
			})
		}
	}

	return *list
}
