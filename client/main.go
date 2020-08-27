package main

import (
	"context"
	"fmt"
	"go-crud-api/userpb"
	"io"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in read env %v", err)
	}
	SERVER_PORT := os.Getenv("SERVER_PORT")
	connection, err := grpc.Dial("localhost:"+SERVER_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer connection.Close()

	client := userpb.NewUserServiceClient(connection)
	var id, name, email, password string

	for true {
		number := 7
		fmt.Println("0 - Exit\n1 - CreateUser\n2 - ReadUser\n3 - UpdateUser \n4 - DeleteUser\n5 - ListUser")
		fmt.Scanf("%d\n", &number)
		if number == 0 {
			break
		}
		if number == 1 {
			fmt.Println("Write name, email and password \n(name email password)")
			fmt.Scanf("%s %s %s\n", &name, &email, &password)
			createUser(name, email, password, client)
		} else if number == 2 {
			fmt.Println("Write id\n(id)")
			fmt.Scanf("%s\n", &id)
			readUser(id, client)
		} else if number == 3 {
			fmt.Println("Write id, name, email e password\n(id name email password)")
			fmt.Scanf("%s, %s, %s, %s, %s\n", &id, &name, &email, &password)
			updateUser(id, name, email, password, client)
		} else if number == 4 {
			fmt.Println("Write id\n(id)")
			fmt.Scanf("%s\n", &id)
			deleteUser(id, client)
		} else if number == 5 {
			listUser(client)
		}
		id = ""
		name = ""
		email = ""
		password = ""
		time.Sleep(2 * time.Second)
	}

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
