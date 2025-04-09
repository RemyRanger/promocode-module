package app

import (
	"APIs/internal/common/config"
	"APIs/internal/common/db"
	"APIs/internal/common/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"

	promocode_handler "APIs/internal/services/promocode/adapters/handler"
	promocode_repository "APIs/internal/services/promocode/adapters/repository"
	promocode_core "APIs/internal/services/promocode/core"
	promocode_ports "APIs/internal/services/promocode/ports"

	weather_client "APIs/internal/services/weather/adapters/openweather_client"
	weather_repository "APIs/internal/services/weather/adapters/repository"
	weather_core "APIs/internal/services/weather/core"

	"gorm.io/gorm"
)

func NewApp(app_name string, app_config config.Config) (*gorm.DB, *chi.Mux) {
	// Init logger
	logger.NewZerolog(app_config)

	// Init db
	gormDb := db.NewPostgres(app_config)

	// Init router
	router := chi.NewRouter()
	pathPrefix := "/v1"

	// Register Middlewares
	router.Use(
		logger.Logger, // Use custom logger

		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		middleware.Timeout(30*time.Second),

		// middlewares
		middleware.Recoverer,

		// response type is forced to JSON
		render.SetContentType(render.ContentTypeJSON),

		// healthcheck
		middleware.Heartbeat(pathPrefix+"/healthcheck"),

		// telemetry
		otelchi.Middleware(app_name, otelchi.WithChiRoutes(router)),
	)

	// Register global services
	weatherService := weather_core.NewService(weather_repository.New(gormDb), weather_client.NewClientAPI(app_config))
	promocode_ports.HandlerWithOptions(promocode_handler.NewHandler(promocode_core.NewService(promocode_repository.New(gormDb), weatherService)), promocode_ports.ChiServerOptions{BaseURL: pathPrefix, BaseRouter: router})

	return gormDb, router
}
