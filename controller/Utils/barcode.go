package Utils

// #cgo darwin pkg-config: zbar
// #cgo LDFLAGS: -lzbar
// #include <zbar.h>
import "C"
import (
	"bytes"
	"github.com/otiai10/gosseract/v2"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"regexp"
	"strings"
)

type RawData string

func imageToBytes(img image.Image) []byte {
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	return buf.Bytes()
}

func extractBarcodeSequence(text string) string {
	// Remove os caracteres especiais
	replaced := strings.ReplaceAll(text, ".", " ")
	replaced = strings.ReplaceAll(replaced, "-", " ")

	// Procura pela primeira sequência
	r := regexp.MustCompile(`(\d{12}\s\d{12}\s\d{12}\s\d{12})`)
	match := r.FindString(replaced)

	// Caso não encontrado, tenta a segunda sequência
	if match == "" {
		r = regexp.MustCompile(`(\d{5}\s\d{5}\s\d{5}\s\d{6}\s\d\s\d{20})`)
		match = r.FindString(replaced)
	}

	// Caso não encontrado, tenta a terceira sequência
	if match == "" {
		r = regexp.MustCompile(`(\d{5}\s\d{5}\s\d{5}\s\d{6}\s\d\s\d{3})`)
		match = r.FindString(replaced)
	}

	// Caso não encontrado, tenta a quarta sequência
	if match == "" {
		r = regexp.MustCompile(`(\d{5}\s\d{5}\s\d{5}\s\d{6}\s\d{5}\s\d{6}\s\d\s\d{14})`)
		match = r.FindString(replaced)
	}

	return match
}

func GetDataFromImage(image image.Image) (results []string, err error) {

	// Tenta extrair o código de barras

	scanner := NewScanner()
	defer scanner.Close()
	scanner.SetConfig(0, C.ZBAR_CFG_ENABLE, 1)
	zImg := NewZbarImage(image)
	defer zImg.Close()
	scanner.Scan(zImg)
	symbol := zImg.GetSymbol()

	for ; symbol != nil; symbol = symbol.Next() {
		if symbol.Type().t == 25 {
			if Mod11(symbol.Data()[0:4]+symbol.Data()[5:44]) != symbol.Data()[4:5] {
				continue
			}
			results = append(results, symbol.Data())
		}
	}

	if len(results) > 0 && len(results[0]) == 44 {
		return results, nil
	}

	// Se a extração do código de barras não retornar um resultado válido, tenta OCR

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(imageToBytes(image))
	ocrText, _ := client.Text()

	typeableLineSequence := extractBarcodeSequence(ocrText)
	typeableLineSequence = strings.ReplaceAll(typeableLineSequence, " ", "")

	if len(typeableLineSequence) == 47 || len(typeableLineSequence) == 48 {
		return []string{typeableLineSequence}, nil
	}

	return results, nil
}
