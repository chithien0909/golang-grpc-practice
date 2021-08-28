package main

import (
	"context"
	"demo-grpc/calculator/calculatorpb"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum called")
	resp := &calculatorpb.SumResponse{
		Result: req.GetA() + req.GetB(),
	}

	return resp, nil
}
func (*server) PrimeNumberDecomposition(
	req *calculatorpb.PNDRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer,
) error {
	k := int32(2)
	n := req.GetNumber()
	for n > 1 {
		if n%k == 0 {
			n = n / k
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
			time.Sleep(500 * time.Millisecond)
			continue
		}
		k++
		fmt.Printf("k increase to %v\n", k)
	}
	return nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")

	if err != nil {
		log.Fatalf("Error while create listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Calculator is running ....")

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Error while server %v", err)
	}

}
