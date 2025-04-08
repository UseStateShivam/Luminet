package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Request represents a request from the proxy server to the client tunnel.
type Request struct {
	Path   string `json:"path"`   // The path of the request
	Method string `json:"method"` // The HTTP method of the request
	Body   string `json:"body"`   // The body of the request
}

// clients is a map of client IDs to their corresponding WebSocket connections.
var clients = make(map[string]*websocket.Conn)

// HandleClientTunnel handles incoming WebSocket connections from client tunnels.
func HandleClientTunnel(c *websocket.Conn) {
	id := c.Params("id")
	clients[id] = c
	log.Printf("🔗 Client [%s] connected.\n", id)

	defer func() {
		log.Printf("❌ Client [%s] disconnected.\n", id)
		delete(clients, id)
		c.Close()
	}()

	// Read messages from the client and discard them.
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

// HandleProxyRequest handles incoming HTTP requests from the proxy server and forwards them to the client tunnel.
func HandleProxyRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	conn, ok := clients[id]
	if !ok {
		return c.Status(502).SendString("❌ No tunnel client connected")
	}

	// Marshal the request into a JSON object.
	req := Request{
		Path:   c.OriginalURL(),
		Method: c.Method(),
		Body:   string(c.Body()),
	}
	data, _ := json.Marshal(req)

	// Send the request to the client tunnel.
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}

	// Read the response from the client tunnel.
	_, resp, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	// Return the response to the proxy server.
	return c.Send(resp)
}

