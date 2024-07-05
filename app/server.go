package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 2048)
	connection.Read(buffer)

	requestString := string(buffer)
	if !strings.HasPrefix(requestString, "GET / HTTP/1.1") {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	response := "HTTP/1.1 200 OK\r\n\r\n"
	connection.Write([]byte(response))
}
