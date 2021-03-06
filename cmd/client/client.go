package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/JoaoVictorPereiraTeixeira/GO-gRPC/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: $v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "joao",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request", err)
	}
	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive de , %v", err)
		}

		fmt.Println("Status:", stream.Status, "-", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Joao",
			Email: "joao@joao.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Joao2",
			Email: "joao2@joao.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Joao3",
			Email: "joao3@joao.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Joao4",
			Email: "joao4@joao.com",
		},
		&pb.User{
			Id:    "w5",
			Name:  "Joao5",
			Email: "joao5@joao.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		fmt.Println(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Joao",
			Email: "joao@joao.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Joao2",
			Email: "joao2@joao.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Joao3",
			Email: "joao3@joao.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Joao4",
			Email: "joao4@joao.com",
		},
		&pb.User{
			Id:    "w5",
			Name:  "Joao5",
			Email: "joao5@joao.com",
		},
	}

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	wait := make(chan int)

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
