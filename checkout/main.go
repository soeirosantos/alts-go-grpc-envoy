package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/soeirosantos/alts/checkout/proto/checkout/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	port         = fmt.Sprintf(":%s", getEnv("PORT", "50054"))
	paymentAddr  = getEnv("PAYMENT_SVC_ADDR", "localhost:9902")
	shippingAddr = getEnv("SHIPPING_SVC_ADDR", "localhost:9903")
)

type checkoutService struct {
	pb.UnimplementedCheckoutServiceServer
}

func (s *checkoutService) PlaceOrder(ctx context.Context, in *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {

	// Dummy logic to simulate an order

	_, err := s.charge(ctx, &pb.Money{}, &pb.CreditCardInfo{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}

	shippingTrackingID, err := s.ship(ctx, &pb.Address{}, []*pb.CartItem{})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "shipping error: %+v", err)
	}

	return &pb.PlaceOrderResponse{Order: &pb.OrderResult{ShippingTrackingId: shippingTrackingID}}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	svc := &checkoutService{}
	pb.RegisterCheckoutServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *checkoutService) ship(ctx context.Context, address *pb.Address, items []*pb.CartItem) (string, error) {

	conn, err := grpc.DialContext(ctx, shippingAddr, grpc.WithInsecure())

	if err != nil {
		return "", fmt.Errorf("failed to connect shipment service: %+v", err)
	}

	defer conn.Close()

	resp, err := pb.NewShippingServiceClient(conn).ShipOrder(ctx, &pb.ShipOrderRequest{
		Address: address,
		Items:   items})
	if err != nil {
		return "", fmt.Errorf("shipment failed: %+v", err)
	}
	return resp.GetTrackingId(), nil
}

func (s *checkoutService) charge(ctx context.Context, amount *pb.Money, paymentInfo *pb.CreditCardInfo) (string, error) {

	// using Envoy
	// altsTC := alts.NewClientCreds(alts.DefaultClientOptions())
	// conn, err := grpc.DialContext(ctx, paymentAddr, grpc.WithTransportCredentials(altsTC))
	conn, err := grpc.Dial(paymentAddr, grpc.WithInsecure())

	if err != nil {
		return "", fmt.Errorf("failed to connect payment service: %+v", err)
	}

	defer conn.Close()

	resp, err := pb.NewPaymentServiceClient(conn).Charge(ctx, &pb.ChargeRequest{
		Amount:     amount,
		CreditCard: paymentInfo})
	if err != nil {
		return "", fmt.Errorf("could not charge the card: %+v", err)
	}

	return resp.GetTransactionId(), nil
}

func (s *checkoutService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *checkoutService) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
