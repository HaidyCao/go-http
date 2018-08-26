package main

import (
	"github.com/HaidyCao/go-http/http"
)

type TlsDial struct {
}

func (dial *TlsDial) CreateGoDial(address string) (*http.GoTcpDial, error) {
	d := &http.GoTcpDial{}
	success, err := d.Connect(address)
	if !success {
		return nil, err
	}

	return http.UpdateDialToSSLTcpDial(d)
}

func main() {
	transport := &http.GoHttpTransport{
		Ntlm: false,
	}
	var tlsCreater http.GoTcpDialCreater
	tlsCreater = &TlsDial{}
	transport.SetTlsCreater(tlsCreater)

	client := &http.GoClient{
		Transport: transport,
	}

	request := &http.GoRequest{
		Url:    "https://www.baidu.com",
		Method: "GET",
	}

	resp, err := http.Request(client, request)

	if err != nil {
		println(err.Error())
		return
	}

	println("responseCode = ", resp.StatueCode)
	str, _ := resp.GetBody().String()
	println(str)
}
