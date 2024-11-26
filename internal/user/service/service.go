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

	formattedDate, err := utils.ParseDate(user.BirthDate)

	if err == nil {
		user.BirthDate = &formattedDate
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	currentUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if user.Email != "" {
		currentUser.Email = user.Email
	}
	if user.Phone != nil {
		currentUser.Phone = user.Phone
	}
	if user.Name != nil {
		currentUser.Name = user.Name
	}
	if user.LastName != nil {
		currentUser.LastName = user.LastName
	}
	if user.SecondName != nil {
		currentUser.SecondName = user.SecondName
	}
	if user.BirthDate != nil {
		currentUser.BirthDate = user.BirthDate
	}
	if user.Sex != nil {
		currentUser.Sex = user.Sex
	}
	currentUser, err = s.userRepo.Update(ctx, currentUser)
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	formattedDate, err := utils.ParseDate(user.BirthDate)
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
	formattedDate, err := utils.ParseDate(user.BirthDate)
	if err == nil {
		user.BirthDate = &formattedDate
	}
	return user, nil
}
