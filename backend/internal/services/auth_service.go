package services

import (
	"context"
	"errors"

	"backend/internal/auth"
	"backend/internal/dtos"
	"backend/internal/models"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(req *dtos.RegisterRequest) (*dtos.AuthResponse, error) {
	ctx := context.Background()

	_, err := gorm.G[models.User](s.DB).Where("email = ?", req.Email).First(ctx)

	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if user.Role == 0 {
		user.Role = models.Regular
	}

	if err = gorm.G[models.User](s.DB).Create(ctx, &user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) Login(req *dtos.LoginRequest) (*dtos.AuthResponse, error) {
	ctx := context.Background()
	user, err := gorm.G[models.User](s.DB).Where("email = ?", req.Email).First(ctx)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err = auth.ComparePasswords(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*dtos.UserResponse, error) {
	ctx := context.Background()
	user, err := gorm.G[models.User](s.DB).Where("id = ?", userID).First(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dtos.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
