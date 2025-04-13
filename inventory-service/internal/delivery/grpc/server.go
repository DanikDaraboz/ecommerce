package grpc

import (
	"google.golang.org/grpc"
	"inventoryservice/pb"
)

type InventoryServer interface {
	pb.InventoryServiceServer
}

func RegisterInventoryServer(grpcServer *grpc.Server, handler InventoryServer) {
	pb.RegisterInventoryServiceServer(grpcServer, handler)
}