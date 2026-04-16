package command_test

import (
	"log/slog"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global variables shared across all command tests (Register, Update, etc.)
var (
	logger *slog.Logger
)

func TestUserCommand(t *testing.T) {
	// Links 'go test' to the Ginkgo runner
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Command Suite")
}

var _ = BeforeSuite(func() {
	// Initialize a logger that outputs to the console during tests
	// This helps you see your slog.Info/Error calls when a test fails
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
})

var _ = AfterSuite(func() {
	// Logic here runs once after ALL tests in the command folder finish
})
