package appcore

import (
	"go.uber.org/fx"
)

// Option represents an option for application configuration
type Option = fx.Option

// Provide creates an option for providing dependency
func Provide(constructors ...interface{}) Option {
	return fx.Provide(constructors...)
}

// Invoke creates an option for invoking function
func Invoke(funcs ...interface{}) Option {
	return fx.Invoke(funcs...)
}

// Options combines multiple options
func Options(opts ...Option) Option {
	return fx.Options(opts...)
}
