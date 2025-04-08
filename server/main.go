package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New() // Initialize a new Fiber app

	// Middleware to allow WebSocket upgrade requests
	app.Use("/ws/:id", func(c *fiber.Ctx) error {
		// Check if the request is a WebSocket upgrade
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next() // Allow the request to proceed
		}
		return fiber.ErrUpgradeRequired // Deny non-WebSocket requests
	})

	// Handle WebSocket connections for tunnel clients
	app.Get("/ws/:id", websocket.New(HandleClientTunnel))

	// Route to forward HTTP requests to corresponding client tunnels
	app.All("/:id/*", HandleProxyRequest)

	// Start the Fiber app and listen on the default port
	log.Fatal(app.Listen(":8080"))
}

