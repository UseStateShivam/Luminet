// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := ParseCLI()

	app := fiber.New()
	SetupProxy(app, config.TcpPort)

	log.Fatal(app.Listen(":" + config.HttpPort))
}
