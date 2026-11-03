package repository

import "gorm.io/gorm"

type Repositories struct {
	UserRepository    *UserRepository
	AuthRepository    *AuthRepository
	SessionRepository *SessionRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository:    NewUserRepository(db),
		AuthRepository:    NewAuthRepository(db),
		SessionRepository: NewSessionRepository(db),
	}
}
