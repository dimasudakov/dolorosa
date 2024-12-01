package main

import (
	"context"
	"dolorosa/internal/pipeline/operations/sbp"
	producer "dolorosa/internal/pkg/kafka"
	"dolorosa/internal/pkg/nirvana_helper"
	"dolorosa/internal/pkg/notifier"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"

	"dolorosa/internal/app/dolorosa"
	"dolorosa/internal/interceptors"
	"dolorosa/internal/pkg/logs_sender"
	"dolorosa/pkg/api/control"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	api "nirvana/pkg/api/nirvana"
)

const (
	nirvanaAddress = "localhost:50052"
)

func main() {
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	conn, err := grpc.NewClient(nirvanaAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	nirvanaClient := api.NewNirvanaClient(conn)

	kafkaProducer, err := producer.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create kafka producer: %v", err)
		return
	}

	logsSender := logs_sender.NewLogsSender(kafkaProducer)
	auditLogInterceptor := interceptors.NewAuditLogsInterceptor(logsSender)
	clientNotifier := notifier.NewNotifier()
	nirvanaHelper := nirvana_helper.NewNirvanaHelper(nirvanaClient)

	sbpChecker := sbp.NewCheckerSbp(sbp.CheckerFields{
		Notifier:         clientNotifier,
		ExceptionChecker: nirvanaHelper,
	})

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auditLogInterceptor.Unary()),
	)

	onlineControl := dolorosa.NewOnlineControlService(sbpChecker)

	control.RegisterOnlineControlServer(grpcServer, onlineControl)

	reflection.Register(grpcServer)

	log.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
