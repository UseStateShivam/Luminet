package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Request struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

var clients = make(map[string]*websocket.Conn)

func HandleClientTunnel(c *websocket.Conn) {
	id := c.Params("id")
	clients[id] = c
	log.Printf("🔗 Client [%s] connected.\n", id)

	defer func() {
		log.Printf("❌ Client [%s] disconnected.\n", id)
		delete(clients, id)
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func HandleProxyRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	conn, ok := clients[id]
	if !ok {
		return c.Status(502).SendString("❌ No tunnel client connected")
	}

	req := Request{
		Path:   c.OriginalURL(),
		Method: c.Method(),
		Body:   string(c.Body()),
	}
	data, _ := json.Marshal(req)

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}

	_, resp, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	return c.Send(resp)
}
