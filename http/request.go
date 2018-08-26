package http

import (
	"bytes"
	"io"
	"strings"
)

type GoRequest struct {
	Method      string
	Url         string
	postData    io.Reader
	ContentType string
	body        io.Reader
	headers     []*GoHeader
}

func (c *GoRequest) AddHeader(header *GoHeader) {
	c.headers = append(c.headers, header)
}

func (c *GoRequest) AddHeaderNameAndValue(name string, value string) {
	header := &GoHeader{
		Name:  name,
		Value: append(make([]string, 0), value),
	}
	c.headers = append(c.headers, header)
}

func NewRequest(url string) *GoRequest {
	return &GoRequest{
		Url:    url,
		Method: "GET",
	}
}

func (req *GoRequest) RequestBodyBytes(data []byte) {
	req.postData = bytes.NewReader(data)
	if req.Method == "" {
		req.Method = "POST"
	}
}

func (req *GoRequest) RequestBodyString(data string) {
	req.postData = strings.NewReader(data)
	if req.Method == "" {
		req.Method = "POST"
	}
}
