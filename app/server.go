package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
)

const BUFFER_SIZE = 8192

func main() {
	pathArg := getPathFromArgs()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(connection, pathArg)
	}

}

func handleConnection(connection net.Conn, filePath string) {
	buffer := make([]byte, BUFFER_SIZE)
	connection.Read(buffer)

	requestString := string(buffer)

	fmt.Println("parsing request")

	verb, path, protocol, headers, body, err := request.ParseRequest(requestString)
	if err != nil {
		connection.Write([]byte(request.InternalServerError()))
		return
	}

	fmt.Println("generating content", verb, path, protocol, headers, body)
	c, err := request.ResponseContent(filePath, verb, path, protocol, headers, body)
	if err != nil {
		connection.Write([]byte(request.InternalServerError()))
		return
	}

	fmt.Println("got response", c)

	connection.Write([]byte(c))
}

func getPathFromArgs() string {
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
	return pathArg
}
