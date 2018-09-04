package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"reflect"
	"syscall"
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
	transport *http.Transport
}

func (client *GoClient) setTransport(transport *http.Transport) {
	client.transport = transport
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

// SilenceSIGPIPE configures the net.Conn in a way that silences SIGPIPEs with
// the SO_NOSIGPIPE socket option.
func SilenceSIGPIPE(c net.Conn) error {
	// use reflection until https://github.com/golang/go/issues/9661 is fixed
	fd := int(reflect.ValueOf(c).Elem().FieldByName("fd").Elem().FieldByName("sysfd").Int())
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_NOSIGPIPE, 1)
}

// NoSIGPIPEDialer returns a dialer that won't SIGPIPE should a connection
// actually SIGPIPE. This prevents the debugger from intercepting the signal
// even though this is normal behaviour.
type NoSIGPIPEDialer net.Dialer

func (d *NoSIGPIPEDialer) handle(c net.Conn, err error) (net.Conn, error) {
	if err != nil {
		return nil, err
	}
	if err := SilenceSIGPIPE(c); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func (d *NoSIGPIPEDialer) Dial(network, address string) (net.Conn, error) {
	c, err := (*net.Dialer)(d).Dial(network, address)
	return d.handle(c, err)
}

func (d *NoSIGPIPEDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, err := (*net.Dialer)(d).DialContext(ctx, network, address)
	return d.handle(c, err)
}

func getTransport(transport *GoHttpTransport) http.RoundTripper {
	if transport != nil {
		ret := &http.Transport{
			DialContext: (&NoSIGPIPEDialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
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
