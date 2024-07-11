package main

import "testing"

func TestParseRequest(t *testing.T) {
	verb, path, protocol, headers, body := parseRequest("GET /user-agent HTTP/1.1\r\nUser-Agent: test-client\r\n\r\nmy body :)")

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
