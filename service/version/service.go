package version

import (
	"runtime"

	micrologger "github.com/giantswarm/microkit/logger"
)

// Config represents the configuration used to create a version service.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new version service
// by best effort.
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

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, maskAnyf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Description == "" {
		return nil, maskAnyf(invalidConfigError, "description commit must not be empty")
	}
	if config.GitCommit == "" {
		return nil, maskAnyf(invalidConfigError, "git commit must not be empty")
	}
	if config.Name == "" {
		return nil, maskAnyf(invalidConfigError, "name must not be empty")
	}
	if config.Source == "" {
		return nil, maskAnyf(invalidConfigError, "name must not be empty")
	}

	newService := &Service{
		Config: config,

		GoVersion: runtime.Version(),
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	Config

	GoVersion string
}

func (s *Service) Get(request Request) (*Response, error) {
	response := DefaultResponse()

	response.Description = s.Description
	response.GitCommit = s.GitCommit
	response.GoVersion = runtime.Version()
	response.Name = s.Name
	response.OSArch = runtime.GOOS + "/" + runtime.GOARCH
	response.Source = s.Source

	return response, nil
}
