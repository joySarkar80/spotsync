package order

import (
	"haddibanga/internal/auth"
	"haddibanga/internal/config"
	"haddibanga/internal/domain/mango"
	"haddibanga/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	orderRepo := NewRepository(db)
	mangoRepo := mango.NewRepository(db)

	svc := NewService(orderRepo, mangoRepo)
	handler := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("/api/v1/orders", middlewares.AuthMiddleware(jwtService))

	api.POST("", handler.CreateOrder)
	api.GET("/me", handler.GetMyOrders)
}
