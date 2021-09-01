package main

import (
	"context"
	"demo-grpc/calculator/calculatorpb"
	"fmt"
	"io"
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
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		k++
		fmt.Printf("k increase to %v\n", k)
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called...")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("Error while Recv Average %v", err)
		}
		log.Printf("Receive req %v\n", req)
		total += req.GetNum()
		count++
	}
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
