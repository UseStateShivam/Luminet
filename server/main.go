package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	// Only allow upgrade if it’s a WebSocket
	app.Use("/ws/:id", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Tunnel endpoint for clients
	app.Get("/ws/:id", websocket.New(HandleClientTunnel))

	// Forward requests from web to the client
	app.All("/:id/*", HandleProxyRequest)

	log.Fatal(app.Listen(":8080")) // Render's default port
}
