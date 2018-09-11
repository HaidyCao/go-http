// +build windows

package http

import "net"

// SilenceSIGPIPE configures the net.Conn in a way that silences SIGPIPEs with
// the SO_NOSIGPIPE socket option.
func SilenceSIGPIPE(c net.Conn) error {
	return nil
}
