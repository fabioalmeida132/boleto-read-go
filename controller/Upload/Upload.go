package Upload

import "C"
import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"rest-go/controller/Boleto"
	"rest-go/controller/OCR"
	"rest-go/controller/Utils"
	"rest-go/models"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/labstack/echo/v4"
	"github.com/signintech/gopdf"
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

	if file.Header.Get("Content-Type") == "image/jpeg" || file.Header.Get("Content-Type") == "image/png" || file.Header.Get("Content-Type") == "image/jpg" {
		// Abrir o arquivo da imagem
		imageFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer imageFile.Close()

		// Ler os dados do arquivo da imagem
		imageData, err := ioutil.ReadAll(imageFile)
		if err != nil {
			log.Fatal(err)
		}

		// Decodificar a imagem
		img, format, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			log.Fatal(err)
		}

		// Criar um novo objeto PDF
		pdf := gopdf.GoPdf{}

		// Definir a configuração do PDF
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.AddPage()

		// Criar um buffer de bytes para armazenar a imagem no formato JPEG ou PNG
		imageBuffer := new(bytes.Buffer)

		// Codificar a imagem para o formato JPEG ou PNG e salvar no buffer
		err = encodeImage(format, img, imageBuffer)
		if err != nil {
			log.Fatal(err)
		}

		// Carregar a imagem a partir dos dados do buffer
		imageObj, err := gopdf.ImageHolderByBytes(imageBuffer.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		// Adicionar a imagem ao PDF
		pdf.ImageByHolder(imageObj, 0, 0, nil)

		// Criar um buffer de bytes em memória para o PDF
		buffer := new(bytes.Buffer)

		// Salvar o PDF no buffer
		err = pdf.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// Atualizar os dados do arquivo `file` com o conteúdo do PDF gerado
		file.Header = make(textproto.MIMEHeader)
		file.Header.Set("Content-Type", "application/pdf")
		file.Size = int64(buffer.Len())
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

	if len(barCode) >= 47 {
		boleto.TypeableLine = barCode
		boleto.BarCode = ""
	} else {
		boleto.BarCode = barCode

		// Read boleto
		typeableLine, err := Boleto.ReadBoleto(file, c.FormValue("password"))
		if err != nil {
			boleto.Message = err.Error()
			return c.JSON(http.StatusBadRequest, boleto)
		}

		// Return boleto with barCode
		boleto.TypeableLine = typeableLine

	}

	var finds []string
	if boleto.TypeableLine != "" {
		finds = append(finds, "typeableLine")
	}
	if boleto.BarCode != "" {
		finds = append(finds, "barCode")

		if boleto.TypeableLine == "" {
			var calc = calcTypeableLine(boleto.BarCode)
			boleto.TypeableLine = calc

			if calc != "" {
				finds = append(finds, "typeableLine")
			}
		}

	}

	boleto.FindTypes = finds

	return c.JSON(http.StatusOK, boleto)
}

func calcTypeableLine(barcode string) string {
	barcodeLine := Utils.Number(barcode)

	if len(barcodeLine) != 44 {
		return ""
	}

	field1 := fmt.Sprint(barcodeLine[0:4], barcodeLine[19:24])
	field2 := fmt.Sprint(barcodeLine[24:29], barcodeLine[29:34])
	field3 := fmt.Sprint(barcodeLine[34:39], barcodeLine[39:44])
	field4 := barcodeLine[4:5]
	field5 := barcodeLine[5:19]

	if Utils.Mod11(barcodeLine[0:4]+barcodeLine[5:44]) != field4 {
		fmt.Println("are not the same")
	}

	return fmt.Sprint(field1, Utils.Mod10(field1), field2, Utils.Mod10(field2), field3, Utils.Mod10(field3), field4, field5)
}

// Função para codificar a imagem no formato correto
func encodeImage(format string, img image.Image, buffer *bytes.Buffer) error {
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		return jpeg.Encode(buffer, img, nil)
	case "png":
		return png.Encode(buffer, img)
	default:
		return fmt.Errorf("Formato de imagem não suportado: %s", format)
	}
}
