package handler

import (
	"net/http"
	"strconv"
	"time"

	"zipfy/internal/generate"
	"zipfy/internal/model"
	"zipfy/internal/validate"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateURL(c echo.Context) error {
	url := new(model.URL)
	if err := c.Bind(url); err != nil {
		c.Logger().Error(err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create url", Error: "Bad request"})
	}
	url.UserID = c.Get("id").(string)

	err := validate.String(&url.URL).Required("url is required").URL("url must be valid url").Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to create url", Error: err.Error()})
	}

	url, err = h.store.URL().Create(url)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to create url", Error: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.URL]{Message: "URL created", Data: url})
}
func (h *Handler) HandleGetAllURLs(c echo.Context) error {
	urls, err := h.store.URL().GetByUserID(c.Get("id").(string))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to get urls", Error: "Internal server error"})
	}

	return c.JSON(http.StatusOK, model.Response[[]*model.URL]{Message: "URLs found", Data: urls})
}
func (h *Handler) HandleGetURLByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to get url", Error: "Bad request"})
	}

	url, err := h.store.URL().GetByID(id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to get url", Error: "Internal server error"})
	}

	return c.JSON(http.StatusOK, model.Response[*model.URL]{Message: "URL found", Data: url})
}
func (h *Handler) HandleUpdateURL(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to update url", Error: "Bad request"})
	}

	url := new(model.URL)
	if err := c.Bind(url); err != nil {
		c.Logger().Error(err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to update url", Error: "Bad request"})
	}
	url.ID = id

	err = validate.String(&url.URL).Required("url is required").URL("url must be valid url").Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to update url", Error: err.Error()})
	}

	url, err = h.store.URL().Update(id, url.URL)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to update url", Error: "Internal server error"})
	}

	return c.JSON(http.StatusOK, model.Response[*model.URL]{Message: "URL updated", Data: url})
}
func (h *Handler) HandleDeleteURL(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to delete url", Error: "Bad request"})
	}

	err = h.store.URL().Delete(id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to delete url", Error: "Internal server error"})
	}

	return c.JSON(http.StatusOK, model.Response[string]{Message: "URL deleted", Data: "URL deleted"})
}
func (h *Handler) HandleRedirect(c echo.Context) error {
	cached, redisErr := h.redis.Get(c.Request().Context(), c.Param("short")).Result()
	if redisErr == nil {
		return c.Redirect(http.StatusMovedPermanently, cached)
	}

	id := generate.ParseIntFromUUID(c.Param("short"))

	url, err := h.store.URL().GetByID(id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to get url", Error: "Internal server error"})
	}

	h.redis.Set(c.Request().Context(), c.Param("short"), url.URL, 48*time.Hour)
	return c.Redirect(http.StatusMovedPermanently, url.URL)
}
