package server

import (
	"spotsync/internal/config"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(&user.User{}, &zone.ParkingZone{}, &reservation.Reservation{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	user.RegisterRoutes(e, db, cfg)
	zone.RegisterRoutes(e, db, cfg)
	reservation.RegisterRoutes(e, db, cfg)

	e.Start(":" + cfg.Port)
}
