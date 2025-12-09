// internal/usecase/user_usecase.go
package usecase

import (
	"auth-service/internal/entity"
	"auth-service/pkg/jwt"
	"auth-service/pkg/password"
	"context"
	"errors"
	"go-auth-backend/internal/config"
	"time"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type userUsecase struct {
	userRepo  repository.UserRepository
	jwtHelper *jwt.JWT
	config    *config.Config
}

func NewUserUsecase(userRepo repository.UserRepository, cfg *config.Config) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtHelper: jwt.NewJWT(cfg.JWTSecret, time.Duration(cfg.JWTExpiredHours)*time.Hour),
		config:    cfg,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// cek email sudah ada?
	exists, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if exists != nil {
		return nil, errors.New("email already registered")
	}

	hashed, _ := password.Hash(req.Password)
	user := &entity.User{
		Email:    req.Email,
		Password: hashed,
		Name:     req.Name,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID.Hex(),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !password.Compare(user.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := u.jwtHelper.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		UserResponse: dto.UserResponse{
			ID:    user.ID.Hex(),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
