package main

import (
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"net"
	"order-service/pkg/config"
	"order-service/pkg/db"
	"order-service/pkg/pb"
	"order-service/pkg/service"
)

func main() {
	cfg := config.LoadConfig()
	db := db.Init(cfg)

	lis, err := net.Listen("tcp", "localhost:"+cfg.TCP_PORT)
	if err != nil {
		panic(err)
	}

	os := service.NewOrderService(db, cfg)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, os)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
