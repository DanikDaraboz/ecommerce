package grpc

import (
	"google.golang.org/grpc"
	"github.com/danikdaraboz/ecommerce/order-service/pb"
)

type OrderServer interface {
	pb.OrderServiceServer
}

func RegisterOrderServer(grpcServer *grpc.Server, handler OrderServer) {
	pb.RegisterOrderServiceServer(grpcServer, handler)
}