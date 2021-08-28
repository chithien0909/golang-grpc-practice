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
	// callPrimeNumberDecomposition(client)
	callAverage(client)
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

func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("Call Average ...")
	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("Call average err %v", err)
	}
	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num: 5,
		},
		calculatorpb.AverageRequest{
			Num: 10,
		},
		calculatorpb.AverageRequest{
			Num: 6,
		},
		calculatorpb.AverageRequest{
			Num: 7,
		},
		calculatorpb.AverageRequest{
			Num: 9,
		},
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Send request error %v", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Receive average response error %v", err)
	}
	log.Printf("Average response %v", resp)
}
