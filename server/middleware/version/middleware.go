package version

import (
	micrologger "github.com/giantswarm/microkit/logger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/giantswarm/microkit-example/service"
)

// Config represents the configuration used to create a version middleware.
type Config struct {
	// Dependencies.
	Logger  micrologger.Logger
	Service *service.Service
}

// DefaultConfig provides a default configuration to create a new version
// middleware by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:  nil,
		Service: nil,
	}
}

// New creates a new configured version middleware.
func New(config Config) (*Middleware, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, maskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.Service == nil {
		return nil, maskAnyf(invalidConfigError, "service must not be empty")
	}

	newMiddleware := &Middleware{
		Config: config,
	}

	return newMiddleware, nil
}

type Middleware struct {
	Config
}

func (m *Middleware) Middleware(next kitendpoint.Endpoint) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		// Your middleware logic goes here.

		return next(ctx, request)
	}
}
