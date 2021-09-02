package main

import (
	"context"
	"demo-grpc/calculator/calculatorpb"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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
func (*server) SumWithDeadline(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum with deadline called")
	for i := 0; i < 3 ; i++ {
		if ctx.Err() == context.Canceled {
			log.Printf("Sum with deadline canceled")
			return nil, status.Errorf(codes.Canceled, "client canceled sum  with deadline req")
		}
		time.Sleep(1 * time.Second)
	}
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

func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("FindMax called...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("FindMax EOF Server")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while Recv FindMax %v", err)
		}
		num := req.GetNum()
		if num > max {
			max = num
		}
		resp := &calculatorpb.FindMaxResponse{
			Result: num,
		}
		err = stream.Send(resp)
		if err != nil {
			log.Fatalf("Send max err %v", err)
			return err
		}
	}
}

func (*server) Square(ctx context.Context,req *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error){
	log.Println("Square called...")
	num := req.GetNum()
	if num < 0 {
		log.Printf("Expect num > 0, req num was %v", num)
		return nil, status.Errorf(codes.InvalidArgument,"Expect num > 0, req num was %v", num)
	}
	return &calculatorpb.SquareResponse{
		SquareRoot: math.Sqrt(float64(req.GetNum())),
	}, nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")

	if err != nil {
		log.Fatalf("Error while create listen %v", err)
	}
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.key"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Create creds ssl error %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)

	s := grpc.NewServer(opts)

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Calculator is running ....")

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Error while server %v", err)
	}

}
