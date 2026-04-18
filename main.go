package main

import (
	"file-uploader/db"
	server "file-uploader/internal/http"
	"file-uploader/models"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const chunkSizeKB = 500
const chunkSizeBytes = 500 * 1024

var rootDir, _ = os.Getwd()

func main() {
	dbConn, err := db.InitDB()
	if err != nil {
		fmt.Printf("Error during creating a database: %s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	//filePath := os.Args[1]

	// file, _ := breakFileIntoChunks(filePath)
	// fmt.Println(file.ID, file.OriginalName, file.Size, file.MimeType)
	// db.AddFile(dbConn, file)
	// restoreFile(fileDirPath, filepath.Base(filePath))
	http.HandleFunc("GET /files", server.GetFilesHandler)
	http.HandleFunc("GET /files/{id}", server.GetFileHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func breakFileIntoChunks(filePath string) (*models.File, error) {
	fileDirPath := uuid.NewString()

	data, _ := os.ReadFile(filePath) // TODO: read file in streams

	if err := os.MkdirAll(fileDirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chunk directory: %w", err)
	}

	var chunkIndex = 0
	for i := 0; i < len(data); i += chunkSizeBytes {
		end := min(i+chunkSizeBytes, len(data))
		chunk := data[i:end]
		chunkPath := filepath.Join(rootDir, fileDirPath, fmt.Sprintf("chunk_%d", chunkIndex))
		os.WriteFile(chunkPath, chunk, 0644)
		chunkIndex++
	}
	return &models.File{
		ID:           fileDirPath,
		OriginalName: filepath.Base(filePath),
		Size:         int64(len(data)),
		MimeType:     mime.TypeByExtension(filepath.Ext(filePath)),
	}, nil
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
// fixes: read the file in streams
// encapsulate file logic in different package
