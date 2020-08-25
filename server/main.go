package main

import (
	"context"
	"fmt"
	"go-crud-api/server/data"
	"go-crud-api/userpb"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Server is running in localhost:50050")

	listener, err := net.Listen("tcp", ":50050")

	if err != nil {
		fmt.Println("Server fails: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &data.UserServiceServer{})

	fmt.Println("Connecting to MongoDB...")
	data.MongoContext = context.Background()

	data.Db, err = mongo.Connect(data.MongoContext, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = data.Db.Ping(data.MongoContext, nil)
	if err != nil {
		log.Fatalf("Error in connection to MongoDb: %v\n", err)
	} else {
		fmt.Println("MongoDB is connected")
	}
	data.UserDb = data.Db.Database("learn-go-crud").Collection("user")

	fmt.Println("Server succesfully started on port: 50050")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	fmt.Println("\nStopping the server...")
	grpcServer.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	data.Db.Disconnect(data.MongoContext)
	fmt.Println("Done.")

}
