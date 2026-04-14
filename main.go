package main

import (
	"file-uploader/db"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const chunkSizeKB = 500
const chunkSizeBytes = 500 * 1024

var rootDir, _ = os.Getwd()

func main() {
	_, err := db.InitDB()
	if err != nil {
		fmt.Printf("Error during creating a database: %s", err)
		os.Exit(1)
	}

	filePath := os.Args[1]
	fileDirPath, _ := breakFileIntoChunks(filePath)
	restoreFile(fileDirPath, filepath.Base(filePath))
}

func breakFileIntoChunks(filePath string) (string, error) {
	data, _ := os.ReadFile(filePath) // TODO: read file in streams

	fileDirPath := uuid.NewString()
	if err := os.MkdirAll(fileDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create chunk directory: %w", err)
	}

	var chunkIndex = 0
	for i := 0; i < len(data); i += chunkSizeBytes {
		end := min(i+chunkSizeBytes, len(data))
		chunk := data[i:end]
		chunkPath := filepath.Join(rootDir, fileDirPath, fmt.Sprintf("chunk_%d", chunkIndex))
		os.WriteFile(chunkPath, chunk, 0644)
		chunkIndex++
	}
	return fileDirPath, nil
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

//TODO:
// sqlite files database
// running the server
// fixes: read the file in streams
