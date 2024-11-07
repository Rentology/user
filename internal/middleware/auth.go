package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"user-service/internal/config"
	us "user-service/internal/user/delivery/http"
	"user-service/lib/sl"
	"user-service/pkg/httpErrors"
)

func (mw *MiddlewareManager) AuthJWTMiddleware(userService us.UserService, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")
			sl.Infof(mw.log, "auth middleware bearerHeader: %s", bearerHeader)
			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					mw.log.Error("auth middleware", "headerParts", "len(headerParts) != 2")
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				//tokenString := headerParts[1]

				return next(c)
			}
			return next(c)
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, userService us.UserService, cfg *config.Config) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	return nil

}
