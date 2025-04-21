package main

import("fmt"
		"log"
		"net"
)

func handleConnection(conn net.Conn) {
	// Log connection timestamp and address
	clientAddr := conn.RemoteAddr().String()
	log.Printf("[+] Connection from %s", clientAddr)

	defer func() {
		// Log disconnection timestamp
		log.Printf("[-] Disconnection from %s", clientAddr)
		conn.Close()
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("[!] Read error from %s: %v", clientAddr, err)
			return
		}
		// Echo the message back to the client
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("[!] Write error to %s: %v", clientAddr, err)
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