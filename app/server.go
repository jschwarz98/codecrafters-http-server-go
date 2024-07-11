package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	pathArg := "./"
	foundPathFlag := false
	for _, arg := range os.Args {
		if foundPathFlag {
			foundPathFlag = false
			pathArg = arg
		}
		if arg == "--directory" {
			foundPathFlag = true
		}
	}

	// TODO pass args and check for path flag. maybe default to ./ otherwhise use the given path
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
		go handleConnection(connection, pathArg)
	}

}

func handleConnection(connection net.Conn, filePath string) {
	buffer := make([]byte, 2048)
	connection.Read(buffer)

	requestString := string(buffer)

	verb, path, protocol, headers, body := parseRequest(requestString)

	c, err := responseContent(connection, filePath, verb, path, protocol, headers, body)
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

func responseContent(connection net.Conn, filePath, verb, path, protocol string, headers map[string]string, body string) (string, error) {
	if protocol != "HTTP/1.1" {
		fmt.Println("Error reading request: only supporting HTTP/1.1 right now")
		return "", errors.New("non HTTP/1.1 Request detected")
	}

	s := "HTTP/1.1 404 Not Found\r\n\r\n"
	if verb == "GET" {
		if path == "/" {
			s = "HTTP/1.1 200 OK\r\n\r\n"
		} else if strings.HasPrefix(path, "/echo/") {
			restOfPath := path[6:]
			s = plainTextResponse(restOfPath)
		} else if path == "/user-agent" {
			s = plainTextResponse(headers["user-agent"])
		} else if strings.HasPrefix(path, "/files/") {
			requestedFile := path[len("/files/"):]
			s = fileReponse(filePath, requestedFile)
		} else {
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	} else if verb == "POST" {
		if strings.HasPrefix(path, "/files/") {
			filename := path[len("/files/"):]
			s = storeFile(filePath, filename, body)
		}
	}
	return s, nil
}

func plainTextResponse(content string) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s)", len(content), content)
}

func fullPath(dir, name string) string {
	fullPath := dir
	if !strings.HasSuffix("/", dir) {
		fullPath += "/"
	}
	fullPath += name
	return fullPath
}

func fileReponse(directory, filename string) string {
	content, err := os.ReadFile(fullPath(directory, filename))
	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s)", len(content), content)
}

func storeFile(path, name, content string) string {
	os.WriteFile(fullPath(path, name), []byte(content), os.FileMode(int(0777)))
	return "HTTP/1.1 201 Created\r\n\r\n"
}
