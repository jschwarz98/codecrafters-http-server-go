package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func internalServerError() string {
	return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
}

func responseContent(filePath, verb, path, protocol string, headers map[string]string, body string) (string, error) {
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
			s = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
	} else if verb == "POST" {
		if strings.HasPrefix(path, "/files/") {
			filename := path[len("/files/"):]
			s = storeFile(filePath, filename, body)
		} else {
			s = "HTTP/1.1 404 Not Found\r\n\r\n"
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
