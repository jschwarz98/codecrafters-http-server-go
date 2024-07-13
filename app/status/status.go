package status

const (
	// 2xx
	OK      = "HTTP/1.1 200 OK\r\n"
	CREATED = "HTTP/1.1 201 Created\r\n"
	// 4xx
	NOT_FOUND = "HTTP/1.1 404 Not Found\r\n"
	// 5xx
	INTERNAL_SERVER_ERROR = "HTTP/1.1 500 Internal Server Error\r\n"
)
