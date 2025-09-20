package orm

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	// Query this custom type is defined to be adaptable
	//to handle any other model driver
	Query *gorm.DB

	SeederItem struct {
		Dependency interface{}   // the modules domain instance, ex. : domain.User{}
		Data       []interface{} // the domain mocked data, base on the related struct
	}

	ISql interface {
		ISqlGeneric
		ISqlTx
	}

	ISqlGeneric interface {
		Init()
		Migrate(path string)
		Seed(items []SeederItem)
		DB() *gorm.DB
		Close()
		Fx(lc fx.Lifecycle) ISqlTx
	}

	ISqlTx interface {
		// Begin transaction
		Begin()
		// Commit commits the transaction.
		Commit() error
		// Rollback rolls back the transaction.
		Rollback() error
		// Resolve commit or rollback transaction by getting the error
		Resolve(err error) error
		// Tx returns the current transaction or the base db if no transaction is active.
		Tx() gorm.DB
	}
)
