package endpoint

import (
	"github.com/giantswarm/microendpoint/endpoint/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/microkit-example/server/middleware"
	"github.com/giantswarm/microkit-example/service"
)

// Config represents the configuration used to create a endpoint.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

// DefaultConfig provides a default configuration to create a new endpoint by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:     nil,
		Middleware: nil,
		Service:    nil,
	}
}

// New creates a new configured endpoint.
func New(config Config) (*Endpoint, error) {
	var err error

	var versionEndpoint *version.Endpoint
	{
		versionConfig := version.DefaultConfig()
		versionConfig.Logger = config.Logger
		versionConfig.Service = config.Service.Version
		versionEndpoint, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newEndpoint := &Endpoint{
		Version: versionEndpoint,
	}

	return newEndpoint, nil
}

// Endpoint is the endpoint collection.
type Endpoint struct {
	Version *version.Endpoint
}
