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

	//createUser("Julii", "juliovcruz0@gmail.com", "12345", client)
	deleteUser("5f37e4e4dbb6cea24257788e", client)

}

func createUser(name string, email string, password string, client userpb.UserServiceClient) {

	request := &userpb.CreateUserRequest{
		User: &userpb.User{
			Name:     name,
			Email:    email,
			Password: password,
		},
	}

	res, err := client.CreateUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Error in execution: %v", err)
	}

	log.Println(res)

}

func deleteUser(id string, client userpb.UserServiceClient) {
	request := &userpb.DeleteUserRequest{
		Id: id,
	}

	res, err := client.DeleteUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Error in execution: %v", err)
	}

	log.Println(res)
}
