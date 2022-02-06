package main

import (
	"context"
	"fmt"
	"github.com/snxl/GO-gRPC/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:5001", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	AddUser(client)
	AddUserVerbose(client)
	AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Test",
		Email: "test@test.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	fmt.Println("running")
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}

}

func AddUsers(client pb.UserServiceClient) {

	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "test 1",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "test 2",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "test 3",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "test 4",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "5",
			Name:  "test 5",
			Email: "test@test.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receave response %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request %v")
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "test 1",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "test 2",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "test 3",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "test 4",
			Email: "test@test.com",
		},
		&pb.User{
			Id:    "5",
			Name:  "test 5",
			Email: "test@test.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
