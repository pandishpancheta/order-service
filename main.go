package main

import (
	"net"
	"order-service/pkg/config"
	"order-service/pkg/db"
	"order-service/pkg/pb"
	"order-service/pkg/service"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	pgdb := db.Init(cfg)
	db.InitTable(pgdb)

	lis, err := net.Listen("tcp", "localhost:"+cfg.TCP_PORT)
	if err != nil {
		panic(err)
	}

	os := service.NewOrderService(pgdb, cfg)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, os)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
