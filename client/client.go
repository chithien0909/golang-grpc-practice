package main

import (
	"context"
	"demo-grpc/calculator/calculatorpb"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error while dial %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	log.Printf("Service client %v", client)
	// callSum(client)
	callPrimeNumberDecomposition(client)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		A: 4,
		B: 6,
	})
	if err != nil {
		log.Fatalf("Call sum api error %v", err)
	}

	fmt.Printf("Sum result: %v", resp.Result)
}

func callPrimeNumberDecomposition(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})
	if err == io.EOF {
		log.Println("server finished streaming")
	}
	if err != nil {
		log.Fatalf("callPND err %v", err)
	}
	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finished streaming")
			return
		}
		log.Printf("prime number %v", resp.GetResult())
	}
}
