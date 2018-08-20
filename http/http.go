package http

import (
	"io"
	"net/http"
)

type Request interface {
	Url(string)
	Headers([]string)
	Client(http.Client)
}

type Response struct {
	StatueCode int64
	Headers map[string]string
	Reader io.Reader
}