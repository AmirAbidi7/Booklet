package dtos

import "backend/internal/models"

type RegisterRequest struct {
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=8"`
	Role     models.Role `json:"role" validate:"lt=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    uint        `json:"id"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}
