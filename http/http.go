package http

import (
	"io"
	"net"
	"net/http"
)

// type GoProxy struct {
// 	Url string
// }

type GoTcpDial struct {
	Address string
	dial    net.Conn
}

type GoHttpTransport interface {
	TcpDial() *GoTcpDial
}

type GoClient struct {
	Transport GoHttpTransport
	Method    string
	Url       string
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

func getClient(goClient GoClient) *http.Client {
	return &http.Client{}
}

func Request(request GoClient) {
	client := &http.Client{
		Transport: getTransport(request.Transport),
	}

	if request.Method == "GET" {
		_, _ = client.Get(request.Url)
	} else if request.Method == "POST" {

	}
}
