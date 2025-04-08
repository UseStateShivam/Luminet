// proxy.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupProxy(app *fiber.App, targetPort string) {
	target := fmt.Sprintf("http://localhost:%s", targetPort)

	app.Use("/*", func(c *fiber.Ctx) error {
		// Convert body to io.Reader
		bodyReader := bytes.NewReader(c.Body())

		// Build the request to forward
		req, err := http.NewRequest(c.Method(), target+c.OriginalURL(), bodyReader)
		if err != nil {
			return err
		}

		// Copy headers from original request
		c.Request().Header.VisitAll(func(k, v []byte) {
			req.Header.Set(string(k), string(v))
		})

		// Send the request to the backend server
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Copy response headers
		for k, v := range resp.Header {
			for _, val := range v {
				c.Set(k, val)
			}
		}

		// Set status code and body
		c.Status(resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return c.Send(respBody)
	})
}
