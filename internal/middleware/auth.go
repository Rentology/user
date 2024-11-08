package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"user-service/internal/config"
	"user-service/lib/sl"
	"user-service/pkg/httpErrors"
	"user-service/pkg/utils"
)

func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookieToken, err := c.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.ErrNoCookie))
				}
				return c.JSON(http.StatusInternalServerError, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			token := cookieToken.Value

			sl.Infof(mw.log, "auth middleware token: %s", token)
			if token != "" {

				if err := mw.validateJWTToken(token, mw.userService, c, mw.cfg); err != nil {
					mw.log.Error("middleware validateJWTToken", "header JWT", err.Error())
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				return next(c)
			}

			cookie, err := c.Cookie("token")
			if err != nil {
				mw.log.Error("c.Cookie", "error", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			if err = mw.validateJWTToken(cookie.Value, mw.userService, c, mw.cfg); err != nil {
				mw.log.Error("validateJWTToken", "error", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			return next(c)
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, userService UserService, c echo.Context, cfg *config.Config) error {
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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userUID, ok := claims["uid"].(float64)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userId := int64(userUID)
		mw.log.Info("uid", "uid", userId)
		user, err := userService.GetByID(c.Request().Context(), userId)
		if err != nil {
			return err
		}
		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil

}
