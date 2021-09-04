package p2pgrpc

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"

	// "net"

	"google.golang.org/grpc"
)

// Protocol is the GRPC-over-libp2p protocol.
const Protocol protocol.ID = "/grpc/0.0.1"

// GrpcHandler is the GRPC-transported protocol handler.
type GrpcHandler struct {
	ctx        context.Context
	host       host.Host
	grpcServer *grpc.Server
	streamCh   chan network.Stream
}

// NewGrpcHandler attaches the GRPC protocol to a host.
func NewGrpcHandler(ctx context.Context, host host.Host) *GrpcHandler {
	grpcServer := grpc.NewServer()
	ghandler := &GrpcHandler{
		ctx:        ctx,
		host:       host,
		grpcServer: grpcServer,
		streamCh:   make(chan network.Stream),
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
func (p *GrpcHandler) HandleStream(stream network.Stream) {
	select {
	case <-p.ctx.Done():
		return
	case p.streamCh <- stream:
	}
}
