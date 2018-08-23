package http

import (
	"crypto/tls"
	"net"
	"time"
)

type GoTcpDial struct {
	Address   string
	SSLConfig GoSSLConfig
	dial      net.Conn
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

func (dial *GoTcpDial) Read(b []byte) (n int, err error) {
	return dial.dial.Read(b)
}

func (dial *GoTcpDial) Write(b []byte) (n int, err error) {
	return dial.dial.Write(b)
}

func (dial *GoTcpDial) Close() error {
	return dial.dial.Close()
}

func (dial *GoTcpDial) SetDeadline(t int64) error {
	deadline := time.Now().Add(time.Duration(t) * time.Second)
	return dial.dial.SetDeadline(deadline)
}

func (dial *GoTcpDial) SetReadDeadline(t int64) error {
	deadline := time.Now().Add(time.Duration(t) * time.Second)
	return dial.dial.SetReadDeadline(deadline)
}

func (dial *GoTcpDial) SetWriteDeadline(t int64) error {
	deadline := time.Now().Add(time.Duration(t) * time.Second)
	return dial.dial.SetReadDeadline(deadline)
}
