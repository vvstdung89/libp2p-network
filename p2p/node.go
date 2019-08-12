package p2p

import (
	"context"
	"fmt"
	"log"
	"node/p2p/chunk"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/crypto"
	p2pPubSub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	p2pGrpc "github.com/paralin/go-libp2p-grpc"

	kaddht "github.com/libp2p/go-libp2p-kad-dht"
)

type Node struct {
	Self         Peer
	Host         host.Host
	Pubsub       *p2pPubSub.PubSub
	GrpcClient   *GRPCService_Client
	GrpcServer   *GRPCService_Server
	Blockchain   blockchainInf
	chunkManager *ChunkManager
	DHT          *kaddht.IpfsDHT
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

	chunkManager := NewChunkManager().SetEngine(chunk.NewSimpleChunk().MaxSize(50 * 1024))

	dht, err := kaddht.New(context.Background(), p2pHost)

	if err := dht.Bootstrap(context.Background()); err != nil {
		log.Println("failed to bootstrap DHT")
	}

	return &Node{
		Host:         p2pHost,
		Self:         selfPeer,
		Pubsub:       pubsub,
		GrpcClient:   grpcClient,
		GrpcServer:   grpcServer,
		Blockchain:   config.Blockchain,
		chunkManager: chunkManager,
		DHT:          dht,
	}
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

//func (node *Node) SubscribeChunk(topic string, handler func(*p2pPubSub.Subscription)) error {
//	subscription, err := node.Pubsub.Subscribe(topic)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	go handler(subscription)
//	return nil
//}
//
//func (node *Node) PublishChunk(topic string, data []byte) error {
//	chunkData, err := node.chunkManager.Split(data)
//	if err != nil {
//		return err
//	}
//	for i := range chunkData {
//		err = node.Pubsub.Publish(topic, chunkData[i])
//	}
//	return err
//}
