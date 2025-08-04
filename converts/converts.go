package converts

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	cgo_avif "github.com/Kagami/go-avif" // This is the one that requires Visual Studio, vcpkg and a bunch of other stuff, delete this one if it runs out of space
	"github.com/chai2010/webp"
	"github.com/gen2brain/avif"

	"bytes"
	"encoding/base64"
)

func ConvertImage(inputFormat, outputFormat, inputPath, outputPath string) error {

	fmt.Println(inputFormat, outputFormat, inputPath, outputPath)

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
		img, _ = webp.Decode(file)
	case "avif":
		img, _ = avif.Decode(file)
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
		return webp.Encode(outFile, img, &webp.Options{
			Lossless: true,
			Quality:  100,
		})
	case "avif":
		return cgo_avif.Encode(outFile, img, nil)
		// return avif.Encode(outFile, img, avif.Options{
		// 	Quality: 100,
		// })
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
	// Encode to desired format
	// or: return png.Encode(outFile, img)
}

func ConvertImageToBase64(path, format string) string {
	var (
		buf bytes.Buffer
		img image.Image
	)

	// Open input file
	file, _ := os.Open(path)

	switch format {
	case "jpg":
		img, _, _ = image.Decode(file)
		_ = jpeg.Encode(&buf, img, nil)

	case "jpeg":
		img, _, _ = image.Decode(file)
		_ = jpeg.Encode(&buf, img, nil)

	case "png":
		img, _, _ = image.Decode(file)
		_ = png.Encode(&buf, img)

	case "webp":
		img, _ = webp.Decode(file)

		_ = webp.Encode(&buf, img, &webp.Options{
			Lossless: true,
			Quality:  100})
	case "avif":
		img, _ = avif.Decode(file)
		_ = avif.Encode(&buf, img, avif.Options{
			Quality: 100})
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
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
