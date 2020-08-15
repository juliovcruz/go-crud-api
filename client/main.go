package main

import (
	"context"
	"learn-go-crud/userpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer connection.Close()

	client := userpb.NewUserServiceClient(connection)

	createUser("Julii", "juliovcruz0@gmail.com", "12345", client)

}

func createUser(name string, email string, password string, c userpb.UserServiceClient) {

	request := &userpb.CreateUserRequest{
		User: &userpb.User{
			Name:     name,
			Email:    email,
			Password: password,
		},
	}

	res, err := c.CreateUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Error in execution: %v", err)
	}

	log.Println(res)

}
