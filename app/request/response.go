package request

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/status"
)

func InternalServerError() string {
	return status.INTERNAL_SERVER_ERROR + "\r\n" // no headers
}

func ResponseContent(filePath, verb, path, protocol string, headers map[string]string, body string) (string, error) {
	if protocol != "HTTP/1.1" {
		fmt.Println("Error reading request: only supporting HTTP/1.1 right now")
		return "", errors.New("non HTTP/1.1 Request detected")
	}

	s := status.NOT_FOUND
	responseHeaders := make(map[string]string)

	fmt.Println("accept encoding: ", headers[status.ACCEPT_ENCODING])
	accepts := strings.Split(headers[status.ACCEPT_ENCODING], ",")
	for _, a := range accepts {
		fmt.Println(a)
		if strings.TrimSpace(a) == "gzip" {
			responseHeaders[status.CONTENT_ENCODING] = "gzip"
		}
	}
	resBody := ""
	if verb == "GET" {
		if path == "/" {
			s = status.OK
		} else if strings.HasPrefix(path, "/echo/") {
			restOfPath := path[len("/echo/"):]
			s, responseHeaders, resBody = plainTextResponse(restOfPath, responseHeaders)
		} else if path == "/user-agent" {
			s, responseHeaders, resBody = plainTextResponse(headers["user-agent"], responseHeaders)
		} else if strings.HasPrefix(path, "/files/") {
			requestedFile := path[len("/files/"):]
			s, responseHeaders, resBody = fileReponse(filePath, requestedFile, responseHeaders)
		} else {
			s = status.NOT_FOUND
		}
	} else if verb == "POST" {
		if strings.HasPrefix(path, "/files/") {
			filename := path[len("/files/"):]
			s = storeFile(filePath, filename, body)
		} else {
			s = status.NOT_FOUND
		}
	}
	if len(resBody) > 0 {
		responseHeaders[status.CONTENT_LENGTH] = fmt.Sprintf("%d", len(resBody))
	}

	response := s

	fmt.Println("Response Headers:", responseHeaders)
	fmt.Println("Response Body:", resBody)
	for key := range responseHeaders {
		val := responseHeaders[key]
		if val != "" {
			response += key + ": " + val + "\r\n"
		}
	}
	response += "\r\n"
	response += resBody

	return response, nil
}

func plainTextResponse(content string, responseHeaders map[string]string) (string, map[string]string, string) {
	responseHeaders[status.CONTENT_TYPE] = "text/plain"
	if responseHeaders[status.CONTENT_ENCODING] == "gzip" {
		// TODO
		content = encodeGZIP(content)
	}
	return status.OK, responseHeaders, content
}

func fullPath(dir, name string) string {
	fullPath := dir
	if !strings.HasSuffix("/", dir) {
		fullPath += "/"
	}
	fullPath += name
	return fullPath
}

func fileReponse(directory, filename string, responseHeaders map[string]string) (string, map[string]string, string) {
	c, err := os.ReadFile(fullPath(directory, filename))
	if err != nil {
		return status.NOT_FOUND, responseHeaders, ""
	}
	content := string(c)
	responseHeaders[status.CONTENT_TYPE] = "application/octet-stream"

	if responseHeaders[status.CONTENT_ENCODING] == "gzip" {
		content = encodeGZIP(content)
	}

	return status.OK, responseHeaders, content
}
func storeFile(path, name, content string) string {
	os.WriteFile(fullPath(path, name), []byte(content), os.FileMode(int(0777)))
	return status.CREATED
}

func encodeGZIP(content string) string {
	// TODO
	return content
}
