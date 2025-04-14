package grpc

import (
	"google.golang.org/grpc"
	"github.com/danikdaraboz/ecommerce/inventory-service/pb"
)

type InventoryServer interface {
	pb.InventoryServiceServer
}

func RegisterInventoryServer(grpcServer *grpc.Server, handler InventoryServer) {
	pb.RegisterInventoryServiceServer(grpcServer, handler)
}