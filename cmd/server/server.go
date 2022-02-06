package main

import (
	"fmt"
	"github.com/snxl/GO-gRPC/pb"
	"github.com/snxl/GO-gRPC/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", "localhost:5001")

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	fmt.Println("server started")

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not serve %v", err)
	}

}
