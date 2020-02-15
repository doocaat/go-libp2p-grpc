package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/golang/protobuf/jsonpb"
	"github.com/paralin/go-libp2p-grpc"
	"github.com/paralin/go-libp2p-grpc/examples/echo/echosvc"
	"google.golang.org/grpc"
)

// Echoer implements the EchoService.
type Echoer struct {
	PeerID peer.ID
}

// Echo asks a node to respond with a message.
func (e *Echoer) Echo(ctx context.Context, req *echosvc.EchoRequest) (*echosvc.EchoReply, error) {
	log.Debugf("%+v", req)
	return &echosvc.EchoReply{
		Message: req.GetMessage(),
		PeerId:  e.PeerID.Pretty(),
	}, nil
}

var (
	log = golog.Logger("echo")
)

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetLogLevel("echo", "debug") // Change to DEBUG for extra info

	// Parse options from the command line
	port := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	echoMsg := flag.String("m", "Hello, world", "message to echo")
	flag.Parse()

	if *port == 0 {
		panic("Please specify server port with -l")
	}
	if *target == "" {
		h := server(*port)
		log.Debug("listening for connections")
		log.Debug("copy below to run a client:")
		log.Debugf("./echo -l %d -d %s -m \"hello echo\\!\"", *port, h.ID().Pretty())
		select {} // hang forever
	} else {
		srvAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", *port))
		if err != nil {
			panic(err)
		}
		srvID, err := peer.Decode(*target)
		if err != nil {
			panic(fmt.Sprintf("parse peer id %s error: %w", *target, err))
		}
		client(srvID, srvAddr, *echoMsg)
	}
}

func server(port int) host.Host {
	// Make a host that listens on the given multiaddress
	h, err := libp2p.New(context.Background(), libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port),
	))
	if err != nil {
		panic(err)
	}
	// Set the grpc protocol handler on it
	grpcProto := p2pgrpc.NewGRPCProtocol(context.Background(), h)

	// Register our echoer GRPC service.
	echosvc.RegisterEchoServiceServer(grpcProto.GetGRPCServer(), &Echoer{PeerID: h.ID()})
	return h
}

func client(srvID peer.ID, srvAddr ma.Multiaddr, echoMsg string) host.Host {
	// Make a host
	h, err := libp2p.New(context.Background())
	if err != nil {
		panic(err)
	}
	// add server infomation in peerstore
	ps := h.Peerstore()
	ps.AddAddr(srvID, srvAddr, 10*time.Second)

	// Set the grpc protocol handler on it
	grpcProto := p2pgrpc.NewGRPCProtocol(context.Background(), h)

	// make a new stream from host B to host A
	log.Debug("dialing via grpc")
	grpcConn, err := grpcProto.Dial(context.Background(), srvID, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	// create our service client
	log.Debug("new grpc client")
	echoClient := echosvc.NewEchoServiceClient(grpcConn)
	echoReply, err := echoClient.Echo(context.Background(), &echosvc.EchoRequest{Message: echoMsg})
	if err != nil {
		panic(err)
	}

	log.Debug("read reply:")
	err = (&jsonpb.Marshaler{EmitDefaults: true, Indent: "\t"}).
		Marshal(os.Stdout, echoReply)
	if err != nil {
		panic(err)
	}
	return h
}
