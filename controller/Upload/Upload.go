package Upload

import "C"
import (
	"github.com/gen2brain/go-fitz"
	"github.com/labstack/echo/v4"
	"net/http"
	"rest-go/controller/Boleto"
	"rest-go/controller/OCR"
	"rest-go/models"
)

func Upload(c echo.Context) error {

	// initialize boleto
	boleto := models.Boleto{}

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		boleto.Message = "File not found"
		return c.JSON(http.StatusBadRequest, boleto)
	}

	// Verify file is PDF
	if file.Header.Get("Content-Type") != "application/pdf" {
		boleto.Message = "File is not a PDF"
		return c.JSON(http.StatusBadRequest, boleto)
	}

	// Open source file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// Verify need password
	_, err = fitz.NewFromReader(src)
	if err != nil {
		boleto.HasPassword = true
	}

	// Read OCR
	barCode, err := OCR.ExtractBarCode(file, c.FormValue("password"))
	if err != nil {
		if err.Error() == "File needs a password" {
			boleto.HasPassword = true
		}
		boleto.Message = err.Error()
		return c.JSON(http.StatusBadRequest, boleto)
	}

	// Read boleto
	typeableLine, err := Boleto.ReadBoleto(file, c.FormValue("password"))
	if err != nil {
		boleto.Message = err.Error()
		return c.JSON(http.StatusBadRequest, boleto)
	}

	// Return boleto with barCode
	boleto.TypeableLine = typeableLine
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
