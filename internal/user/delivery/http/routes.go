package http

import (
	"github.com/labstack/echo/v4"
)

type UserHandlers interface {
	CreateUser() echo.HandlerFunc
	GetUserById() echo.HandlerFunc
	GetUserByEmail() echo.HandlerFunc
}

func MapUserRoutes(userGroup *echo.Group, h UserHandlers) {
	userGroup.GET("/:id", h.GetUserById())
	userGroup.GET("", h.GetUserByEmail())
	userGroup.POST("", h.CreateUser())
}
