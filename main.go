package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const chunkSizeKB = 500

var rootDir, _ = os.Getwd()

func main() {
	rootDir, _ := os.Getwd()
	path := filepath.Join(rootDir, "image.jpg")
	data, _ := os.ReadFile(path)

	os.Mkdir("image_jpg", 0755)

	var chunk = make([]byte, 0, chunkSizeKB*1024)
	var chunkIndex = 0
	for _, b := range data {
		chunk = append(chunk, b)
		if len(chunk) >= chunkSizeKB*1024 {
			chunkPath := filepath.Join(rootDir, "image_jpg", fmt.Sprintf("chunk_%d", chunkIndex))
			os.WriteFile(chunkPath, chunk, 0644)
			chunkIndex++
			chunk = make([]byte, 0, chunkSizeKB*1024)
		}
	}

	if len(chunk) > 0 {
		chunkPath := filepath.Join(rootDir, "image_jpg", fmt.Sprintf("chunk_%d", chunkIndex))
		os.WriteFile(chunkPath, chunk, 0644)
	}

	restoreFile(filepath.Join(rootDir, "image_jpg"))
}

func restoreFile(fileDir string) {
	var fileData []byte
	entities, _ := os.ReadDir(fileDir)

	for _, entity := range entities {
		filePath := filepath.Join(fileDir, entity.Name())
		chunkData, _ := os.ReadFile(filePath)
		fileData = append(fileData, chunkData...)
	}

	resultFilePath := filepath.Join(rootDir, "restored_image.jpg")
	os.WriteFile(resultFilePath, fileData, 0644)
}
