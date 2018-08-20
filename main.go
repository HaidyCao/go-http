package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
)

func main() {
	config := &tls.Config{
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			println(base64.StdEncoding.EncodeToString(rawCerts[0]))
			println(base64.StdEncoding.EncodeToString(verifiedChains[0][0].Raw))
			return nil
		},
	}

	dial, _ := tls.Dial("tcp", "www.baidu.com:443", config)
	println(dial.ConnectionState)
}
