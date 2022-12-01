package Boleto

import (
	"code.sajari.com/docconv"
	"errors"
	"os"
	"regexp"
	"rest-go/models"
	"strings"
)

func ReadBoleto(id string) (models.Boleto, error) {
	res, err := docconv.ConvertPath("tmp/" + id + ".pdf")
	if err != nil {
		return models.Boleto{}, err
	}

	err = os.Remove("tmp/" + id + ".pdf")
	if err != nil {
		return models.Boleto{}, err
	}

	result := strings.Split(res.Body, "\n")
	var typeableLine string
	for _, line := range result {
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
	if typeableLine != "" {
		return models.Boleto{TypeableLine: typeableLine}, nil
	} else {
		return models.Boleto{}, errors.New("typeable line not found")
	}
}
