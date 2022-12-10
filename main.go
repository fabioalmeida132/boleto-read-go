package main

import (
	"github.com/labstack/echo/v4"
	"os"
	"rest-go/controller/upload"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.POST("/upload", upload.Upload)
	e.Logger.Fatal(e.Start(":" + port))
}
