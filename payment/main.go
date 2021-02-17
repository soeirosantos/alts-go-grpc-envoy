package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	"github.com/soeirosantos/alts/payment/proto/checkout/pb"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/credentials/alts"
)

var (
	port   = fmt.Sprintf(":%s", getEnv("PORT", "50052"))
	altsSA = getEnv("ALTS_SERVICE_ACCOUNT", "")
)

type paymentService struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *paymentService) Charge(ctx context.Context, in *pb.ChargeRequest) (*pb.ChargeResponse, error) {

	// NOOP Stub

	txID, _ := uuid.NewUUID()

	log.Printf("payment processed (transaction_id: %s)", txID.String())

	return &pb.ChargeResponse{TransactionId: txID.String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Enable ALTS
	altsTC := alts.NewServerCreds(alts.DefaultServerOptions())
	srv := grpc.NewServer(grpc.Creds(altsTC))
	// srv := grpc.NewServer()

	svc := &paymentService{}

	pb.RegisterPaymentServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (s *paymentService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *paymentService) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
