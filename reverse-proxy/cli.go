// cli.go
package main

import (
	"flag"
	"fmt"
)

type Config struct {
	HttpPort string
	TcpPort  string
}

func ParseCLI() Config {
	httpPort := flag.String("port", "8080", "HTTP port for reverse proxy")
	tcpPort := flag.String("tcpport", "8000", "TCP server to forward to")
	flag.Parse()

	fmt.Printf("🌐 Reverse Proxy will listen on :%s\n", *httpPort)
	fmt.Printf("🔌 Will forward TCP requests to localhost:%s\n\n", *tcpPort)

	return Config{
		HttpPort: *httpPort,
		TcpPort:  *tcpPort,
	}
}
