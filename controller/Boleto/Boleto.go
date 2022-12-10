package Boleto

import (
	"bytes"
	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"mime/multipart"
	"regexp"
	"rest-go/models"
)

func ReadBoleto(file *multipart.FileHeader, password string) (models.Boleto, error) {

	// Open source file
	src, err := file.Open()
	if err != nil {
		return models.Boleto{}, err
	}
	defer src.Close()

	var typeableLine = ""
	if password != "" {
		conf := pdfcpu.NewAESConfiguration("upw", "opw", 256)
		conf.UserPW = password
		conf.UserPWNew = nil
		buf := new(bytes.Buffer)
		newFile := api.Decrypt(src, buf, conf)
		if newFile == nil {
			doc, err := fitz.NewFromMemory(buf.Bytes())
			defer doc.Close()
			if err != nil {
				return models.Boleto{}, errors.New("File needs a password")
			}
			typeableLine, err = GetTypeableLine(doc)
			if err != nil {
				return models.Boleto{}, err
			}
			return models.Boleto{TypeableLine: typeableLine}, nil
		}

		return models.Boleto{}, errors.New("File needs a correct password")
	}

	doc, err := fitz.NewFromReader(src)
	if err != nil {
		return models.Boleto{}, err
	}
	defer doc.Close()

	typeableLine, err = GetTypeableLine(doc)
	if err != nil {
		return models.Boleto{}, err
	}

	return models.Boleto{TypeableLine: typeableLine}, nil
}

func GetTypeableLine(doc *fitz.Document) (string, error) {
	var typeableLine = ""
	for n := 0; n < doc.NumPage(); n++ {
		text, err := doc.Text(n)
		if err != nil {
			return "", err
		}

		lines := regexp.MustCompile(`\r?\n`).Split(text, -1)
		for _, line := range lines {
			re := regexp.MustCompile("[^0-9]+")
			onlyNumbers := re.ReplaceAllString(line, "")
			if len(onlyNumbers) == 51 {
				onlyNumbers = onlyNumbers[3:]
			}
			if len(onlyNumbers) == 47 || len(onlyNumbers) == 48 {
				typeableLine = onlyNumbers
				break
			}
		}
	}
	return typeableLine, nil
}
