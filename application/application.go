package application

import (
	"context"
	"sut-storage-go/lib/pkg/db"

	"google.golang.org/api/drive/v3"
	"google.golang.org/grpc"
)

type Application struct {
	DbClients    db.Handler
	GrpcServer   *grpc.Server
	GrpcClients  map[string]*grpc.ClientConn
	Context      context.Context
	DriveHandler *drive.Service
}
