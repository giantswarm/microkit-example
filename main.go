package main

import (
	"fmt"

	"github.com/giantswarm/microkit/command"
	"github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"

	"github.com/giantswarm/microkit-example/server"
	"github.com/giantswarm/microkit-example/service"
)

var (
	description    string = "This is an example microservice using the microkit framework."
	gitCommit      string = "n/a"
	name           string = "microkit-example"
	projectVersion string = "n/a"
	source         string = "https://github.com/giantswarm/microkit-example"
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger logger.Logger
	{
		newLogger, err = logger.New(logger.DefaultConfig())
		if err != nil {
			panic(err)
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func() microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.Logger = newLogger

			serviceConfig.Description = description
			serviceConfig.GitCommit = gitCommit
			serviceConfig.Name = name
			serviceConfig.ProjectVersion = projectVersion
			serviceConfig.Source = source

			newService, err = service.New(serviceConfig)
			if err != nil {
				panic(err)
			}
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			serverConfig := server.DefaultConfig()

			serverConfig.Logger = newLogger
			serverConfig.Service = newService

			newServer, err = server.New(serverConfig)
			if err != nil {
				panic(err)
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger
		commandConfig.ServerFactory = newServerFactory

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.ProjectVersion = projectVersion
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
		if err != nil {
			fmt.Printf("%#v\n", err)
			panic(err)
		}
	}

	newCommand.New().Execute()
}
