package http

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/axgle/mahonia"
)

// Go Response

type GoResponse struct {
	StatueCode int
	Response   *http.Response
}

func (resp *GoResponse) GetHeader(name string) *GoHeader {
	value := resp.Response.Header[name]

	return &GoHeader{
		Name:  name,
		Count: len(value),
		Value: value,
	}
}

func (resp *GoResponse) GetBody() *GoBody {
	return &GoBody{
		Body: resp.Response.Body,
		read: false,
	}
}

func (resp *GoResponse) Read(buffer []byte) (int, error) {
	if resp.Response.Body == nil {
		return 0, errors.New("Reader is nil")
	}

	return resp.Response.Body.Read(buffer)
}

// GoBody

type GoBody struct {
	Body io.ReadCloser
	Byte []byte
	read bool
}

func newGoBody(body io.ReadCloser) GoBody {
	return GoBody{
		Body: body,
		Byte: make([]byte, 0),
		read: false,
	}
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func (body *GoBody) String(charset string) (string, error) {
	if body.read {
		ret := string(body.Byte)
		if charset == "" || strings.EqualFold(charset, "utf-8") {
			return ConvertToString(ret, charset, "utf-8"), nil
		}
		return ret, nil
	}

	body.read = true
	buffer := make([]byte, 1024*4)
	str := make([]byte, 0)
	for {
		len, err := body.Body.Read(buffer)
		if err != nil {
			return "", err
		}
		str = append(str, buffer[:len]...)
		if len == -1 {
			break
		}
	}

	body.Byte = str
	// dec := mahonia.NewDecoder("gbk")

	return string(str), nil
}
