package converts

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nickalie/go-webpbin"
)

func ConvertImage(format, inputPath, outputPath string) error {
	// Open input file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image (auto-detects format)
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "png":
		return png.Encode(outFile, img)
	case "jpeg":
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100})
	case "jpg":
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100})
	case "webp":
		return webpbin.Encode(outFile, img)
	default:
		return nil
	}
	// Encode to desired format
	// or: return png.Encode(outFile, img)
}
