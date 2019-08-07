package p2p

import (
	"context"
	"google.golang.org/grpc"
)

type blockchainInf interface {
}

type GRPCService_Server struct {
	Blockchain blockchainInf
}

func (self *GRPCService_Server) registerServices(grpsServer *grpc.Server) {
	RegisterGreetServiceServer(grpsServer, self)
}

func (self *GRPCService_Server) Greet(ctx context.Context, req *GreetRequest) (*GreetResponse, error) {
	//fmt.Println("invoke with Greet")
	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	res := &GreetResponse{
		Result: result,
	}
	return res, nil
}
