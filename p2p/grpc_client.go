package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	p2pgrpc "github.com/paralin/go-libp2p-grpc"
	"google.golang.org/grpc"
	"log"
)

type GRPCService_Client struct {
	p2pgrpc *p2pgrpc.GRPCProtocol
}

func (self *GRPCService_Client) Greet(ctx context.Context, peerID peer.ID, first string) (string, error) {

	grpcConn, err := self.p2pgrpc.Dial(ctx, peerID, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", err
	}

	greetClient := NewGreetServiceClient(grpcConn)

	greetRely, err := greetClient.Greet(ctx, &GreetRequest{Greeting: &Greeting{FirstName: first}})
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return greetRely.GetResult(), nil
}
