package dto

import (
	"github.com/google/uuid"
	"microservice/internal/domain"
)

type TenantUuid struct {
	Uuid string `header:"X.TENANT.UUID" json:"uuid"`
}

func (dto *TenantUuid) ToDomain() domain.Tenant {
	uid := uuid.MustParse(dto.Uuid)
	d := domain.NewTenant()
	d.SetUUID(uid)
	return *d
}

type ListQryRequest struct {
	Page   int    `query:"page" json:"page" validate:"omitempty,numeric"`            // integer value
	Limit  int    `query:"limit" json:"limit" validate:"omitempty,numeric"`          // integer value
	Sort   string `query:"sort" json:"sort" validate:"omitempty,pagination_sort"`    // the Sql DB columns are accepted
	Order  string `query:"order" json:"order" validate:"omitempty,pagination_order"` // "asc" or "desc"
	Search string `query:"search" json:"search" validate:"omitempty,fa_alphanum"`    // will search the value along with "username", "firstname", "lastname"
}

func (dto *ListQryRequest) EvalBaseQry() domain.ReqBaseQryParam {
	d := new(domain.ReqBaseQryParam)

	if dto.Page != 0 {
		d.SetPage(dto.Page)
	}

	if dto.Limit != 0 {
		d.SetLimit(dto.Limit)
	}

	if len(dto.Sort) > 0 {
		d.SetSort(dto.Sort)
	}

	if len(dto.Order) > 0 {
		d.SetOrder(dto.Order)
	}

	if len(dto.Search) > 0 {
		d.SetSearch(dto.Search)
	}

	return *d
}

type ListBaseResponse struct {
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"10"`
	Pages int   `json:"pages" example:"5"`
	Total int64 `json:"total" example:"45"`
}
