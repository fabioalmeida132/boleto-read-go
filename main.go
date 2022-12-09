package main

import (
	"github.com/labstack/echo/v4"
	"rest-go/controller/upload"
)

func main() {

	e := echo.New()

	e.POST("/upload", upload.Upload)
	e.Logger.Fatal(e.Start(":80"))
}
