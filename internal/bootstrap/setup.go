package bootstrap

import (
	"users/internal/router"

	appcore "github.com/mortired/appsap-core"

	logging "github.com/mortired/appsap-logging"

	"github.com/labstack/echo/v4"
)

func Setup() *appcore.Application {
	options := []appcore.Option{
		// Infrastructure
		appcore.PostgresModule,
		appcore.HMACModule,

		// Repositories
		appcore.Provide(ProvideUserRepository),

		// Services
		appcore.Provide(ProvideUserService),

		// Controllers
		appcore.Provide(ProvideUserController),

		// Echo module
		appcore.EchoModule,

		// Setup Echo middleware with logger
		appcore.Invoke(SetupEchoMiddleware),

		// Router
		appcore.Invoke(router.SetupRoutes),

		// HTTP Server
		appcore.EchoServer,
	}

	return appcore.New(options...)
}

// SetupEchoMiddleware sets up Echo middleware with logger
func SetupEchoMiddleware(e *echo.Echo, logger *logging.Logger) {
	appcore.SetupEchoMiddleware(e, logger)
}
