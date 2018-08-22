package http

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
)

// type GoProxy struct {
// 	Url string
// }

type GoSSLConfig interface {
	VerifyPeerCertificate(rawCerts []byte) error
}

type GoTcpDial struct {
	Address   string
	SSLConfig GoSSLConfig
	dial      net.Conn
}

type GoHttpTransport struct {
	TcpDial *GoTcpDial
}

type GoHeaderReader interface {
	ReadHeader() *GoHeader
	HasMore() bool
}

type GoClient struct {
	Transport *GoHttpTransport

	Method      string
	Url         string
	ContentType string
	Body        GoIoReader
	body        io.Reader
	headers     []*GoHeader
}

func (c *GoClient) AddHeader(header *GoHeader) {
	c.headers = append(c.headers, header)
}

type GoIoReader interface {
	Read(buffer []byte) (int, error)
}

// GoHeader

type GoHeader struct {
	Name  string
	Count int
	Value []string
}

func (h *GoHeader) GetHeader(index int) (string, error) {
	if index >= h.Count {
		return "", errors.New("out of range")
	}
	return h.Value[index], nil
}

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

func (resp *GoResponse) Read(buffer []byte) (int, error) {
	if resp.Response.Body == nil {
		return 0, errors.New("Reader is nil")
	}

	return resp.Response.Body.Read(buffer)
}

func GetGoTcpDial(address string) (*GoTcpDial, error) {
	dial, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &GoTcpDial{
		Address: address,
		dial:    dial,
	}, nil
}

func getVerifyPeerCertificateFunc(sslConfig GoSSLConfig) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if sslConfig == nil {
		return nil
	}
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		for i := 0; i < len(rawCerts); i++ {
			err := sslConfig.VerifyPeerCertificate(rawCerts[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func UpdateDialToSSLTcpDial(tcpDial *GoTcpDial) (*GoTcpDial, error) {
	originDial := tcpDial.dial

	config := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: getVerifyPeerCertificateFunc(tcpDial.SSLConfig),
	}
	return &GoTcpDial{
		Address: tcpDial.Address,
		dial:    tls.Client(originDial, config),
	}, nil
}

func getTransport(transport *GoHttpTransport) http.RoundTripper {
	if transport != nil {
		return &http.Transport{
			Dial: func(network string, addr string) (net.Conn, error) {
				tcpDial := transport.TcpDial
				if tcpDial.Address != "" {
					return net.Dial("tcp", tcpDial.Address)
				}
				return net.Dial("tcp", tcpDial.Address)
			},
		}
	}
	return http.DefaultTransport
}

func getClient(goClient *GoClient) *http.Client {
	return &http.Client{}
}

func Request(request *GoClient) (*GoResponse, error) {
	client := &http.Client{
		Transport: getTransport(request.Transport),
	}

	httpRequest, err := http.NewRequest(request.Method, request.Url, nil)

	if request.headers != nil {
		for i := 0; i < len(request.headers); i++ {
			header := request.headers[i]

			for j := 0; j < len(header.Value); j++ {
				httpRequest.Header.Add(header.Name, header.Value[j])
			}
		}
	}

	if request.ContentType != "" {
		httpRequest.Header.Set("Content-Type", request.ContentType)
	}

	if err != nil {
		return nil, err
	}

	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return &GoResponse{
		StatueCode: response.StatusCode,
		Response:   response,
	}, nil
}
