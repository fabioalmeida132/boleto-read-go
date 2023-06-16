package main

import (
	"os"
	"rest-go/controller/Upload"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	e := echo.New()
	e.Use(middleware.CORS())

	e.POST("/upload", Upload.Upload)
	e.Logger.Fatal(e.Start(":" + port))
}
