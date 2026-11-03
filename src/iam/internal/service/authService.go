package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sujithsss01/email_server/iam/internal/apperrors"
	iampb "github.com/sujithsss01/email_server/iam/internal/pb"
	"github.com/sujithsss01/email_server/iam/internal/repository"
	"github.com/sujithsss01/email_server/iam/internal/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	iampb.UnimplementedAuthServiceServer
	authRepository    *repository.AuthRepository
	sessionRepository *repository.SessionRepository
}

func (s *AuthService) GetRefreshToken(ctx context.Context, request *iampb.GetRefreshTokenRequest) (*iampb.GetRefreshTokenResponse, error) {
	res, err := s.authRepository.GetUserPassHash(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	err = utils.VerifyPasswordAndHash(request.Password, *res.PassHash)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	hashedRefreshToken := utils.HashRefreshToken(refreshToken)
	sessionId, err := s.sessionRepository.StoreRefreshToken(ctx, res.UserId, nil, hashedRefreshToken)
	return &iampb.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
		SessionId:    sessionId,
	}, nil
}

func (s *AuthService) RotateRefreshToken(ctx context.Context, request *iampb.RotateRefreshTokenRequest) (*iampb.RotateRefreshTokenResponse, error) {
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	sessionId, err := uuid.Parse(request.SessionId)
	if err != nil {
		return nil, err
	}
	hashedRefreshToken := utils.HashRefreshToken(newRefreshToken)
	err = s.sessionRepository.RotateRefreshToken(ctx, sessionId, request.RefreshToken, hashedRefreshToken)
	if err != nil {
		return nil, err
	}
	return &iampb.RotateRefreshTokenResponse{
		RefreshToken: newRefreshToken,
		SessionId:    request.SessionId,
	}, nil
}

func (s *AuthService) GetAuthToken(ctx context.Context, request *iampb.GetAuthTokenRequest) (*iampb.GetAuthTokenResponse, error) {
	hashedRefreshToken := utils.HashRefreshToken(request.RefreshToken)
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, apperrors.ErrInvalidID
	}
	session, err := s.sessionRepository.FetchSessionDetailsByToken(ctx, userId, hashedRefreshToken)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateJWT(request.UserId, session.SessionId.String())
	if err != nil {
		return nil, err
	}
	return &iampb.GetAuthTokenResponse{
		AuthToken: token,
	}, nil
}

func (s *AuthService) LogoutSession(ctx context.Context, request *iampb.LogoutSessionRequest) (*emptypb.Empty, error) {
	claim, err := utils.ValidateJwt(request.JwtToken)
	if err != nil {
		return nil, err
	}
	sessionId, err := uuid.Parse(claim.SessionId)
	if err != nil {
		return nil, apperrors.ErrInvalidID
	}
	err = s.sessionRepository.RevokeSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) LogoutAllSession(ctx context.Context, request *iampb.LogoutAllSessionRequest) (*emptypb.Empty, error) {
	claim, err := utils.ValidateJwt(request.JwtToken)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return nil, apperrors.ErrInvalidID
	}
	err = s.sessionRepository.RevokeTokensByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func NewAuthService(authRepository *repository.AuthRepository, sessionRepository *repository.SessionRepository) *AuthService {
	return &AuthService{
		authRepository:    authRepository,
		sessionRepository: sessionRepository,
	}
}
