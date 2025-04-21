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

	// Set 30s inactivity timeout
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Log connection duration
	duration := time.Now()

	defer func() {
		// Logs & recovers from panics
		if r := recover(); r != nil {
			log.Printf("[!] Recovered from panic with client %s after %s: %v", clientAddr, duration.Round(time.Second), r)
		}
		conn.Close()
		log.Printf("[-] Disconnection from %s (duration %s)", clientAddr, duration.Round(time.Second))  // Log disconnection timestamp
	}()

	buf := make([]byte, 1024)
	for {
		// Read data from client
		n, err := conn.Read(buf)
		if err != nil {
			// Handle read errors:
			switch {
			case err == io.EOF:
				// Client closed connection normally
				log.Printf("[-] Client %s disconnected", clientAddr)
			case errors.Is(err, net.ErrClosed):
				// Connection was closed by server
				log.Printf("[-] Connection %s closed by server", clientAddr)
			case isTimeout(err):
				// Read operation timed out
				log.Printf("[-] Client %s timed out", clientAddr)
			default:
				// Unexpected read error
				log.Printf("[!] Read error from %s: %v", clientAddr, err)
			}
			return
		}

		// Validate received data
		if n == 0 {
			log.Printf("[!] Client %s sent empty payload", clientAddr)
			continue
		}
		// Reset deadline after successful operation
		conn.SetDeadline(time.Now().Add(30 * time.Second))

		//  Safely Echo back
		if _, err = conn.Write(buf[:n]); err != nil {
			// Handle different write errors:
			switch {
			case err == io.EOF:
				// Client disconnected during write
				log.Printf("[-] Client %s disconnected during write", clientAddr)
			case errors.Is(err, net.ErrClosed):
				// Connection was closed
				log.Printf("[-] Connection %s closed during write", clientAddr)
			case isTimeout(err):
				// Write operation timed out
				log.Printf("[-] Write to %s timed out", clientAddr)
			default:
				// Unexpected write error
				log.Printf("[!] Write error to %s: %v", clientAddr, err)
			}
			return
		}
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