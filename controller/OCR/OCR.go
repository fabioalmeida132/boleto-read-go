package OCR

import (
	"C"
	"bytes"
	"mime/multipart"
	"rest-go/controller/Utils"

	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
)

func ExtractBarCode(file *multipart.FileHeader, password string) (string, error) {
	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var findCode = ""
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
				return "", errors.New("File needs a password")
			}
			findCode, err = GetBarCode(doc)
			if err != nil {
				return "", err
			}

			return findCode, nil
		}

		if newFile != nil {
			if newFile.Error() == "pdfcpu: this file is not encrypted" {
				return "", errors.New("File does not need a password")
			}
		}

		return "", errors.New("File needs a correct password")
	}

	doc, err := fitz.NewFromReader(src)
	if err != nil {
		return "", errors.New("File needs a password")
	}

	findCode, err = GetBarCode(doc)
	if err != nil {
		return "", err
	}

	return findCode, nil
}

func GetBarCode(doc *fitz.Document) (string, error) {
	var findCode = ""
	for n := 0; n < doc.NumPage(); n++ {
		img, _ := doc.Image(n)

		results, _ := Utils.GetDataFromImage(img)
		for _, result := range results {
			if len(result) == 44 || len(result) == 48 {
				findCode = result
				break
			}
		}

		if findCode != "" {
			break
		}

	}
	return findCode, nil
}
