package utils

import "os"

func GetFileSize(path string) float32 {
	fileInfo, _ := os.Stat(path)

	return float32(fileInfo.Size()) / 1024 / 1024
}

func ChangeFileName(oldPath, newPath string) {
	os.Rename(oldPath, newPath)
}

func CreateDirectoryIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func CountDigits(number int) (numberOfDigits int) {
	var counter int
	for i := 1; i <= number; i *= 10 {
		counter++
	}
	return counter
}
