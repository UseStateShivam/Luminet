package main

import (
	"net"

	"github.com/gofiber/fiber/v2"
)

// A reverse proxy is like a middleman between a client and your backend server. But unlike a forward proxy (which sits in front of the client), the reverse proxy sits in front of your servers, and clients think they're talking to the main server — when really, the proxy is just passing stuff back and forth
// Think of it as the bouncer at a club. People (clients) come to the front, talk to the bouncer (reverse proxy), and the bouncer decides:
// Who gets in
// Where they go (User service? Payment? Auth?)
// Maybe even adds some checks (ID = auth, rate-limits, logs)
// main is the entry point of the program. It sets up a reverse proxy using the Fiber framework.
func main() {
	// app is the main app instance of the Fiber framework.
	app := fiber.New()

	// app.Use is a middleware function that is called before each request. It adds a route to the app that listens for requests on "/ping".
	// The second argument is a handler function that is called when a request is received on the "/ping" route.
	app.Use("/ping", func(c *fiber.Ctx) error {
		// conn is a TCP connection to the server at localhost:8000. The connection is established using net.Dial.
		conn, dialError := net.Dial("tcp", "localhost:8000")
		if dialError != nil {
			// If there is an error establishing the connection, it is returned as the error of the handler function.
			return dialError
		}
		// conn.Write is called to send the string "PING" to the server.
		conn.Write([]byte("PING"))
		// buffer is a slice of bytes that is used to store the response from the server.
		buffer := make([]byte, 1024)
		// n is the number of bytes read from the server and readError is an error that is returned if there is an error reading from the server.
		n, readError := conn.Read(buffer)
		if readError != nil {
			// If there is an error reading from the server, it is returned as the error of the handler function.
			return readError
		}
		// message is a string that is created from the bytes read from the server.
		message := string(buffer[:n])
		// c.SendString is called to send the response back to the client as a string.
		return c.SendString(message)
	})

	// app.Listen is called to start the server. The argument is the address that the server should listen on.
	// The error returned by app.Listen is checked and if it is not nil, the program panics with the error.
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
	// The server has started successfully, so a message is printed to the console to indicate this.
	println("Server started on port 8080")
}
