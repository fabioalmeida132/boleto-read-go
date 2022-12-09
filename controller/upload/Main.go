package upload

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
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
		return c.String(http.StatusBadRequest, "File is not a PDF")
	}

	// Open source file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create id to file
	id := uuid.New().String()

	// Create destination file
	dst, err := os.OpenFile("tmp/"+id+".pdf", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy source file to destination file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Read OCR
	var barCode = ""
	barCode, err = OCR.ExtractBarCode("tmp/"+id+".pdf", id, c.FormValue("password"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Read boleto
	boleto, err := Boleto.ReadBoleto("tmp/" + id + ".pdf")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

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
