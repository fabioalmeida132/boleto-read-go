package main

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"os"
	"rest-go/controller/Upload"
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
