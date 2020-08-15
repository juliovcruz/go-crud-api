package main

import (
	"context"
	"fmt"
	"learn-go-crud/userpb"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var db *mongo.Client
var userDb *mongo.Collection
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

	result, err := userDb.InsertOne(mongoContext, dataUser)
	if err != nil {
		return nil, err
	}

	idResult := result.InsertedID.(primitive.ObjectID)
	userRequest.Id = idResult.Hex()

	return &userpb.CreateUserResponse{User: userRequest}, nil

}

func (s *UserServiceServer) ReadUser(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	result := userDb.FindOne(ctx, bson.M{"_id": id})
	data := User{}
	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	response := &userpb.ReadUserResponse{
		User: &userpb.User{
			Id:       id.Hex(),
			Name:     data.Name,
			Email:    data.Email,
			Password: data.Password,
		},
	}

	return response, nil

}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	_, err = userDb.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil

}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	user := req.GetUser()

	id, err := primitive.ObjectIDFromHex(user.GetId())
	if err != nil {
		return nil, err
	}

	data := bson.M{
		"name":     user.GetName(),
		"email":    user.GetEmail(),
		"password": user.GetPassword(),
	}

	result := userDb.FindOneAndUpdate(ctx, bson.M{"_id": id},
		bson.M{"$set": data}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := User{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponse{
		User: &userpb.User{
			Id:       decoded.ID.Hex(),
			Name:     decoded.Name,
			Email:    decoded.Email,
			Password: decoded.Password,
		},
	}, nil

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

	userDb = db.Database("learn-go-crud").Collection("user")

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
