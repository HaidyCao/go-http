package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"

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
	data []byte
	read bool
}

func newGoBody(body io.ReadCloser) GoBody {
	return GoBody{
		Body: body,
		data: make([]byte, 0),
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

func (body *GoBody) GetData() ([]byte, error) {
	b, err := ioutil.ReadAll(body.Body)
	return b, err
}

func (body *GoBody) String() (string, error) {
	data, err := body.GetData()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (body *GoBody) Close() error {
	if body.Body != nil {
		return body.Body.Close()
	}
	return nil
}
