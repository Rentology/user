package http

import (
	"github.com/labstack/echo/v4"
	"user-service/internal/middleware"
)

type UserHandlers interface {
	CreateUser() echo.HandlerFunc
	GetUserById() echo.HandlerFunc
	GetUserByEmail() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
}

func MapUserRoutes(userGroup *echo.Group, h UserHandlers, mw *middleware.MiddlewareManager) {
	userGroup.GET("/:id", h.GetUserById(), mw.AuthJWTMiddleware())
	userGroup.GET("", h.GetUserByEmail(), mw.AuthJWTMiddleware())
	userGroup.POST("", h.CreateUser(), mw.AuthJWTMiddleware())
	userGroup.PATCH("", h.UpdateUser(), mw.AuthJWTMiddleware())
}
