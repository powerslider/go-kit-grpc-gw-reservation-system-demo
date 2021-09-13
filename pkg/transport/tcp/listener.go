package tcp

import (
	"crypto/tls"
	"net"
)

func NewTCPListener(addr string, tlsConfig ...*tls.Config) (net.Listener, error) {
	var listener net.Listener
	var err error
	if len(tlsConfig) > 0 {
		listener, err = tls.Listen("tcp", addr, tlsConfig[0])
	} else {
		listener, err = net.Listen("tcp", addr)
	}
	if err != nil {
		return nil, err
	}
	return listener, nil
}
