package p2p

//TODO: how to strictly allow subscribe
// must include signature
// monitor and add to blacklist
//

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/crypto"
	p2pPubSub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	p2pGrpc "github.com/paralin/go-libp2p-grpc"
	"log"
)

type Node struct {
	Self       Peer
	Host       host.Host
	Pubsub     *p2pPubSub.PubSub
	GrpcClient *GRPCService_Client
	GrpcServer *GRPCService_Server
	Blockchain blockchainInf
}

type NodeConfig struct {
	PublicIP   string
	Port       int
	PrivateKey crypto.PrivKey
	Blockchain blockchainInf
}

func NewNode(config NodeConfig) *Node {

	listenAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.PublicIP, config.Port))
	catchError(err)

	ctx := context.Background()
	p2pHost, err := libp2p.New(ctx,
		libp2p.ListenAddrs(listenAddr), libp2p.Identity(config.PrivateKey),
	)
	catchError(err)

	selfPeer := Peer{
		PeerID:        p2pHost.ID(),
		IP:            config.PublicIP,
		Port:          config.Port,
		TargetAddress: append([]ma.Multiaddr{}, listenAddr),
	}

	//create pubsub protocol
	pubsub, err := p2pPubSub.NewGossipSub(ctx, p2pHost)
	catchError(err)

	//create grpc libp2p protocol
	p2pgrpc := p2pGrpc.NewGRPCProtocol(context.Background(), p2pHost)

	//create grpc server
	grpcServer := &GRPCService_Server{Blockchain: config.Blockchain}
	grpcServer.registerServices(p2pgrpc.GetGRPCServer())

	//create grpc client
	grpcClient := &GRPCService_Client{
		p2pgrpc: p2pgrpc,
	}
	return &Node{Host: p2pHost, Self: selfPeer, Pubsub: pubsub, GrpcClient: grpcClient, GrpcServer: grpcServer, Blockchain: config.Blockchain}
}

func (node *Node) Subscribe(topic string, handler func(*p2pPubSub.Subscription)) error {
	subscription, err := node.Pubsub.Subscribe(topic)
	if err != nil {
		log.Println(err)
		return err
	}
	go handler(subscription)
	return nil
}

func (node *Node) Publish(topic string, data []byte) error {
	return node.Pubsub.Publish(topic, data)
}
