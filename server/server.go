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
	Logger  micrologger.Logger
	Service *service.Service
}

// DefaultConfig provides a default configuration to create a new server object
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:  nil,
		Service: nil,
	}
}

// New creates a new configured server object.
func New(config Config) (microserver.Server, error) {
	var err error

	var middlewareCollection *middleware.Middleware
	{
		middlewareConfig := middleware.DefaultConfig()
		middlewareConfig.Logger = config.Logger
		middlewareConfig.Service = config.Service
		middlewareCollection, err = middleware.New(middlewareConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var endpointCollection *endpoint.Endpoint
	{
		endpointConfig := endpoint.DefaultConfig()
		endpointConfig.Logger = config.Logger
		endpointConfig.Middleware = middlewareCollection
		endpointConfig.Service = config.Service
		endpointCollection, err = endpoint.New(endpointConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	newServer := &server{
		// Dependencies.
		endpoints: []microserver.Endpoint{
			endpointCollection.Version,
		},
		logger: config.Logger,

		// Internals.
		bootOnce:     sync.Once{},
		shutdownOnce: sync.Once{},
	}

	return newServer, nil
}

// server manages the transport logic and endpoint registration.
type server struct {
	// Dependencies.
	endpoints []microserver.Endpoint
	logger    micrologger.Logger

	// Internals.
	bootOnce     sync.Once
	shutdownOnce sync.Once
}

func (s *server) Boot() {
	s.bootOnce.Do(func() {
		// Here goes your custom boot logic for your server/endpoint/middleware, if
		// any.
	})
}

func (s *server) Endpoints() []microserver.Endpoint {
	return s.endpoints
}

// ErrorEncoder is a global error handler used for all endpoints. Errors
// received here are encoded by go-kit and express in which area the error was
// emitted. The underlying error defines the HTTP status code and the encoded
// error message. The response is always a JSON object containing an error field
// describing the error.
func (s *server) ErrorEncoder() kithttp.ErrorEncoder {
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

func (s *server) RequestFuncs() []kithttp.RequestFunc {
	return []kithttp.RequestFunc{
		func(ctx context.Context, r *http.Request) context.Context {
			// Your custom logic to enrich the request context with request specific
			// information goes here.
			return ctx
		},
	}
}

func (s *server) Shutdown() {
	s.shutdownOnce.Do(func() {
		// Here goes your custom shutdown logic for your server/endpoint/middleware,
		// if any.
	})
}
