package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	// Log connection timestamp and address
	clientAddr := conn.RemoteAddr().String()
	log.Printf("[+] Connection from %s", clientAddr)
	startTime := time.Now()	// Record connection start time

	// Set 30s inactivity timeout
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	
	defer func() {
		duration := time.Since(startTime) // Calculate connection duration

		// Logs & recovers from panics
		if r := recover(); r != nil {
			log.Printf("[!] Recovered from panic with client %s after %s: %v", clientAddr, duration.Round(time.Millisecond), r)
		}
		
		// Close connection and log any errors
		if err := conn.Close(); err != nil {
			log.Printf("[!] Error closing %s: %v", clientAddr, err) 
		}

		// Normal disconnection
		log.Printf("[-] Disconnected %s (duration: %s)", clientAddr, duration.Round(time.Millisecond))
	}()

	const readBufferSize = 1024
	buf := make([]byte, readBufferSize)	// Buffer for reading client data
	
	for {
		// Read data from client
		n, err := conn.Read(buf)
		if err != nil {
			handleDisconnect(clientAddr, err, "read")
			return
		}

		if n == 0 {
			continue // Skip empty packets
		}

		// Reset deadline after successful operation
		conn.SetDeadline(time.Now().Add(30 * time.Second))

		//  Safely Echo back
		if _, err = conn.Write(buf[:n]); err != nil {
			handleDisconnect(clientAddr, err, "write")
			return
		}
	}
}

// Helper func to handle disconnect/error handling 
func handleDisconnect(clientAddr string, err error, op string) {
	switch {
	case err == io.EOF:
		log.Printf("[-] Client %s disconnected", clientAddr)
	case errors.Is(err, net.ErrClosed):
		log.Printf("[-] Connection %s closed", clientAddr)
	case isTimeout(err):
		log.Printf("[-] %s timeout for %s", op, clientAddr)
	default:
		log.Printf("[!] %s error with %s: %v", op, clientAddr, err)
	}
}

// Helper function to check for timeout errors
func isTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func main() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("[!] Failed to start server: %v", err)
	}

	defer listener.Close()

	fmt.Println("Server listening on :4000")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[!] Accept error: %v", err)
				continue
			}
			go handleConnection(conn)
		}
}