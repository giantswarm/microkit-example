// Package service implements business logic of the micro service.
package service

import (
	kitlog "github.com/go-kit/kit/log"

	"github.com/giantswarm/microkit-example/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	Logger kitlog.Logger

	// Settings.
	Description    string
	GitCommit      string
	Name           string
	ProjectVersion string
	Source         string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Description:    "",
		GitCommit:      "",
		Name:           "",
		ProjectVersion: "",
		Source:         "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	var err error

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Logger = config.Logger

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.ProjectVersion = config.ProjectVersion
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, maskAny(err)
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