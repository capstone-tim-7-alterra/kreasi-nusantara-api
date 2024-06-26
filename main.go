package main

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/drivers/database"
	"kreasi-nusantara-api/routes"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var db *gorm.DB
var v *validation.Validator

func init() {
	config.LoadEnv()
	config.InitConfigDB()
	db = database.ConnectDB(config.InitConfigDB())
	v = validation.NewValidator()
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://kreasi-nusantara-batik-website-testing.vercel.app",
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

	routes.InitRoute(e, db, v)

	e.Logger.Fatal(e.Start(":8080"))
}