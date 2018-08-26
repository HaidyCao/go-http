package http

import (
	"crypto/tls"
	"errors"
	"net"
	"time"
)

type GoTcpDialCreater interface {
	CreateGoDial(address string) (*GoTcpDial, error)
}

type GoTcpDial struct {
	address   string
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
		address: tcpDial.address,
		dial:    tls.Client(originDial, config),
	}, nil
}

func GetGoTcpDial(address string) (*GoTcpDial, error) {
	dial, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &GoTcpDial{
		address: address,
		dial:    dial,
	}, nil
}

func (dial *GoTcpDial) Connect(address string) (bool, error) {
	d, err := net.Dial("tcp", address)
	if err != nil {
		return false, err
	}
	dial.dial = d

	return true, nil
}

func (dial *GoTcpDial) Read(b []byte) (n int, err error) {
	if dial.dial != nil {
		return 0, errors.New("dial not connect")
	}
	return dial.dial.Read(b)
}

func (dial *GoTcpDial) Write(b []byte) (n int, err error) {
	if dial.dial != nil {
		return 0, errors.New("dial not connect")
	}
	return dial.dial.Write(b)
}

func (dial *GoTcpDial) Close() error {
	if dial.dial != nil {
		return errors.New("dial not connect")
	}
	return dial.dial.Close()
}

func (dial *GoTcpDial) SetDeadline(t int64) error {
	if dial.dial != nil {
		return errors.New("dial not connect")
	}
	deadline := time.Now().Add(time.Duration(t) * time.Second)
	return dial.dial.SetDeadline(deadline)
}

func (dial *GoTcpDial) SetReadDeadline(t int64) error {
	if dial.dial != nil {
		return errors.New("dial not connect")
	}
	deadline := time.Now().Add(time.Duration(t) * time.Second)
	return dial.dial.SetReadDeadline(deadline)
}

func (dial *GoTcpDial) SetWriteDeadline(t int64) error {
	if dial.dial != nil {
		return errors.New("dial not connect")
	}
	deadline := time.Now().Add(time.Duration(t) * time.Second)
	return dial.dial.SetReadDeadline(deadline)
}
