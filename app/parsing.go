package main

import (
	"errors"
	"fmt"
	"strings"
)

func parseRequest(requestString string) (string, string, string, map[string]string, string, error) {
	headers := make(map[string]string)
	i := strings.Index(requestString, "\r\n")

	if i == -1 {
		fmt.Println("Error reading request: no next line was found")
		return "", "", "", headers, "", errors.New("no next line")
	}
	firstLine := requestString[0:i]
	fields := strings.Fields(firstLine)
	if len(fields) != 3 {
		fmt.Println("Error reading request: expected 3 parameters, got ", len(fields))
		return "", "", "", headers, "", errors.New("not 3 req params in first line")
	}

	verb := fields[0]
	path := fields[1]
	protocol := fields[2]

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
	endIndex := strings.Index(body, "\x00")
	if endIndex != -1 {
		body = body[0:endIndex]
	}

	return verb, path, protocol, headers, body, nil
}
