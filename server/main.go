// server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// startReverseProxy starts a reverse proxy server on the specified listen address
// and forwards incoming requests to the target URL.
func startReverseProxy(listenAddr, targetURL string) error {
	// Parse the target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("error parsing target URL: %v", err)
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Optional: Customize the director (e.g., to modify headers)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	// Start the proxy server
	log.Printf("Reverse proxy running on http://%s → %s", listenAddr, targetURL)
	return http.ListenAndServe(listenAddr, proxy)
}

func main() {
	// Define the target port to forward requests to
	target := "http://localhost:3000"

	// Start the reverse proxy server
	err := startReverseProxy(":8080", target)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

