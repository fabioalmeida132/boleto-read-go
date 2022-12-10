package upload

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rest-go/controller/Boleto"
	"rest-go/controller/OCR"
)

func Upload(c echo.Context) error {

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, `File not found`)
	}

	// Verify file is PDF
	if file.Header.Get("Content-Type") != "application/pdf" {
		return echo.NewHTTPError(http.StatusBadRequest, "File is not a PDF")
	}

	// Open source file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read OCR
	barCode, err := OCR.ExtractBarCode(file, c.FormValue("password"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Read boleto
	boleto, err := Boleto.ReadBoleto(file, c.FormValue("password"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//
	// Return boleto with barCode
	boleto.BarCode = barCode

	var finds []string
	if boleto.TypeableLine != "" {
		finds = append(finds, "typeableLine")
	}
	if boleto.BarCode != "" {
		finds = append(finds, "barCode")
	}

	boleto.FindTypes = finds

	return c.JSON(http.StatusOK, boleto)
}
