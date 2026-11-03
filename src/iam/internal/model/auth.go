package model

import (
	"github.com/google/uuid"
)

type AuthType int16

const (
	AuthPasswordHash AuthType = iota + 1
)

type Auth struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId   uuid.UUID `gorm:"type:uuid; not null"`
	AuthType AuthType  `gorm:"type:smallint;not null"`
	PassHash *string

	User User `gorm:"foreignKey:UserId;references:ID"`
}

func (*Auth) TableName() string {
	return "auth"
}
