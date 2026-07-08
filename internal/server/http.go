package server

import (
	"haddibanga/internal/config"
	"haddibanga/internal/domain/mango"
	"haddibanga/internal/domain/order"
	"haddibanga/internal/domain/user"

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
	db.AutoMigrate(&user.User{}, &mango.Mango{}, &order.Order{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	user.RegisterRoutes(e, db, cfg)
	mango.RegisterRoutes(e, db, cfg)
	order.RegisterRoutes(e, db, cfg)

	e.Start(":" + cfg.Port)
}
