package main

import (
	"testing"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
)

func TestParseRequest(t *testing.T) {
	verb, path, protocol, headers, body, err := request.ParseRequest("GET /user-agent HTTP/1.1\r\nUser-Agent: test-client\r\n\r\nmy body :)")

	if err != nil {
		t.Errorf("should not error out: %s", err.Error())
	}

	if verb != "GET" {
		t.Errorf("verb should be GET but got %s", verb)
	}

	if path != "/user-agent" {
		t.Errorf("path should be /user-agent but got %s", path)
	}

	if protocol != "HTTP/1.1" {
		t.Errorf("protocol should be HTTP/1.1 but got %s", protocol)
	}
	if body != "my body :)" {
		t.Errorf("body should be 'my body :)' but got %s", body)
	}

	if headers["user-agent"] != "test-client" {
		t.Errorf("user agent should be test-client but got %s", headers["user-agent"])
	}

}
