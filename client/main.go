package main

import (
	"context"
	"log"
	"time"

	"github.com/soeirosantos/alts/client/proto/checkout/pb"
	"google.golang.org/grpc"
)

const (
	checkoutAddr = "localhost:50054"
)

func main() {
	conn, err := grpc.Dial(checkoutAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCheckoutServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.PlaceOrder(ctx, &pb.PlaceOrderRequest{})
	if err != nil {
		log.Fatalf("could not place order: %v", err)
	}
	log.Printf("Order placed (tracking_id: %s)", r.GetOrder().GetShippingTrackingId())
}
