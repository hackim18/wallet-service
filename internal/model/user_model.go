package model

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Email       string     `json:"email,omitempty"`
	AccessToken string     `json:"access_token,omitempty"`
}

type UserLogin struct {
	AccessToken string `json:"access_token"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `json:"-" validate:"required,max=100"`
	Password string    `json:"password,omitempty" validate:"max=100"`
	Name     string    `json:"name,omitempty" validate:"max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required,max=100"`
}
