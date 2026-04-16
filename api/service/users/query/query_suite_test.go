package query_test

import (
	"log/slog"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global variables shared across all query tests (GetUser, ListUsers, etc.)
var (
	logger *slog.Logger
)

func TestUserQuery(t *testing.T) {
	// Links the 'go test' command to the Ginkgo runner
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Query Suite")
}

var _ = BeforeSuite(func() {
	// Initialize a standard logger for your handlers
	// LevelDebug allows you to see detailed output if a test fails
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
})

var _ = AfterSuite(func() {
	// No database connections to close since we are mocking the Repository layer
})
