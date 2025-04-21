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

	// Set 30s inactivity tiemout
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	defer func() {
		// Logs & recovers from panics
		if r := recover(); r != nil {
			log.Printf("[!] Recovered form panic with client %s: %v", clientAddr, r)
		}
		conn.Close()
		log.Printf("[-] Disconnection from %s", clientAddr)  // Log disconnection timestamp
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("[-] Client %s disconnected", clientAddr)
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("[-] Client %s timed out", clientAddr)
			} else {
				log.Printf("[!] Read error from %s: %v", clientAddr, err)
			}
			return
		}

		// Validate recieved data
		if n == 0 {
			log.Printf("[!] Client %s sent empty payload", clientAddr)
			continue
		}
		// Reset deadline after successful operation
		conn.SetDeadline(time.Now().Add(30 * time.Second))

		//  Safely Echo back
		if _, err = conn.Write(buf[:n]); err != nil {
			if err == io.EOF || errors.Is(err, net.ErrClosed) {
				log.Printf("[-] Client %s disconnected during write", clientAddr)	
			} else {
				log.Printf("[!] Write error to %s: %v", clientAddr, err)
			}
			return
		}
	}
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