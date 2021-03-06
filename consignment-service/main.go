// consignement-service/main.go
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/midnightrun/shippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository holds all consignments
type Repository struct {
	consignement []*pb.Consignment
}

// Create creates a new consignment and adds it to the repo
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignement, consignment)
	repo.consignement = updated
	return consignment, nil
}

// GetAll returns all current consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignement
}

// Service satisfy the proto interface
type Service struct {
	repo IRepository
}

// CreateConsignment creates a new consignment
func (s *Service) CreateConsignment(ctx context.Context, request *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(request)

	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

// GetConsignment gets all consignments
func (s *Service) GetConsignment(ctx context.Context, request *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()

	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &Service{repo})

	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("err: %v", err)
	}
}
