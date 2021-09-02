package main

import (
	"contact/contactpb"
	"google.golang.org/grpc"
	"log"
)

func main() {

	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error while dial %v", err)
	}
	defer cc.Close()

	client := contactpb.NewContactServiceClient(cc)

	log.Printf("Service client %v", client)
}
