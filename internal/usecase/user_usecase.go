package usecase

import (
	"context"
	"errors"
	"go-auth-backend/internal/dto"
	"go-auth-backend/internal/entity"
	"go-auth-backend/internal/repository/mongodb"
	"go-auth-backend/pkg/jwt"
	"go-auth-backend/pkg/password"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type userUsecase struct {
	repo mongodb.UserRepository
	jwt  *jwt.JWT
}

func NewUserUsecase(repo mongodb.UserRepository, secret string, expiredHours int) UserUsecase {
	return &userUsecase{
		repo: repo,
		jwt:  jwt.New(secret, expiredHours),
	}
}

func (u *userUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// cek email sudah ada?
	if existing, _ := u.repo.FindByEmail(ctx, req.Email); existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, _ := password.Hash(req.Password)
	user := &entity.User{
		Email:    req.Email,
		Password: hashed,
		Name:     req.Name,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("email atau password salah")
	}

	if !password.Compare(user.Password, req.Password) {
		return nil, errors.New("email atau password salah")
	}

	token, err := u.jwt.Generate(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
