// +build !windows

package http

import (
	"net"
	"reflect"
	"runtime"
	"syscall"
)

// SilenceSIGPIPE configures the net.Conn in a way that silences SIGPIPEs with
// the SO_NOSIGPIPE socket option.
// SilenceSIGPIPE configures the net.Conn in a way that silences SIGPIPEs with
// the SO_NOSIGPIPE socket option.
func SilenceSIGPIPE(c net.Conn) error {
	// use reflection until https://github.com/golang/go/issues/9661 is fixed
	netFDField := reflect.ValueOf(c).Elem()
	fd := int(netFDField.FieldByName("fd").Elem().FieldByName("pfd").FieldByName("Sysfd").Int())
	if runtime.GOOS != "windows" {
		return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_NOSIGPIPE, 1)
	} else {
		return nil
	}
}
