package main

import (
	"book-service_gc2p3/config"
	"book-service_gc2p3/pb"
	"book-service_gc2p3/service"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	config.InitViper()
	config.InitMongo()

	grpcServer := grpc.NewServer()

	bookServer := service.NewBookService(config.DB)
	pb.RegisterBookServiceServer(grpcServer, bookServer)

	gRPCPort := config.Viper.GetString("GRPC_PORT")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on port :", gRPCPort)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}

}
