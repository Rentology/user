package middleware

import (
	"context"
	"log/slog"
	"user-service/internal/config"
	"user-service/internal/models"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type MiddlewareManager struct {
	cfg         *config.Config
	log         *slog.Logger
	userService UserService
}

func NewMiddlewareManager(cfg *config.Config, log *slog.Logger, service UserService) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, log: log, userService: service}
}
