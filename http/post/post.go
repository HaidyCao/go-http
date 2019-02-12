package post

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// GoResetReader GoResetReader
type GoResetReader interface {
	io.ReadCloser
	Reset() error
}

// GoPostStream post stream
type GoPostStream struct {
	ContentLength int64
	body          io.Reader
	GoResetReader
}

// NewPostStream NewPostStream
func NewPostStream(resetReader GoResetReader, contentLength int64) *GoPostStream {
	return &GoPostStream{
		ContentLength: contentLength,
		body:          resetReader,
	}
}

func NewReaderPostStream(reader io.Reader) *GoPostStream {
	return &GoPostStream{
		body: reader,
	}
}

// get reader
func (g *GoPostStream) GetReader() io.Reader {
	return g.body
}

func (g *GoPostStream) Read(p []byte) (int, error) {
	return g.body.Read(p)
}

func (g *GoPostStream) IsCustomReader() bool {
	switch g.body.(type) {
	case *bytes.Buffer:
		return false
	case *bytes.Reader:
		return false
	case *strings.Reader:
		return false
	default:
		return true
	}
}

// Reset Reset
func (g *GoPostStream) Reset() error {
	switch v := g.body.(type) {
	case GoResetReader:
		v.Reset()
		return nil
	default:
		return errors.New("body not support reset")
	}
}

// Close Close
func (g *GoPostStream) Close() error {
	switch v := g.body.(type) {
	case GoResetReader:
		v.Close()
		return nil
	default:
		return errors.New("body not support reset")
	}
}
