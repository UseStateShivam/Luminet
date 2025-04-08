package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// SetupProxy sets up the Fiber app to act as a reverse proxy.
//
// It takes a target port number as a parameter, which is the port number that the
// reverse proxy will forward incoming HTTP requests to.
func SetupProxy(app *fiber.App, targetPort string) {
	// Set up the target URL that the reverse proxy will forward requests to.
	target := fmt.Sprintf("http://localhost:%s", targetPort) // TODO: Replace with the actual target URL

	// Set up a middleware function that will be called for every incoming request.
	app.Use("/*", func(c *fiber.Ctx) error {
		// Convert the request body to an io.Reader.
		bodyReader := bytes.NewReader(c.Body())

		// Build the request to forward to the target URL.
		req, err := http.NewRequest(c.Method(), target+c.OriginalURL(), bodyReader)
		if err != nil {
			return err
		}

		// Copy all headers from the original request to the new request.
		// This ensures that the target server receives the same headers as the
		// original request.
		c.Request().Header.VisitAll(func(k, v []byte) {
			req.Header.Set(string(k), string(v))
		})

		// Send the request to the target server.
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Copy all response headers from the target server to the Fiber response.
		// This ensures that the client receives the same headers as the target
		// server sent.
		for k, v := range resp.Header {
			for _, val := range v {
				c.Set(k, val)
			}
		}

		// Set the status code and body of the Fiber response to match the
		// response from the target server.
		c.Status(resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return c.Send(respBody)
	})
}

