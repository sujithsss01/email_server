package main

import (
	"context"
	"github.com/sujithsss01/email_server/iam/internal/grpc_client"
	"github.com/sujithsss01/email_server/iam/internal/postgres_client"
)

func main() {
	ctx := context.Background()
	sqlDB := postgresclient.NewConnection(ctx)
	defer sqlDB.Close()
	postgresclient.RunMigrations()
	gormDB := postgresclient.NewGORM(sqlDB)
	grpcclient.StartGrpcServer(gormDB)
	_ = gormDB
}
