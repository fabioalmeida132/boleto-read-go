package OCR

import "C"
import (
	"fmt"
	"github.com/karmdip-mi/go-fitz"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func ExtractBarCode(pathFile string, id string) (string, error) {

	doc, err := fitz.New(pathFile)
	if err != nil {
		panic(err)
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

		//err = os.Remove(filepath.Join(folder+"/", fmt.Sprintf(id+"-%05d.jpeg", n)))
		//if err != nil {
		//	log.Fatal(err)
		//}

		if findCode != "" {
			break
		}

	}

	return findCode, nil
}
