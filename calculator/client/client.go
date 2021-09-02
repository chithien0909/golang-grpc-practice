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
	"time"

	"google.golang.org/grpc"
)

func main() {
	certFile := "ssl/ca.crt"
	keyFile := "ssl/ca.key"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Create creds ssl error %v\n", sslErr)
		return
	}

	cc, err := grpc.Dial("localhost:50069", grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("Error while dial %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	log.Printf("Service client %v", client)
	callSum(client)
	// callPrimeNumberDecomposition(client)
	//callAverage(client)
	//callFindMax(client)
	//callSquareRoot(client, -234)
	//callSumWithDeadline(client, 4 * time.Second)
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
func callSumWithDeadline(c calculatorpb.CalculatorServiceClient, timeout time.Duration) {
	log.Println("Call sum with deadline ...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()


	resp, err := c.SumWithDeadline(ctx, &calculatorpb.SumRequest{
		A: 4,
		B: 6,
	})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("Call sum with deadline  DeadlineExceeded")
				return
			} else {
				log.Printf("Call sum with deadline  api err %v", err)
				return
			}
		} else {
			log.Printf("Call sum with deadline unknown error %v", err)
			return
		}
	}
	log.Printf("Sum with deadline result: %v", resp.Result)
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

func callFindMax(c calculatorpb.CalculatorServiceClient)  {

	log.Println("Call Find Max ...")
	stream, err := c.FindMax(context.Background())
	waitc := make(chan struct{})
	go func() {
		listReq := []calculatorpb.FindMaxRequest{
			calculatorpb.FindMaxRequest{
				Num: 5,
			},
			calculatorpb.FindMaxRequest{
				Num: 10,
			},
			calculatorpb.FindMaxRequest{
				Num: 6,
			},
			calculatorpb.FindMaxRequest{
				Num: 7,
			},
			calculatorpb.FindMaxRequest{
				Num: 9,
			},
		}
		for _, v := range listReq {
			_ = stream.Send(&v)
			if err != nil {
				log.Fatalf("Send find max request err %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		_ = stream.CloseSend()
	}()
	go func() {
		for {
			resv, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Ending find max api ...")
				break
			}
			if err != nil {

				log.Fatalf("Receive find max response error %v", err)
			}
			fmt.Printf("Max value: %v\n", resv.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func callSquareRoot(c calculatorpb.CalculatorServiceClient, num int32)  {
	log.Println("Call Square Root ...")
	resp, err := c.Square(context.Background(), &calculatorpb.SquareRequest{
		Num: num,
	})
	if err != nil {
		if errStatus, ok := status.FromError(err); ok{
			log.Printf("Error msg: %v\n", errStatus.Message())
			log.Printf("Error code: %v\n", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("InvalidArgument num: %v\n", num)
				return
			}
		}
		log.Fatalf("call square root api err %v\n", err)
	}
	log.Printf("Square root response %v\n", resp.GetSquareRoot())
}