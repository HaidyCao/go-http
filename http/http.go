package http

import (
	"crypto/tls"
	"crypto/x509"
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

type GoHttpTransport interface {
	TcpDial() *GoTcpDial
}

type GoHeader struct {
	Name  string
	Value string
}

type GoHeaderReader interface {
	ReadHeader() *GoHeader
	HasMore() bool
}

type GoClient struct {
	Transport   GoHttpTransport
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

type GoResponse struct {
	StatueCode int64
	Headers    map[string]string
	Reader     io.Reader
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

func updateDialToSSLTcpDial(tcpDial *GoTcpDial) (*GoTcpDial, error) {
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

// TODO tcp 协议升级

// func getProxy(proxy GoProxy) *http.Pro

func getTransport(transport GoHttpTransport) *http.Transport {
	return &http.Transport{
		Dial: func(network string, addr string) (net.Conn, error) {
			tcpDial := transport.TcpDial()
			if tcpDial.Address != "" {
				return net.Dial("tcp", tcpDial.Address)
			}
			return net.Dial("tcp", tcpDial.Address)
		},
	}
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
			httpRequest.Header.Add(header.Name, header.Value)
		}
	}

	if request.ContentType != "" {
		httpRequest.Header.Set("Content-Type", request.ContentType)
	}

	if err != nil {
		return nil, err
	}

	_, _ = client.Do(httpRequest)

	return nil, nil
}
