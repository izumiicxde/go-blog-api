package types

import (
	"gorm.io/gorm"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(u RegisterUserPayload) error
}

type User struct {
	gorm.Model
	ID        int    `json:"id,omitempty"`
	FistName  string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=30"`
	LastName  string `json:"last_name" validate:"required,max=30"`
	Email     string `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	AvatarUrl string `json:"avatar_url" validate:"required"`
}
