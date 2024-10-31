package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"

	pb "user-service/userpb"

	"google.golang.org/grpc"
)

var users = make(map[string]*pb.UserResponse) // In-memory storage
var idCounter = 1

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	id := strconv.Itoa(idCounter)
	idCounter++
	user := &pb.UserResponse{Id: id, Name: req.Name, Email: req.Email}
	users[id] = user
	return user, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, exists := users[req.Id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, exists := users[req.Id]
	if !exists {
		return nil, errors.New("user not found")
	}
	user.Name = req.Name
	user.Email = req.Email
	return user, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, exists := users[req.Id]
	if !exists {
		return &pb.DeleteUserResponse{Success: false}, errors.New("user not found")
	}
	delete(users, req.Id)
	return &pb.DeleteUserResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Println("User Service is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
