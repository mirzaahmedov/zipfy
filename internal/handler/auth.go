package handler

import (
	"net/http"

	"zipfy/internal/model"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenCookie, err := c.Cookie("token")
		if err != nil || tokenCookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Unauthorized", Error: "Unauthorized"})
		}

		id, err := ParseToken(tokenCookie.Value, h.jwtSecret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Unauthorized", Error: "Unauthorized"})
		}

		if id == "" {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Unauthorized", Error: "Unauthorized"})
		}

		c.Set("id", id)

		return next(c)
	}
}
