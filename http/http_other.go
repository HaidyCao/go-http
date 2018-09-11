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
	fd := int(reflect.ValueOf(c).Elem().FieldByName("fd").Elem().FieldByName("sysfd").Int())
	if runtime.GOOS != "windows" {
		return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_NOSIGPIPE, 1)
	} else {
		return nil
	}
}
