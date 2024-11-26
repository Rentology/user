package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
	"user-service/internal/config"
	"user-service/internal/models"
	"user-service/pkg/httpErrors"
	"user-service/pkg/utils"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userHandlers struct {
	cfg         *config.Config
	userService UserService
	log         *slog.Logger
}

func NewUserHandlers(cfg *config.Config, userService UserService, log *slog.Logger) UserHandlers {
	return &userHandlers{cfg, userService, log}
}

func (h *userHandlers) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling Create", slog.String("request_id", requestID))
		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		user, err := h.userService.Create(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, user)
	}
}

func (h *userHandlers) GetUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling GetUserById", slog.String("request_id", requestID), slog.String("id", c.Param("id")))
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		user, err := h.userService.GetByID(ctx, id)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, user)
	}
}

func (h *userHandlers) GetUserByEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling GetUserByEmail", slog.String("request_id", requestID), slog.String("id", c.Param("id")))
		email := c.QueryParam("email")
		if email == "" {
			utils.LogResponseError(c, h.log, httpErrors.NewBadRequestError("email is required"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("email is required"))
		}
		user, err := h.userService.GetByEmail(ctx, email)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, user)
	}
}

func (h *userHandlers) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling UpdateUser", slog.String("request_id", requestID), slog.String("id", c.Param("id")))
		user := &models.User{}
		if err := utils.ReadRequestForUpdate(c, user); err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if user.ID == 0 {
			utils.LogResponseError(c, h.log, httpErrors.NewBadRequestError(httpErrors.BadRequest))
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError(httpErrors.BadRequest)))
		}

		userClaims := c.Get("user").(map[string]interface{})
		userIdFromClaims := userClaims["uid"].(float64)
		if (int64(userIdFromClaims)) != user.ID {
			utils.LogResponseError(c, h.log, httpErrors.Unauthorized)
			return c.JSON(http.StatusUnauthorized, httpErrors.Unauthorized)
		}

		user, err := h.userService.Update(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, user)
	}
}
