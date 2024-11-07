package middleware

import (
	"log/slog"
	"user-service/internal/config"
)

type MiddlewareManager struct {
	cfg *config.Config
	log *slog.Logger
}

func newMiddlewareManager() *MiddlewareManager {
	return nil
}
