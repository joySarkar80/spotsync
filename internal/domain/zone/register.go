package zone

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	zoneRepository := NewRepository(db)
	zoneService := NewService(zoneRepository)
	zoneHandler := NewHandler(zoneService)

	jwtService := auth.NewJWTService(cfg.JwtSecret)
	authMw := middlewares.AuthMiddleware(jwtService)
	adminMw := middlewares.AdminOnly()

	api := e.Group("/api/v1/zones")

	api.GET("", zoneHandler.GetAllZones)
	api.GET("/:id", zoneHandler.GetZoneByID)
	api.POST("", zoneHandler.CreateZone, authMw, adminMw)
}
