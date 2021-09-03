gen-cal:
	protoc gatewaypb/gateway.proto --go_out=plugins=grpc:.

run-server:
	go run server/server.go server/models.go

run-client:
	go run client/client.go
