package http

import (
	"bytes"
	"io"
	"strings"

	"github.com/HaidyCao/go-http/http/post"
)

type GoRequest struct {
	Method      string
	Url         string
	postData    *post.GoPostStream
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
	req.postData = post.NewReaderPostStream(bytes.NewReader(data))
	if req.Method == "" {
		req.Method = "POST"
	}
}

func (req *GoRequest) RequestBodyString(data string) {
	req.postData = post.NewReaderPostStream(strings.NewReader(data))
	if req.Method == "" {
		req.Method = "POST"
	}
}

func (req *GoRequest) RequestPostStream(stream *post.GoPostStream) {
	req.postData = stream
}
