package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sujithsss01/email_server/iam/internal/apperrors"
	"github.com/sujithsss01/email_server/iam/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
	username string,
	email string,
	phoneNumber string,
	password string,
) (*model.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, apperrors.ErrUUIDCreationFailed
	}
	user := &model.User{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		Username:    username,
		Email:       email,
		DateJoined:  time.Now(),
		PhoneNumber: phoneNumber,
	}
	authId, err := uuid.NewV7()
	if err != nil {
		return nil, apperrors.ErrUUIDCreationFailed
	}

	auth := &model.Auth{
		ID:       authId,
		UserId:   user.ID,
		AuthType: model.AuthPasswordHash,
		PassHash: &password,
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&auth).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id=?", userId).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
