package server

import (
	"encoding/json"
	"net/http"
	"sync"

	micrologger "github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"

	"github.com/giantswarm/microkit-example/server/endpoint"
	"github.com/giantswarm/microkit-example/server/middleware"
	"github.com/giantswarm/microkit-example/service"
)

// Config represents the configuration used to create a new server object.
type Config struct {
	// Dependencies.
	MicroServerConfig microserver.Config
	Service           *service.Service
}

// DefaultConfig provides a default configuration to create a new server object
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		MicroServerConfig: microserver.DefaultConfig(),
		Service:           nil,
	}
}

// New creates a new configured server object.
func New(config Config) (microserver.Server, error) {
	var err error

	var middlewareCollection *middleware.Middleware
	{
		middlewareConfig := middleware.DefaultConfig()
		middlewareConfig.Logger = config.MicroServerConfig.Logger
		middlewareConfig.Service = config.Service
		middlewareCollection, err = middleware.New(middlewareConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var endpointCollection *endpoint.Endpoint
	{
		endpointConfig := endpoint.DefaultConfig()
		endpointConfig.Logger = config.MicroServerConfig.Logger
		endpointConfig.Middleware = middlewareCollection
		endpointConfig.Service = config.Service
		endpointCollection, err = endpoint.New(endpointConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	// Create our custom server.
	newServer := &server{
		// Dependencies.
		logger: config.MicroServerConfig.Logger,

		// Internals.
		bootOnce:     sync.Once{},
		config:       config.MicroServerConfig,
		shutdownOnce: sync.Once{},
	}

	// Apply internals to the micro server config.
	newServer.config.Endpoints = []microserver.Endpoint{
		endpointCollection.Version,
	}
	newServer.config.ErrorEncoder = newServer.newErrorEncoder()
	newServer.config.RequestFuncs = newServer.newRequestFuncs()

	return newServer, nil
}

// server manages the transport logic and endpoint registration.
type server struct {
	// Dependencies.
	endpoints []microserver.Endpoint
	logger    micrologger.Logger

	// Internals.
	bootOnce     sync.Once
	config       microserver.Config
	shutdownOnce sync.Once
}

func (s *server) Boot() {
	s.bootOnce.Do(func() {
		// Here goes your custom boot logic for your server/endpoint/middleware, if
		// any.
	})
}

func (s *server) Config() microserver.Config {
	return s.config
}

func (s *server) Shutdown() {
	s.shutdownOnce.Do(func() {
		// Here goes your custom shutdown logic for your server/endpoint/middleware,
		// if any.
	})
}

// ErrorEncoder is a global error handler used for all endpoints. Errors
// received here are encoded by go-kit and express in which area the error was
// emitted. The underlying error defines the HTTP status code and the encoded
// error message. The response is always a JSON object containing an error field
// describing the error.
func (s *server) newErrorEncoder() kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		switch e := err.(type) {
		case kithttp.Error:
			err = e.Err

			switch e.Domain {
			case kithttp.DomainEncode:
				w.WriteHeader(http.StatusBadRequest)
			case kithttp.DomainDecode:
				w.WriteHeader(http.StatusBadRequest)
			case kithttp.DomainDo:
				// Your custom service errors go here.
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (s *server) newRequestFuncs() []kithttp.RequestFunc {
	return []kithttp.RequestFunc{
		func(ctx context.Context, r *http.Request) context.Context {
			// Your custom logic to enrich the request context with request specific
			// information goes here.
			return ctx
		},
	}
}
