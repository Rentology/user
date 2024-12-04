package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	apiMiddlewares "user-service/internal/middleware"
	userHttp "user-service/internal/user/delivery/http"
	"user-service/internal/user/repository"
	"user-service/internal/user/service"
	"user-service/pkg/utils"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	userRepo := repository.NewUserRepository(s.db)

	userService := service.NewUserService(s.cfg, userRepo, s.log)

	userHandlers := userHttp.NewUserHandlers(s.cfg, userService, s.log)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, s.log, userService)

	allowedOrigins := "http://localhost:3000"
	if s.cfg.App.Env == "prod" {
		allowedOrigins = "http://localhost:80"
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{allowedOrigins},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true, // разрешает отправку учетных данных
	}))

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	userGroup := v1.Group("/users")

	userHttp.MapUserRoutes(userGroup, userHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.log.Info(fmt.Sprintf("Health check RequestID: %s", utils.GetRequestID(c)))
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return nil
}
