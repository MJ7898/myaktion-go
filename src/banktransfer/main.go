package main

import (
	"fmt"
	"github.com/MJ7898/myaktion-go/src/banktransfer/grpc/banktransfer"
	"github.com/MJ7898/myaktion-go/src/banktransfer/kafka"
	"github.com/MJ7898/myaktion-go/src/banktransfer/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

var grpcPort = 9111

func init() {
	defer kafka.EnsureTransactionTopic()
	// init logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Info("Log level not specified, set default to: INFO")
		log.SetLevel(log.InfoLevel)
		return
	}
	log.SetLevel(level)
}

func main() {
	log.Info("Starting Banktransfer server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on grpc port %d: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer()
	transferService := service.NewBankTransferService()
	transferService.Start()
	defer transferService.Stop()
	banktransfer.RegisterBankTransferServer(grpcServer, transferService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}