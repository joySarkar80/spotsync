package mango

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("/api/v1/mangoes")

	api.GET("", handler.GetMangoes)
	api.GET("/:id", handler.GetMangoByID)
	api.POST("", handler.CreateMango, middlewares.AuthMiddleware(jwtService))
	api.PATCH("/:id", handler.UpdateMango, middlewares.AuthMiddleware(jwtService))
}
