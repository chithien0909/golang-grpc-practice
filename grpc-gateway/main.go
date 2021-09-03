

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-gateway/gatewaypb"
	"log"
	"net"
)

type server struct {
}

func (server) Echo(ctx context.Context, msg *gatewaypb.StringMessage) (*gatewaypb.StringMessage, error)  {
	log.Printf("Receive message %s\n", msg.GetMessage())
	return msg, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50080")

	s := grpc.NewServer()

	gatewaypb.RegisterGatewayServiceServer(s, &server{})

	fmt.Println("Gateway servioce is running ....")

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Error while server %v", err)
	}

}
