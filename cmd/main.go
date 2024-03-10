package main

import (
	"log"

	"zipfy/internal/config"
	"zipfy/internal/handler"
	"zipfy/internal/store/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := echo.New()
	c, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	s := postgres.NewPostgresStore(c.Postgres.User, c.Postgres.PWD, c.Postgres.Host, c.Postgres.Port, c.Postgres.DB, c.Postgres.SSL)
	err = s.Open()
	if err != nil {
		log.Fatalf("Error opening store: %v", err)
	}
	defer func() {
		err = s.Close()
		if err != nil {
			log.Fatalf("Error closing store: %v", err)
		}
	}()

	r := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.PWD,
		DB:       c.Redis.DB,
	})
	h := handler.NewHandler(s, c.JWT.Secret, r)

	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/:short", h.HandleRedirect)

	api := e.Group("/api/v1")

	api.POST("/auth/register", h.HandleRegister)
	api.POST("/auth/login", h.HandleLogin)
	api.GET("/auth/logout", h.HandleLogout)
	api.GET("/auth/check", h.HandleCheckToken, h.AuthMiddleware)

	p := api.Group("/urls", h.AuthMiddleware)

	p.POST("", h.HandleCreateURL)
	p.GET("", h.HandleGetAllURLs)
	p.GET("/:id", h.HandleGetURLByID)
	p.PUT("/:id", h.HandleUpdateURL)
	p.DELETE("/:id", h.HandleDeleteURL)

	log.Fatal(e.Start(c.HTTP.Addr))
}
