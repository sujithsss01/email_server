package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	FirstName   string    `gorm:"not null"`
	LastName    string    `gorm:"not null"`
	Username    string    `gorm:"not null:unique"`
	Email       string    `gorm:"type:citext;unique;not null"`
	DateJoined  time.Time `gorm:"not null"`
	PhoneNumber string    `gorm:"not null"`
}
