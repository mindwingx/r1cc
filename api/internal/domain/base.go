package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Base struct {
	BaseRelation
	id        uint
	uid       uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func (bd *Base) ID() uint {
	return bd.id
}

func (bd *Base) SetID(id uint) {
	bd.id = id
}

func (bd *Base) UUID() uuid.UUID {
	return bd.uid
}

func (bd *Base) SetUUID(id uuid.UUID) {
	bd.uid = id
}

func (bd *Base) CreatedAt() time.Time {
	return bd.createdAt
}

func (bd *Base) SetCreatedAt(t time.Time) {
	bd.createdAt = t
}

func (bd *Base) UpdatedAt() time.Time {
	return bd.updatedAt
}

func (bd *Base) SetUpdatedAt(t time.Time) {
	bd.updatedAt = t
}

func (bd *Base) DeletedAt() time.Time {
	return bd.deletedAt
}

func (bd *Base) SetDeletedAt(t time.Time) {
	bd.deletedAt = t
}

//

type BaseList struct {
	total int64
}

func (bl *BaseList) Total() int64 { return bl.total }

func (bl *BaseList) SetTotal(total int64) { bl.total = total }

// collection default query params

type ReqBaseQryParam struct {
	BaseRelation
	page   int
	limit  int
	order  string
	sort   string
	search string
	relId  uint
	items  []uuid.UUID
}

func (bc *ReqBaseQryParam) Page() int {
	if bc.page != 0 {
		return bc.page
	}

	return 1
}

func (bc *ReqBaseQryParam) SetPage(page int) {
	bc.page = page
}

func (bc *ReqBaseQryParam) Limit() int {
	if bc.limit != 0 {
		if bc.limit > 50 {
			bc.limit = 50
		}

		return bc.limit
	}

	return 10 // default items count per page
}

func (bc *ReqBaseQryParam) SetLimit(limit int) {
	bc.limit = limit
}

func (bc *ReqBaseQryParam) Order() string {
	if bc.order != "" {
		return bc.order
	}

	return "desc"
}

// SetOrder default order: desc
func (bc *ReqBaseQryParam) SetOrder(order string) {
	bc.order = order
}

// Sort default "created_at"
func (bc *ReqBaseQryParam) Sort() string {
	if bc.sort != "" {
		return bc.sort
	}

	return "created_at"
}

// SetSort related field to be sorted
func (bc *ReqBaseQryParam) SetSort(sort string) { bc.sort = sort }

func (bc *ReqBaseQryParam) Search() string {
	return bc.search
}

func (bc *ReqBaseQryParam) SetSearch(search string) { bc.search = search }

func (bc *ReqBaseQryParam) RelId() uint {
	return bc.relId
}

func (bc *ReqBaseQryParam) SetRelId(relId uint) {
	bc.relId = relId
}

func (bc *ReqBaseQryParam) Items() []uuid.UUID {
	return bc.items
}

func (bc *ReqBaseQryParam) SetItems(items ...uuid.UUID) {
	bc.items = items
}

//

func (bc *ReqBaseQryParam) Offset() int {
	return (bc.Page() - 1) * bc.Limit()
}

func (bc *ReqBaseQryParam) SortOrder() string {
	return fmt.Sprintf("%s %s", bc.Sort(), bc.Order())
}

// General Purpose

type BaseRelation struct {
	dbRelations []string
}

// GetRelations (general purpose) use in repo layer to check retrieve relations or not
func (br *BaseRelation) GetRelations() []string {
	return br.dbRelations
}

// SetRelations (general purpose) by default sets to TRUE, if it needs to retrieve relations in repo layer
func (br *BaseRelation) SetRelations(relations ...string) {
	if len(relations) == 0 {
		all := "*"
		br.dbRelations = append(br.dbRelations, all)
		return
	}

	for _, r := range relations {
		br.dbRelations = append(br.dbRelations, r)
	}
}
