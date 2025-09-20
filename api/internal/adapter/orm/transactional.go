package orm

import (
	"gorm.io/gorm"
)

type transactional struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewTransaction(db *gorm.DB) ISqlTx {
	return &transactional{
		db: db,
	}
}

// Begin starts a new transaction.
func (u *transactional) Begin() {
	u.tx = u.db.Begin()
}

// Commit commits the transaction.
func (u *transactional) Commit() (err error) {
	if u.tx != nil {
		err = u.tx.Commit().Error
	}
	return
}

// Rollback rolls back the transaction.
func (u *transactional) Rollback() (err error) {
	if u.tx != nil {
		err = u.tx.Rollback().Error
	}
	return
}

// Resolve commit or rollback transaction by getting the error
func (u *transactional) Resolve(dbErr error) (err error) {
	if dbErr != nil {
		err = u.Rollback()
		u.tx = nil
		return
	}

	dbErr = u.Commit()
	u.tx = nil
	return
}

// Tx returns the current transaction or the base db if no transaction is active.
func (u *transactional) Tx() gorm.DB {
	if u.tx != nil {
		return *u.tx
	}
	return *u.db
}
