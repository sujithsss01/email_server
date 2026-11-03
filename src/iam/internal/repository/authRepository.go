package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sujithsss01/email_server/iam/internal/apperrors"
	"github.com/sujithsss01/email_server/iam/internal/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

type GetUserPassHashResponse struct {
	UserId   uuid.UUID `gorm:"column:id"`
	PassHash *string   `gorm:"column:pass_hash"`
}

func (r *AuthRepository) GetUserPassHash(
	ctx context.Context,
	username string,
) (*GetUserPassHashResponse, error) {
	res := &GetUserPassHashResponse{}
	err := r.db.WithContext(ctx).
		Table("auth a").
		Select("u.id, a.pass_hash").
		Joins("JOIN users u ON u.id = a.user_id").
		Where("a.auth_type=? AND ( u.username=? OR u.email=? OR u.phone_number=? )", model.AuthPasswordHash, username, username, username).
		First(res).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if res.PassHash == nil {
		return nil, apperrors.ErrPasswordNotFound
	}
	return res, nil
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
