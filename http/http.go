package http

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"time"

	ntlmssp "github.com/Azure/go-ntlmssp"
)

// type GoProxy struct {
// 	Url string
// }

type GoSSLConfig interface {
	VerifyPeerCertificate(rawCerts []byte) error
}

type GoHttpTransport struct {
	tcpCreater GoTcpDialCreater
	tcpDial    *GoTcpDial
	tlsCreater GoTcpDialCreater
	tlsDial    *GoTcpDial
	Ntlm       bool
	Basic      bool
	Username   string
	Password   string
}

func (transport *GoHttpTransport) SetTlsCreater(creater GoTcpDialCreater) {
	transport.tlsCreater = creater
}

func (transport *GoHttpTransport) SetTcpCreater(creater GoTcpDialCreater) {
	transport.tcpCreater = creater
}

type GoHeaderReader interface {
	ReadHeader() *GoHeader
	HasMore() bool
}

// Go Client
type GoClient struct {
	Transport *GoHttpTransport
	Jar       *GoCookieJar
	Timeout   int
}

func NewGoClient() *GoClient {
	return &GoClient{
		Timeout: 0,
	}
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
				if transport.tcpCreater == nil {
					return net.Dial(network, addr)
				}
				tcpDial, err := transport.tcpCreater.CreateGoDial(addr)
				if err != nil {
					return nil, err
				}
				if tcpDial != nil {
					return tcpDial.dial, nil
				}

				dial, err := net.Dial(network, addr)

				return dial, err
			},
			DialTLS: func(network, addr string) (net.Conn, error) {
				if transport.tlsCreater == nil {
					return tls.Dial(network, addr, nil)
				}
				tcpDial, err := transport.tlsCreater.CreateGoDial(addr)
				if err != nil {
					return nil, err
				}
				if tcpDial != nil {
					return tcpDial.dial, nil
				}
				dial, err := tls.Dial(network, addr, nil)

				return dial, err
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

func Request(goClient *GoClient, request *GoRequest) (*GoResponse, error) {
	client := &http.Client{
		Transport: getTransport(goClient.Transport),
		Jar:       getCookieJar(goClient.Jar),
		Timeout:   time.Duration(goClient.Timeout) * time.Second,
	}

	httpRequest, err := http.NewRequest(request.Method, request.Url, request.postData)

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

	if goClient.Transport != nil {
		if goClient.Transport.Ntlm || goClient.Transport.Basic {
			httpRequest.SetBasicAuth(goClient.Transport.Username, goClient.Transport.Password)
		}
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
