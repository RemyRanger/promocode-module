package control_test

import (
	"APIs/internal/app"
	"APIs/internal/common/config"
	test_utils "APIs/test/integration/test-utils"
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

const app_name = "API_test"

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Suite")
}

var (
	db         *gorm.DB
	router     *chi.Mux
	dbHost     string
	dbPort     string
	pgShutdown func()
)

var _ = BeforeSuite(func() {
	var err error
	dbHost, dbPort, pgShutdown, err = test_utils.SetupPgContainer()
	Expect(err).ToNot(HaveOccurred())

	// Initialize handler to test
	app_config := config.Config{
		Server: config.Server{
			Env: "DEV",
		},
		Db: config.Db{
			Addr: fmt.Sprintf("host=%s user=test password=test dbname=testdb port=%s sslmode=disable TimeZone=UTC", dbHost, dbPort),
		},
		Logs: config.Logs{
			Level: "DEBUG",
		},
	}
	db, router = app.NewApp(app_name, app_config)
})

var _ = AfterSuite(func() {
	// Clean up the test environment
	pgShutdown()
})
