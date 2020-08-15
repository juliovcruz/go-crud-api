package main

import (
	"context"
	"io"
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

	//createUser("Victor", "victor@gmail.com", "12345", client)
	//deleteUser("5f37e4e4dbb6cea24257788e", client)
	//readUser("5f37e6f5dbb6cea24257788f", client)
	//updateUser("5f37e6f5dbb6cea24257788f", "Julios", "juliocruz.dev@gmail.com", "12345", client)
	listUser(client)

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

func readUser(id string, client userpb.UserServiceClient) {
	request := &userpb.ReadUserRequest{
		Id: id,
	}

	res, err := client.ReadUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Error in execution: %v", err)
	}

	log.Println(res)
}

func updateUser(id string, name string, email string, password string, client userpb.UserServiceClient) {
	request := &userpb.UpdateUserRequest{
		User: &userpb.User{
			Id:       id,
			Name:     name,
			Email:    email,
			Password: password,
		},
	}

	res, err := client.UpdateUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Error in execution: %v", err)
	}

	log.Println(res)
}

func listUser(client userpb.UserServiceClient) {
	request := &userpb.ListUserRequest{}

	stream, err := client.ListUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Eror in execution: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Eror in list: %v", err)
		}
		log.Println(res)
	}

}
