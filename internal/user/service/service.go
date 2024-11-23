package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"user-service/internal/config"
	"user-service/internal/models"
	handlers "user-service/internal/user/delivery/http"
	"user-service/pkg/utils"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userService struct {
	cfg      *config.Config
	userRepo UserRepository
	log      *slog.Logger
}

func NewUserService(cfg *config.Config, userRepo UserRepository, log *slog.Logger) handlers.UserService {
	return &userService{cfg, userRepo, log}
}

func (s *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := s.GetByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return nil, err
		}
	}
	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	formattedDate, err := utils.ParseDate(*user.BirthDate)

	if err == nil {
		user.BirthDate = &formattedDate
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	panic("implement me")
}

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	formattedDate, err := utils.ParseDate(*user.BirthDate)
	if err == nil {
		user.BirthDate = &formattedDate
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	formattedDate, err := utils.ParseDate(*user.BirthDate)
	if err == nil {
		user.BirthDate = &formattedDate
	}
	return user, nil
}
