package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// Request represents a simplified HTTP request structure.
type Request struct {
	Path   string `json:"path"`   // The path of the request.
	Method string `json:"method"` // The HTTP method (GET, POST, etc.).
	Body   string `json:"body"`   // The body of the request.
}

func main() {
	// Define the tunnel ID and the local port to forward requests to.
	tunnelID := "byteworks" // TODO: Replace with the actual tunnel ID
	localPort := "5173" // TODO: Replace with the actual local port
	
	// Construct the WebSocket URL for the remote server.
	remoteServer := "ws://localhost:8080/ws/" + tunnelID // TODO: Replace with the actual remote server URL

	// Establish a WebSocket connection to the remote server.
	conn, _, err := websocket.DefaultDialer.Dial(remoteServer, nil)
	if err != nil {
		log.Fatal("❌ Failed to connect tunnel:", err)
	}
	defer conn.Close() // Ensure the connection is closed when the main function exits.

	log.Println("🔗 Tunnel established! Listening for requests...")

	for {
		var req Request
		// Read a JSON-encoded request from the WebSocket connection.
		if err := conn.ReadJSON(&req); err != nil {
			log.Println("🔌 Tunnel closed:", err)
			os.Exit(1) // Exit the program if the tunnel is closed.
		}

		// Construct the URL to forward the request to the local server.
		url := "http://localhost:" + localPort + req.Path // TODO: Replace with the actual local server URL
		// Create a new HTTP request with the method, URL, and body from the received request.
		httpReq, _ := http.NewRequest(req.Method, url, bytes.NewBuffer([]byte(req.Body)))
		client := &http.Client{} // Create an HTTP client to send the request.

		// Send the request to the local server and receive a response.
		resp, err := client.Do(httpReq)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("❌ Forwarding failed"))
			continue // Skip further processing if forwarding fails.
		}

		// Read the response body from the local server.
		body, _ := io.ReadAll(resp.Body)

		// Send the response body back through the WebSocket to the original requester.
		conn.WriteMessage(websocket.TextMessage, body)
	}
}

