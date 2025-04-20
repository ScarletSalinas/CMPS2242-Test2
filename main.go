package main

import("fmt"
		"net"
)

// Use listen method when creating a server
// Listen on port 4000
// Make port 4000 open
// Listen for tcp 
// if tcp package sent, recieve package

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		// Echo the message back to the client
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}
}

func main() {
	// Define the target host and port we want to connect to
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	fmt.Println("Server listening on :4000")

	// Our program runs an infinite loop
	//cant come back to first loop
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting :", err)
				continue
			}
			go handleConnection(conn)
		}
}