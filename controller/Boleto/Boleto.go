package Boleto

import (
	"github.com/gen2brain/go-fitz"
	"regexp"
	"rest-go/controller/Utils"
	"rest-go/models"
)

func ReadBoleto(pathFile string) (models.Boleto, error) {

	var typeableLine string
	doc, err := fitz.New(pathFile)
	if err != nil {
		Utils.RemoveFile(pathFile)
	}

	for n := 0; n < doc.NumPage(); n++ {
		text, err := doc.Text(n)
		if err != nil {
			return models.Boleto{}, err
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

	err = Utils.RemoveFile(pathFile)
	if err != nil {
		return models.Boleto{}, err
	}

	return models.Boleto{TypeableLine: typeableLine}, nil
}
