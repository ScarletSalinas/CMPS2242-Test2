package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

	const maxMessageSize = 1024
	reader := bufio.NewReader(conn)	
		for {
			// Read until newline
			line, err := reader.ReadString('\n')
			if err != nil {
				handleDisconnect(clientAddr, err, "read")
				return
			}
	
			// Clean the input
			original := line
			truncated := false
			if len(original) > maxMessageSize {
				original = original[:maxMessageSize]
				truncated = true
			}
			clean := strings.TrimSpace(original)
	
			if truncated {
				log.Printf("[!] Truncated message from %s to %d bytes", clientAddr, maxMessageSize)
				_, _ = conn.Write([]byte("[!] Warning: message too long. Truncated to 1024 bytes.\n"))
			}

			// Reset deadline after successful operation
			conn.SetDeadline(time.Now().Add(30 * time.Second))

			//  Safely Echo back + newline
			if _, err = conn.Write([]byte(clean + "\n")); err != nil {
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