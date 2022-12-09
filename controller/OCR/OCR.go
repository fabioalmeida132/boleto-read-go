package OCR

import (
	"C"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"rest-go/controller/Utils"
)

func ExtractBarCode(pathFile string, id string, password string) (string, error) {

	if password != "" {
		conf := pdfcpu.NewAESConfiguration("upw", "opw", 256)
		conf.UserPW = password
		err := api.DecryptFile(pathFile, pathFile, conf)
		if err != nil {
			if err.Error() == "pdfcpu: please provide the correct password" {
				// remove file
				err = Utils.RemoveFile(pathFile)
				if err != nil {
					return "", err
				}
				return "", errors.New("Password is incorrect")
			}
			Utils.RemoveFile(pathFile)
			return "", errors.New("File does not need password")
		}

	}

	doc, err := fitz.New(pathFile)
	if err != nil {
		Utils.RemoveFile(pathFile)
		return "", errors.New("File needs a password")
	}

	folder := "tmp/"

	// Extract pages as images
	var findCode = ""
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(folder+"/", fmt.Sprintf(id+"-%05d.jpeg", n)))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}

		results, err := GetDataFromFile(filepath.Join(folder+"/", fmt.Sprintf(id+"-%05d.jpeg", n)))
		if err != nil {
			log.Fatal(err)
		}

		for _, result := range results {
			findCode = result
			break
		}

		f.Close()

		err = os.Remove(filepath.Join(folder+"/", fmt.Sprintf(id+"-%05d.jpeg", n)))
		if err != nil {
			log.Fatal(err)
		}

		if findCode != "" {
			break
		}

	}

	return findCode, nil
}
