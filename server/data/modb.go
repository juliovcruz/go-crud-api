package data

import (
	"context"
	"go-crud-api/userpb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Client
var UserDb *mongo.Collection
var MongoContext context.Context

type UserServiceServer struct{}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
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

	result, err := UserDb.InsertOne(MongoContext, dataUser)
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

	result := UserDb.FindOne(ctx, bson.M{"_id": id})
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

	_, err = UserDb.DeleteOne(ctx, bson.M{"_id": id})
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

	result := UserDb.FindOneAndUpdate(ctx, bson.M{"_id": id},
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

func (s *UserServiceServer) ListUser(req *userpb.ListUserRequest, stream userpb.UserService_ListUserServer) error {
	data := &User{}

	cursor, err := UserDb.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return err
		}

		stream.Send(&userpb.ListUserResponse{
			User: &userpb.User{
				Id:       data.ID.Hex(),
				Name:     data.Name,
				Email:    data.Email,
				Password: data.Password,
			},
		})
		if err := cursor.Err(); err != nil {
			return err
		}

	}
	return nil

}
