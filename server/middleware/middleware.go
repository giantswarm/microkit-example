package middleware

import (
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/microkit-example/service"
)

// Config represents the configuration used to create a middleware.
type Config struct {
	// Dependencies.
	Logger  micrologger.Logger
	Service *service.Service
}

// DefaultConfig provides a default configuration to create a new
// middleware by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:  nil,
		Service: nil,
	}
}

// New creates a new configured middleware.
func New(config Config) (*Middleware, error) {
	newMiddleware := &Middleware{}

	return newMiddleware, nil
}

// Middleware is middleware collection.
type Middleware struct {
}
