package OCR

import (
	"C"
	"bytes"
	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"mime/multipart"
	"rest-go/controller/Utils"
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
		img, err := doc.Image(n)

		results, err := Utils.GetDataFromImage(img)
		if err != nil {
			return "", errors.New("Error to extract barcode")
		}

		for _, result := range results {
			findCode = result
			break
		}

		if findCode != "" {
			break
		}

	}
	return findCode, nil
}
