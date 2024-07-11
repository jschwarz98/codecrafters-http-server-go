package main

import (
	"errors"
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

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(connection)
	}

}

func handleConnection(connection net.Conn) {
	buffer := make([]byte, 2048)
	connection.Read(buffer)

	requestString := string(buffer)

	verb, path, protocol, headers, body := parseRequest(requestString)

	c, err := responseContent(connection, verb, path, protocol, headers, body)
	if err != nil {
		fmt.Println("error during response generation:", err.Error())
		os.Exit(1)
	}

	connection.Write([]byte(c))
}
func parseRequest(requestString string) (string, string, string, map[string]string, string) {
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

	headers := make(map[string]string)

	for {
		restOfRequest := requestString[i+2:]
		newI := strings.Index(restOfRequest, "\r\n")
		if newI == 0 {
			i += 2
			break
		}
		headerString := requestString[i+2 : i+2+newI]

		headerVals := strings.SplitN(headerString, ": ", 2)
		if len(headerVals) != 2 {
			fmt.Println("Error reading headers: expected key and value, got ", headerString)
			i += newI + 2
			continue
		}
		headers[strings.ToLower(headerVals[0])] = headerVals[1]
		i += newI + 2
	}

	body := requestString[i+2:]

	return verb, path, protocol, headers, body
}

func responseContent(connection net.Conn, verb, path, protocol string, headers map[string]string, body string) (string, error) {
	_ = body
	if verb != "GET" {
		fmt.Println("Error reading request: non GET request detected. only supporting GET right now")
		return "", errors.New("non Get Request detected")
	}
	if protocol != "HTTP/1.1" {
		fmt.Println("Error reading request: only supporting HTTP/1.1 right now")
		return "", errors.New("non HTTP/1.1 Request detected")
	}

	s := ""
	if path == "/" {
		s = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path, "/echo/") {
		restOfPath := path[6:]
		s = plainTextResponse(restOfPath)
	} else if path == "/user-agent" {
		s = plainTextResponse(headers["user-agent"])
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	return s, nil
}

func plainTextResponse(content string) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s)", len(content), content)
}
