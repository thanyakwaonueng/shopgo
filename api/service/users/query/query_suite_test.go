package query_test

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var logger *slog.Logger

func TestUserQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Query Suite")
}

var _ = BeforeSuite(func() {
	logger = slog.Default()
})
