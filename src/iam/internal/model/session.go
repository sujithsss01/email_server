package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	UserId       uuid.UUID `gorm:"not null"`
	RefreshToken string    `gorm:"type:varchar(128);not null;unique"`
	ExpiresAt    time.Time `gorm:"not null"`
	RevokedAt    *time.Time
	CreatedAt    time.Time `gorm:"not null"`
	SessionId    uuid.UUID `gorm:"not null"`

	User User `gorm:"foreignKey:UserId;references:ID"`
}

func (*Session) TableName() string {
	return "session"
}
