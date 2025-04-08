package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// main is the entry point of the reverse proxy server.
func main() {
	// Parse the command-line flags and get the configuration.
	config := ParseCLI()

	// Create a new Fiber app instance.
	app := fiber.New()

	// Set up the proxy routes for the app.
	//
	// The proxy routes will forward incoming HTTP requests to the target TCP
	// server specified in the configuration.
	SetupProxy(app, config.TcpPort)

	// Start the Fiber app and listen on the specified HTTP port.
	//
	// The HTTP port is specified in the configuration and defaults to 8080.
	log.Fatal(app.Listen(":" + config.HttpPort))
}

