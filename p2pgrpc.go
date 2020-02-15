package p2pgrpc

import (
	"context"
	// "net"

	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"google.golang.org/grpc"
)

// Protocol is the GRPC-over-libp2p protocol.
const Protocol protocol.ID = "/grpc/0.0.1"

// GrpcHandler is the GRPC-transported protocol handler.
type GrpcHandler struct {
	ctx        context.Context
	host       host.Host
	grpcServer *grpc.Server
	streamCh   chan inet.Stream
}

// NewGrpcHandler attaches the GRPC protocol to a host.
func NewGrpcHandler(ctx context.Context, host host.Host) *GrpcHandler {
	grpcServer := grpc.NewServer()
	ghandler := &GrpcHandler{
		ctx:        ctx,
		host:       host,
		grpcServer: grpcServer,
		streamCh:   make(chan inet.Stream),
	}
	host.SetStreamHandler(Protocol, ghandler.HandleStream)
	// Serve will not return until Accept fails, when the ctx is canceled.
	go grpcServer.Serve(newGrpcListener(ghandler))
	return ghandler
}

// GetGRPCServer returns the grpc server.
func (p *GrpcHandler) GetGRPCServer() *grpc.Server {
	return p.grpcServer
}

// HandleStream handles an incoming stream.
func (p *GrpcHandler) HandleStream(stream inet.Stream) {
	select {
	case <-p.ctx.Done():
		return
	case p.streamCh <- stream:
	}
}
