package service

import (
	"context"

	"github.com/google/uuid"
	iampb "github.com/sujithsss01/email_server/iam/internal/pb"
	"github.com/sujithsss01/email_server/iam/internal/repository"
	"github.com/sujithsss01/email_server/iam/internal/utils"
)

type UserService struct {
	iampb.UnimplementedUserServiceServer
	repository *repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, request *iampb.CreateUserRequest) (*iampb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}
	user, err := s.repository.CreateUser(ctx, request.FirstName, request.LastName, request.Username, request.Email, request.PhoneNumber, hashedPassword)
	if err != nil {
		return nil, err
	}
	return &iampb.CreateUserResponse{
		Status: "Created",
		UserId: user.ID.String(),
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, request *iampb.GetUserRequest) (*iampb.GetUserResponse, error) {
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, err
	}
	user, err := s.repository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &iampb.GetUserResponse{
		UserId:      user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}
