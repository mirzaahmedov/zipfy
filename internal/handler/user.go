package handler

import (
	"fmt"
	"net/http"
	"time"

	"zipfy/internal/model"
	"zipfy/internal/validate"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) HandleRegister(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		c.Logger().Error(err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create user", Error: "Bad request"})
	}

	err := validate.String(&user.Name).Required("name is required").Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create user", Error: err.Error()})
	}
	err = validate.String(&user.Email).Required("email is required").Email("email must be valid email address").Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create user", Error: err.Error()})
	}
	err = validate.String(&user.Password).Required("password is required").MinLength(8, "password must be at least 8 characters").Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create user", Error: err.Error()})
	}

	email := user.Email
	encryptPassword, err := EncryptPassword(user.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create user", Error: "Bad request"})
	}
	user.Password = encryptPassword

	user, err = h.store.User().Create(user)
	if err != nil {
		c.Logger().Error(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Constraint == "unique_user_email" {
				return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create user", Error: fmt.Sprintf("User with email %s already exists", email)})
			}
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to create user", Error: "Internal server error"})
	}

	token, err := GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to create user", Error: "Internal server error"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	return c.JSON(http.StatusCreated, model.Response[*model.User]{Message: "User created", Data: user})
}
func (h *Handler) HandleLogin(c echo.Context) error {
	body := new(model.User)
	if err := c.Bind(body); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to login", Error: "Bad request"})
	}

	err := validate.String(&body.Email).Required("email is required").Email("email must be valid email address").Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to login", Error: err.Error()})
	}

	user, err := h.store.User().GetByEmail(body.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to login", Error: "Internal server error"})
	}

	if err := ComparePasswords(user.Password, body.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Failed to login", Error: "Unauthorized"})
	}

	token, err := GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to login", Error: "Internal server error"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	return c.JSON(200, model.Response[*model.User]{Message: "User logged in", Data: user})
}
func (h *Handler) HandleLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1,
	})
	return c.JSON(http.StatusOK, model.Response[interface{}]{Message: "User logged out", Data: nil})
}
func (h *Handler) HandleCheckToken(c echo.Context) error {
	id := c.Get("id").(string)
	user, err := h.store.User().GetByID(id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to check token", Error: "Internal server error"})
	}
	return c.JSON(http.StatusOK, model.Response[*model.User]{Message: "Token is valid", Data: user})
}
