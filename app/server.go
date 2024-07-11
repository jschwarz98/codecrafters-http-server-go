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
	i := strings.Index(requestString, "\r\n")
	if i == -1 {
		fmt.Println("Error reading request: no next line was found")
		os.Exit(1)
	}

	firstLine := requestString[0:i]
	fields := strings.Fields(firstLine)
	if len(fields) != 3 {
		fmt.Println("Error reading request: expected 3 parameters, got ", len(fields))
		os.Exit(1)
	}

	verb := fields[0]
	path := fields[1]
	protocol := fields[2]

	if verb != "GET" {
		fmt.Println("Error reading request: non GET request detected. only supporting GET right now")
		os.Exit(1)
	}
	if protocol != "HTTP/1.1" {
		fmt.Println("Error reading request: only supporting HTTP/1.1 right now")
		os.Exit(1)
	}

	if path == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo/") {
		restOfPath := path[6:]
		msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s)", len(restOfPath), restOfPath)
		connection.Write([]byte(msg))
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
