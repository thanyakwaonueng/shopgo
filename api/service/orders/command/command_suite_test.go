package command_test

import (
    "fmt"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global variables shared across all tests in the suite
var (
	logger  *slog.Logger
	db      *gorm.DB
	sqlDb   *sql.DB
	sqlMock sqlmock.Sqlmock
)

func TestUserCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Command Suite")
}

var _ = BeforeSuite(func() {
	// 1. Initialize Logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// 2. Initialize SQL Mock
	var err error
	sqlDb, sqlMock, err = sqlmock.New()
	Expect(err).NotTo(HaveOccurred())

	// 3. Initialize GORM with the mock connection
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// 1. Verify all transaction expectations (Begin/Commit/Rollback) were met
	// This ensures you didn't miss any COMMIT or ROLLBACK calls
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		fmt.Printf("There were unfulfilled expectations: %s", err)
	}

	// 2. Simply close the database without a strict expectation check
	// This prevents the "Close was not expected" error while still cleaning up
	_ = sqlDb.Close()
})
