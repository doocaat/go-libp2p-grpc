package p2pgrpc

import (
	"github.com/libp2p/go-libp2p-core/network"
	netaddr "github.com/multiformats/go-multiaddr/net"
	"net"
)

// streamConn represents a net.Conn wrapped to be compatible with net.conn
type streamConn struct {
	network.Stream
}

// LocalAddr returns the local address.
func (c *streamConn) LocalAddr() net.Addr {
	addr, err := netaddr.ToNetAddr(c.Stream.Conn().LocalMultiaddr())
	if err != nil {
		return fakeLocalAddr()
	}
	return addr
}

// RemoteAddr returns the remote address.
func (c *streamConn) RemoteAddr() net.Addr {
	addr, err := netaddr.ToNetAddr(c.Stream.Conn().RemoteMultiaddr())
	if err != nil {
		return fakeRemoteAddr()
	}
	return addr
}

var _ net.Conn = &streamConn{}
