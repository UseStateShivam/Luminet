package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Request struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

func main() {
	tunnelID := "byteworks"
	localPort := "5173"
	remoteServer := "ws://localhost:8080/ws/" + tunnelID

	conn, _, err := websocket.DefaultDialer.Dial(remoteServer, nil)
	if err != nil {
		log.Fatal("❌ Failed to connect tunnel:", err)
	}
	defer conn.Close()

	log.Println("🔗 Tunnel established! Listening for requests...")

	for {
		var req Request
		if err := conn.ReadJSON(&req); err != nil {
			log.Println("🔌 Tunnel closed:", err)
			os.Exit(1)
		}

		// Forward request to localhost
		url := "http://localhost:" + localPort + req.Path
		httpReq, _ := http.NewRequest(req.Method, url, bytes.NewBuffer([]byte(req.Body)))
		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("❌ Forwarding failed"))
			continue
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		conn.WriteMessage(websocket.TextMessage, body)
	}
}
