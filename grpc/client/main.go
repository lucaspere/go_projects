package main

import (
	"context"
	"fmt"
	"log"
	"os"

	users "github.com/lucaspere/grpc/service"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func getUser(client users.UsersClient, u *users.UserGetRequest) (*users.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func createUserRequest(jsonQuery string) (*users.UserGetRequest, error) {
	u := users.UserGetRequest{}
	input := []byte(jsonQuery)
	fmt.Println(u)
	return &u, protojson.Unmarshal(input, &u)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal(
			"Must specify a gRPC server address and search query",
		)
	}
	fmt.Println(os.Args)
	conn, err := setupGrpcConn(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	u, err := createUserRequest(os.Args[2])
	if err != nil {
		log.Fatalf("Bad user input: %v", err)
	}
	fmt.Println(u)
	result, err := getUser(c, u)
	if err != nil {
		log.Fatal(err)
	}

	data, err := protojson.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(
		os.Stdout, string(data),
	)
}
