package p2pgrpc

import "net"

// fakeLocalAddr returns a dummy local address.
func fakeLocalAddr() net.Addr {
	localIP := net.ParseIP("127.0.0.1")
	return &net.TCPAddr{IP: localIP, Port: 0}
}

// fakeRemoteAddr returns a dummy remote address.
func fakeRemoteAddr() net.Addr {
	remoteIP := net.ParseIP("127.1.0.1")
	return &net.TCPAddr{IP: remoteIP, Port: 0}
}
