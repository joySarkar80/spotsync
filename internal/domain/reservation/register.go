package reservation

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	reservationRepository := NewRepository(db)
	reservationService := NewService(reservationRepository)
	reservationHandler := NewHandler(reservationService)

	jwtService := auth.NewJWTService(cfg.JwtSecret)
	authMw := middlewares.AuthMiddleware(jwtService)
	adminMw := middlewares.AdminOnly()

	api := e.Group("/api/v1/reservations", authMw)

	api.POST("", reservationHandler.Reserve)
	api.GET("/my-reservations", reservationHandler.GetMyReservations)
	api.DELETE("/:id", reservationHandler.CancelReservation)
	api.GET("", reservationHandler.GetAllReservations, adminMw)
}
