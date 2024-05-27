package main

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/drivers/database"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	config.LoadEnv()
	config.InitConfigDB()
	db = database.ConnectDB(config.InitConfigDB())
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXCSRFToken,
		},
		AllowMethods: []string{
			echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
	e.Logger.Fatal(e.Start(":8080"))
}