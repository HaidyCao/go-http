package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"

	ntlmssp "github.com/Azure/go-ntlmssp"
)

// type GoProxy struct {
// 	Url string
// }

type GoSSLConfig interface {
	VerifyPeerCertificate(rawCerts []byte) error
}

type GoHttpTransport struct {
	TcpDial  *GoTcpDial
	TlsDial  *GoTcpDial
	Ntlm     bool
	Basic    bool
	Username string
	Password string
}

type GoHeaderReader interface {
	ReadHeader() *GoHeader
	HasMore() bool
}

// Go Client
type GoClient struct {
	Transport    *GoHttpTransport
	Jar          *GoCookieJar
	Method       string
	Url          string
	ContentType  string
	body         io.Reader
	headers      []*GoHeader
	PostData     []byte
	ConnTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

func NewGoClient() *GoClient {

	return &GoClient{
		ConnTimeout:  30,
		ReadTimeout:  30,
		WriteTimeout: 30,
	}
}

func (c *GoClient) AddHeader(header *GoHeader) {
	c.headers = append(c.headers, header)
}

func (c *GoClient) AddHeaderNameAndValue(name string, value string) {
	header := &GoHeader{
		Name:  name,
		Value: append(make([]string, 0), value),
	}
	c.headers = append(c.headers, header)
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

func getTransport(transport *GoHttpTransport) http.RoundTripper {
	if transport != nil {
		ret := &http.Transport{
			Dial: func(network string, addr string) (net.Conn, error) {
				tcpDial := transport.TcpDial
				if tcpDial != nil {
					return tcpDial.dial, nil
				}
				return net.Dial(network, addr)
			},
			DialTLS: func(network, addr string) (net.Conn, error) {
				tcpDial := transport.TlsDial
				if tcpDial != nil {
					return tcpDial.dial, nil
				}
				return tls.Dial(network, addr, nil)
			},
		}

		if transport.Ntlm {
			return ntlmssp.Negotiator{
				RoundTripper: ret,
			}
		}
		return ret
	}
	return http.DefaultTransport
}

func getCookieJar(cookieJar *GoCookieJar) http.CookieJar {
	if cookieJar != nil {
		return cookieJar.Jar
	}
	return nil
}

func getClient(goClient *GoClient) *http.Client {
	return &http.Client{}
}

func Request(request *GoClient) (*GoResponse, error) {
	client := &http.Client{
		Transport: getTransport(request.Transport),
		Jar:       getCookieJar(request.Jar),
	}

	var postData io.Reader
	if request.PostData != nil {
		postData = bytes.NewReader(request.PostData)
	} else {
		postData = nil
	}

	httpRequest, err := http.NewRequest(request.Method, request.Url, postData)

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

	if request.Transport != nil {
		if request.Transport.Ntlm || request.Transport.Basic {
			httpRequest.SetBasicAuth(request.Transport.Username, request.Transport.Password)
		}
	}

	if request.PostData != nil {
		request.body = bytes.NewReader(request.PostData)
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
