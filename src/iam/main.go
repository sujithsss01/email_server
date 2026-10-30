package main

import (
	"context"

	"github.com/sujithsss01/email_server/iam/internal/postgres_client"
)

func main() {
	ctx := context.Background()
	sqlDB := postgres.NewConnection(ctx)
	defer sqlDB.Close()
	postgres.RunMigrations()
	gormDB := postgres.NewGORM(sqlDB)
	_ = gormDB
}
