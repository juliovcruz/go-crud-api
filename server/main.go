package main

import (
	"context"
	"fmt"
	"learn-go-crud/userpb"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var db *mongo.Client
var userdb *mongo.Collection
var mongoContext context.Context

type UserServiceServer struct{}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"email"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	userRequest := req.GetUser()

	dataUser := User{
		Name:     userRequest.GetName(),
		Email:    userRequest.GetEmail(),
		Password: userRequest.GetPassword(),
	}

	result, err := userdb.InsertOne(mongoContext, dataUser)
	if err != nil {
		return nil, err
	}

	idResult := result.InsertedID.(primitive.ObjectID)
	userRequest.Id = idResult.Hex()

	return &userpb.CreateUserResponse{User: userRequest}, nil

}

func main() {
	fmt.Println("Server is running in localhost:5050")

	listener, err := net.Listen("tcp", ":5050")

	if err != nil {
		fmt.Println("Server fails: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &UserServiceServer{})

	fmt.Println("Connecting to MongoDB...")
	mongoContext = context.Background()

	db, err = mongo.Connect(mongoContext, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(mongoContext, nil)
	if err != nil {
		log.Fatalf("Error in connection to MongoDb: %v\n", err)
	} else {
		fmt.Println("MongoDB is connected")
	}

	userdb = db.Database("learn-go-crud").Collection("user")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	fmt.Println("Server succesfully started on port: 5050")

	fmt.Println("\nStopping the server...")
	grpcServer.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoContext)
	fmt.Println("Done.")

}