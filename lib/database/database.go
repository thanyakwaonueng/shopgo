package database

import (
	"database/sql"
	//"github.com/thanyakwaonueng/shopgo/lib/environment"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	//"time"

	"github.com/DATA-DOG/go-sqlmock"
	extraClausePlugin "github.com/WinterYukky/gorm-extra-clause-plugin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(dsn string) *gorm.DB {
	var conf *gorm.Config
	if util.IsTestMode() {
		conf = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	} else {
		conf = &gorm.Config{
			// Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), conf)
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}

	// Use plug-in
	if err := db.Use(extraClausePlugin.New()); err != nil {
		panic("Failed to initialize plugin: " + err.Error())
	}
    
    //config db e.g. set max open connection
	//postgresDb, _ := db.DB()

	return db
}

func NewMockDb() (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	sqlDb, sqlMock, err := sqlmock.New()
	if err != nil {
		panic("Failed to open sqlmock database: " + err.Error())
	}

	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDb,
		PreferSimpleProtocol: true,
	})
	mockDb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("Failed to open gorm database: " + err.Error())
	}

	return mockDb, sqlDb, sqlMock
}
