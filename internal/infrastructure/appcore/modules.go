package appcore

import (
	"users/internal/infrastructure/postgres"
)

var EchoModule = Options(
	Provide(ProvideEchoServer),
)

var EchoServer = Options(
	Invoke(StartEchoServer),
)

var PostgresModule = Options(
	Provide(postgres.ProvidePostgresConfig),
	Provide(postgres.New),
)

// HMACModule optional module for HMAC authentication
var HMACModule = Options(
	Provide(ProvideHMACConfig),
	Provide(ProvideHMACMiddleware),
	Invoke(SetupHMACMiddleware),
)

// LoggingModule provides logger instance
var LoggingModule = Options(
	Provide(ProvideLogger),
)

// LifecycleModule provides lifecycle hooks
var LifecycleModule = Options(
	Provide(NewLifecycleHooks),
	Invoke(SetupGracefulShutdown),
)
