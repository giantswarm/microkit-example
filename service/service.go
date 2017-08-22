// Package service implements business logic of the micro service.
package service

import (
	kitlog "github.com/go-kit/kit/log"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	Logger kitlog.Logger

	// Settings.
	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	var err error

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		Version: versionService,
	}

	return newService, nil
}

// Service bundles other services.
type Service struct {
	Version *version.Service
}
