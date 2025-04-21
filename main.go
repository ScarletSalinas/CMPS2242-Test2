package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)
var port = flag.Int("port", 4000, "Port number to listen on")

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
		
		// Clean input
		clean := strings.TrimSpace(line)

		// Handle Command input
		switch {
		case strings.HasPrefix(clean, "/time"):
			currentTime := time.Now().Format("2006-01-02 15:04:05 CST")	// Get current time
			conn.Write([]byte("Server time: " + currentTime + "\n"))
			continue  // Maintain connection
		
		case strings.HasPrefix(clean, "/quit"):
			conn.Write([]byte("Closing connection...\n"))
			return //close connection

		case strings.HasPrefix(clean, "/echo "):
			message := strings.TrimPrefix(clean, "/echo ") // Remove "/echo" from input
			conn.Write([]byte(message + "\n"))
			continue // Maintain conn.

		default:

			// Handle special input 
			switch strings.ToLower(clean) {
			case "hello":
				if _, err := conn.Write([]byte("Hi, there!\n")); err != nil {
					handleDisconnect(clientAddr, err, "write")
				}
				continue // Maintains connection

			case "":
				if _, err := conn.Write([]byte("Say something...\n")); err != nil {
					handleDisconnect(clientAddr, err, "write")
				}
				continue // Maintains connection

			case "bye":
				if _, err := conn.Write([]byte("So long, and thanks for all the fish!\n")); err != nil {
					handleDisconnect(clientAddr, err, "write")
				}
				return // Disconnects

			default:
				truncated := false
				// Handle normal messages with truncation
				if len(clean) > maxMessageSize {
					clean = clean[:maxMessageSize]
					truncated = true
				}

				if truncated {
					log.Printf("[!] Truncated message from %s", clientAddr)
					if _, err := conn.Write([]byte("[!] Message truncated\n")); err != nil {
						handleDisconnect(clientAddr, err, "write")
						return
					}
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
	}
}

// Helper func to handle disconnect/error handling 
func handleDisconnect(clientAddr string, err error, op string) {
	switch {
	case err == io.EOF:
		return
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
	flag.Parse()
	
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("[!] Failed to start server: %v", err)
	}

	defer listener.Close()

	log.Printf("Server listening on :%d", *port)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[!] Accept error: %v", err)
				continue
			}
			go handleConnection(conn)
		}
}