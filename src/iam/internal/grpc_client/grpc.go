package grpcclient

import (
	"context"
	"log"
	"net"

	"buf.build/go/protovalidate"
	"github.com/sujithsss01/email_server/config"
	"github.com/sujithsss01/email_server/iam/internal/pb"
	"github.com/sujithsss01/email_server/iam/internal/repository"
	"github.com/sujithsss01/email_server/iam/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func validateRequestServerInterceptor(validator protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		msg, ok := req.(proto.Message)
		if !ok {
			return nil, status.Error(
				codes.Internal,
				"invalid protobuf request",
			)
		}
		if err := validator.Validate(msg); err != nil {
			return nil, status.Error(
				codes.InvalidArgument,
				err.Error(),
			)
		}
		return handler(ctx, req)
	}
}

func registerServices(server *grpc.Server, db *gorm.DB) {
	repositories := repository.NewRepositories(db)
	iampb.RegisterUserServiceServer(
		server,
		service.NewUserService(repositories.UserRepository))
	iampb.RegisterAuthServiceServer(
		server,
		service.NewAuthService(repositories.AuthRepository, repositories.SessionRepository),
	)
}

func StartGrpcServer(gormDB *gorm.DB) {
	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	cf := config.GetConfig()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println("Listening failed")
		panic(err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(validateRequestServerInterceptor(validator)),
	)
	registerServices(server, gormDB)
	log.Println("started listening at %v:", lis.Addr())
	if cf.Environment == "local" {
		reflection.Register(server)
	}
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
