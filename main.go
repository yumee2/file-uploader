package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const chunkSizeKB = 500
const chunkSizeBytes = 500 * 1024

var rootDir, _ = os.Getwd()

func main() {
	filePath := os.Args[1]
	fileDirPath := breakFileIntoChunks(filePath)
	restoreFile(fileDirPath, filepath.Base(filePath))
}

func breakFileIntoChunks(filePath string) (fileDirPath string) {
	fileName := filepath.Base(filePath)
	data, _ := os.ReadFile(filePath)

	fileDirPath = strings.Join(strings.Split(fileName, "."), "_")
	os.Mkdir(fileDirPath, 0755)

	var chunkIndex = 0
	for i := 0; i < len(data); i += chunkSizeBytes {
		end := min(i+chunkSizeBytes, len(data))
		chunk := data[i:end]
		chunkPath := filepath.Join(rootDir, fileDirPath, fmt.Sprintf("chunk_%d", chunkIndex))
		os.WriteFile(chunkPath, chunk, 0644)
		chunkIndex++
	}
	return fileDirPath
}

func restoreFile(fileDir, fileName string) {
	var fileData []byte
	entities, _ := os.ReadDir(fileDir)

	for _, entity := range entities {
		filePath := filepath.Join(fileDir, entity.Name())
		chunkData, _ := os.ReadFile(filePath)
		fileData = append(fileData, chunkData...)
	}

	resultFilePath := filepath.Join(rootDir, fileName)
	os.WriteFile(resultFilePath, fileData, 0644)
}
