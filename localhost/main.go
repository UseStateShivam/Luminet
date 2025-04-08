package main

import (
	"io"
	"net"
	"time"
)

func main() {
	// TCP
	// - Works at lower level (TCP layer)
	// - Sends/receives raw bytes
	// - No knowledge of URLs, headers
	// - Great for tunneling / flexibility
	// Why start with TCP? Because tunneling works best when it's not bound to HTTP only — you can tunnel anything (HTTP, SSH, gRPC, etc.) if the base is TCP.
	// Net.Listen func returns a listener and an error
	listener, netListenerError := net.Listen("tcp", "localhost:8000")
	// Listener is like a receptionist at a hotel — just waiting to greet new guests and hand them off to the bellboys (handlers).
	if netListenerError != nil {
		panic(netListenerError)
	}
	// Print the listening address
	println("Listening on: " + listener.Addr().String())
	// Looping for all incoming connections
	for {
		// The listener can be used to accept incoming connections on that TCP port
		// We can use this connection to send and receive data
		// This line blocks — it halts until someone actually connects.
		conn, connectionAcceptError := listener.Accept()
		// When a client connects, it returns a net.Conn — a full duplex stream (read/write) between the server and the client.
		// Conn is like a phone line between you and one client.
		// You can talk, they can talk, and both ends can hang up.
		if connectionAcceptError != nil {
			panic(connectionAcceptError)
		}
		// A go routine is created for each incoming connection ensuring that the server can handle multiple connections at the same time
		// Every time a new user connects, we spawn a goroutine to handle them independently
		// This makes the server non-blocking and able to handle multiple connections simultaneously.
		go handleConnection(conn)
		// This is like spawning a new worker for each client. We don’t wait for one person to finish before helping the next.
	}
}

// The server has now opened a 2-way communication line to a client.
func handleConnection(conn net.Conn) error{
	// Every connection will likely need:
	// A unique ID (for tracking)
	// A timestamp for last heartbeat
	// A channel or buffer to store messages
	// The remote address info (from conn.RemoteAddr())
	defer conn.Close()
	// Deffer is used to ensure that the connection is closed when the function ends
	heartBeatTimeStamp := time.Now()
	// Heartbeat is used to check if the connection is still alive
	// If the time limit exceeds the max timeout, the connection will be closed
	// Setting the max timeout to 1 minute
	// If the connection is idle for 1 minute, the server will close the connection
	maxTimeout := time.Minute
	// The max timeout can be tweaked as per the requirements of the application and the environment
	// This is a go routine that runs in the background
	go handleTimeout(conn, heartBeatTimeStamp, maxTimeout)
	// The handleTimeout function is called to check if the connection is still alive
	// If the connection is idle, the server will close the connection
	// conn.Read() is the gateway into the TCP stream.
	// TCP is a stream. It doesn’t know what a “message” is — just a sequence of bytes.
	// The server may receive:
	// Partial messages (need to stitch them)
	// Multiple messages at once (need to split them)
	// That’s our job to handle at the application layer.
	buffer := make([]byte, 1024)
	// Requires a buffer to read into
	n, readError := conn.Read(buffer)
	if readError != nil {
		if readError == io.EOF {
			// If the connection is closed, the server will close the connection
			conn.Close()
		}
		return readError
	}
	message := string(buffer[:n])
	print(message)
	// Returns how many bytes were read
	// And also an error (which tells if the connection closed, timed out, etc.)
	return nil
}

func handleTimeout(conn net.Conn, heartBeatTimeStamp time.Time, maxTimeout time.Duration) {
	// If the connection is idle for max timeout, the server will close the connection
	if time.Since(heartBeatTimeStamp) > maxTimeout {
		conn.Close()
	}
}
