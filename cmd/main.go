package main

import (
	"log"
	"net"

	"grpc_with_some_methods/protos/go"
	"grpc_with_some_methods/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	dsn := "test_user:test_password@tcp(localhost:3306)/grpc_example?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&_go.Entity{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	entityService := &services.EntityHandler{
		UnimplementedEntityServiceServer: _go.UnimplementedEntityServiceServer{},
		DB:                               db,
	}
	_go.RegisterEntityServiceServer(s, entityService)

	log.Printf("Starting gRPC server on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
