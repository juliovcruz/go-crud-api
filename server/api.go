package server

import (
	"context"
	"fmt"
	"go-crud-api/server/data"
	"go-crud-api/userpb"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

func Run() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in read env %v", err)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	SERVER_PORT := os.Getenv("SERVER_PORT")

	fmt.Println("Server is running in localhost:" + SERVER_PORT)

	listener, err := net.Listen("tcp", ":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Server fails: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &data.UserServiceServer{})

	fmt.Println("Connecting to MongoDB...")
	data.MongoContext = context.Background()

	data.Db, err = mongo.Connect(data.MongoContext,
		options.Client().ApplyURI("mongodb://"+DB_USER+":"+DB_PASS+"@"+DB_HOST+":"+DB_PORT))
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

	fmt.Println("Server successfully started on port:" + SERVER_PORT)

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
