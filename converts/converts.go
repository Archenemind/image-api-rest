package converts

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nickalie/go-webpbin"
)

func ConvertImage(inputFormat, outputFormat, inputPath, outputPath string) error {
	// Open input file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var img image.Image

	switch inputFormat {
	case "jpg":
		img, _, _ = image.Decode(file)
	case "jpeg":
		img, _, _ = image.Decode(file)
	case "png":
		img, _, _ = image.Decode(file)
	case "webp":
		img, _ = webpbin.Decode(file)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch outputFormat {
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

func DeleteImages(paths []string) error {
	for i := range paths {
		err := os.Remove(paths[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFileSize(path string) float32 {
	fileInfo, _ := os.Stat(path)

	return float32(fileInfo.Size()) / 1024 / 1024
}

func ChangeFileName(oldPath, newPath string) {
	os.Rename(oldPath, newPath)
}
