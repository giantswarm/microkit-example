package middleware

import (
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/microkit-example/server/middleware/version"
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
	var err error

	var versionMiddleware *version.Middleware
	{
		versionConfig := version.DefaultConfig()
		versionConfig.Logger = config.Logger
		versionConfig.Service = config.Service
		versionMiddleware, err = version.New(versionConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	newMiddleware := &Middleware{
		Version: versionMiddleware,
	}

	return newMiddleware, nil
}

// Middleware is middleware collection.
type Middleware struct {
	Version *version.Middleware
}
