package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	// healthpb "google.golang.org/grpc/health/grpc_health_v1"
	pstonks "stonks-service/proto"

	"stonks-service/stonks/config"
	"stonks-service/stonks/service"
)

func main() {
	log := service.GetLog(context.Background())
	log.Print("creating config...")
	err := config.CreateConfig()
	if err != nil {
		log.Fatalf("Unable to initialize service config: %v", err)
	}

	svc, err := service.NewStonksService(config.Config)
	if err != nil {
		log.Fatalf("Unable to start stonks-service: %v", err)
	}

	grpcErrChan := make(chan error, 0)
	log.Print("starting stonks-service...")
	RunServer(svc, config.Config.AppPort, grpcErrChan)

	select {
	case svcErr := <-grpcErrChan:
		log.Fatalf("error processing: %v", svcErr)
	}

	log.Println("doooooone.")
}

func RunServer(svc *service.StonksService, port string, errChan chan error) {

	go func(s *service.StonksService, p string) {
		url := fmt.Sprintf("0.0.0.0:%v", p)
		fmt.Println("listening on...", url)
		svcListen, err := net.Listen("tcp", url)
		if err != nil {
			err = fmt.Errorf("Unable to listen on: %v", svcListen)
			errChan <- err
		}
		var opts []grpc.ServerOption
		server := grpc.NewServer(opts...)

		pstonks.RegisterStonksServiceServer(server, s)
		err = server.Serve(svcListen)
		if err != nil {
			err = fmt.Errorf("Failed to run server: %v", err)
		}
		errChan <- err
	}(svc, port)

}
