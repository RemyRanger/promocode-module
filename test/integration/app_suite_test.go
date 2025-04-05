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
	"github.com/onsi/gomega/ghttp"
	"gorm.io/gorm"
)

const app_name = "app_test"

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Suite")
}

var (
	db            *gorm.DB
	router        *chi.Mux
	weatherServer *ghttp.Server
	dbHost        string
	dbPort        string
	pgShutdown    func()
)

var _ = BeforeSuite(func() {
	var err error
	dbHost, dbPort, pgShutdown, err = test_utils.SetupPgContainer()
	Expect(err).ToNot(HaveOccurred())

	weatherServer = test_utils.BootstrapUpstreamServer()
	DeferCleanup(func() {
		weatherServer.Close()
	})

	// Initialize handler to test
	app_config := config.Config{
		Server: config.Server{
			Env: "DEV",
		},
		Db: config.Db{
			Addr: fmt.Sprintf("host=%s user=test password=test dbname=testdb port=%s sslmode=disable TimeZone=UTC", dbHost, dbPort),
		},
		Openweather: config.Openweather{
			Url:    weatherServer.URL(),
			Apikey: "weather_api_key",
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
