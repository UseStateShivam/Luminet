package main

import (
	"flag"
	"fmt"
)

// Config represents the configuration for the reverse proxy.
type Config struct {
	// HttpPort is the port number that the reverse proxy will listen on.
	HttpPort string
	// TcpPort is the port number that the reverse proxy will forward TCP requests to.
	TcpPort string
}

// ParseCLI parses the command-line flags and returns a Config.
func ParseCLI() Config {
	// The port number that the reverse proxy will listen on.
	httpPort := flag.String("port", "8080", "HTTP port for reverse proxy")
	// The port number that the reverse proxy will forward TCP requests to.
	tcpPort := flag.String("tcpport", "8000", "TCP server to forward to")
	// Parse the command-line flags.
	flag.Parse()

	// Print a message indicating the port numbers that the reverse proxy
	// will listen on and forward to.
	fmt.Printf("🌐 Reverse Proxy will listen on :%s\n", *httpPort)
	fmt.Printf("🔌 Will forward TCP requests to localhost:%s\n\n", *tcpPort)

	// Return the configuration.
	return Config{
		HttpPort: *httpPort,
		TcpPort:  *tcpPort,
	}
}

