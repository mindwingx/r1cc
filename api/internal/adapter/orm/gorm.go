package orm

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"microservice/config"
	"microservice/internal/adapter/registry"
	"microservice/pkg/utils"
	"os"
	"sort"
	"time"
)

type sql struct {
	config config.Database
	db     gorm.DB
}

func New(service *config.Service, registry registry.IRegistry) ISqlGeneric {
	db := new(sql)

	if err := registry.Parse(&db.config); err != nil {
		utils.PrintStd(utils.StdPanic, "database", "config parse err: %s", err)
	}

	db.config.Debug = service.Debug
	return db
}

func (q *sql) Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		q.config.Host,
		q.config.Username,
		q.config.Password,
		q.config.Database,
		q.config.Port,
		q.config.Ssl,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 q.newGormLog(q.config.SlowSqlThreshold),
		NowFunc:                func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		utils.PrintStd(utils.StdPanic, "database", "connection err: %s", err)
	}

	sqlDatabase, err := database.DB()
	if err != nil {
		utils.PrintStd(utils.StdPanic, "database", "init err: %s", err)
	}

	if q.config.MaxIdleConnections != 0 {
		sqlDatabase.SetMaxIdleConns(q.config.MaxIdleConnections)
	}

	if q.config.MaxOpenConnections != 0 {
		sqlDatabase.SetMaxOpenConns(q.config.MaxOpenConnections)
	}

	if q.config.MaxLifetimeSeconds != 0 {
		sqlDatabase.SetConnMaxLifetime(time.Second * time.Duration(q.config.MaxLifetimeSeconds))
	}

	if q.config.Debug {
		database = database.Debug()
		utils.PrintStd(utils.StdLog, "database", "debug is enabled")
	}

	q.db = *database
}

func (q *sql) Migrate(path string) {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		utils.PrintStd(utils.StdPanic, "database", "migrations dir scan err: %s", err)
	}

	defer func() { _ = dir.Close() }()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		utils.PrintStd(utils.StdPanic, "database", "reading directory contents err: %s", err)
	}

	// Sort the entries alphabetically by name - model file order by numeric(01, 02, etc)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	// Iterate over the file info slice and print the file names
	for _, fileInfo := range fileInfos {
		if fileInfo.Mode().IsRegular() {
			if err = q.db.Exec(q.parseSqlFile(path, fileInfo)).Error; err != nil {
				utils.PrintStd(utils.StdPanic, "database", "migrate err: %s", err)
			}
		}
	}
}

func (q *sql) Seed(items []SeederItem) {
	if len(items) > 0 {
		var count int64

		for _, item := range items {
			instance := q.db.Model(&item.Dependency)
			result := instance.Count(&count)

			if result.Error != nil {
				utils.PrintStd(utils.StdPanic, "database", "seed prepare err: %s", result.Error)
			}

			if (count == 0) && (len(item.Data) > 0) {
				utils.PrintStd(utils.StdLog, "database", "seeding started...")

				for _, data := range item.Data {
					create := instance.Create(data)
					if create.Error != nil {
						utils.PrintStd(utils.StdPanic, "database", "seeding err: %s", create.Error)
					}
				}

				utils.PrintStd(utils.StdLog, "database", "seeding finished")
			}
		}
	}
}

func (q *sql) DB() *gorm.DB {
	return &q.db
}

func (q *sql) Fx(lc fx.Lifecycle) ISqlTx {
	uow := NewTransaction(&q.db)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "database", "initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "database", "stopping...")
			defer utils.PrintStd(utils.StdLog, "database", "stopped")

			q.Close()
			return
		},
	})

	return uow
}

func (q *sql) Close() {
	sqlDatabase, err := q.db.DB()
	if err != nil {
		utils.PrintStd(utils.StdLog, "database", "connection close retrieve", err)
	}

	err = sqlDatabase.Close()
	if err != nil {
		utils.PrintStd(utils.StdLog, "database", "connection close failure", err)
	}
}

// HELPER METHODS

func (q *sql) newGormLog(SlowSqlThreshold int) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Duration(SlowSqlThreshold) * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn,                                   // Log level
			IgnoreRecordNotFoundError: false,                                         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                          // Disable color
		})
}

func (q *sql) parseSqlFile(path string, fileInfo os.FileInfo) string {
	sqlFile := fmt.Sprintf("%s/%s", path, fileInfo.Name())
	sqlBytes, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		utils.PrintStd(utils.StdLog, "database", "SQL file parse", err)
	}
	// Convert SQL file contents to string
	return string(sqlBytes)
}
