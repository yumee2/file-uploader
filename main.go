package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const chunkSizeKB = 500

func main() {
	rootDir, _ := os.Getwd()
	path := filepath.Join(rootDir, "image.jpg")
	data, _ := os.ReadFile(path)

	fmt.Println(len(data), len(data)/1024, len(data)/(1024*1024))

	var chunk = make([]byte, 0, chunkSizeKB*1024)
	var chunkIndex = 0
	for _, b := range data {
		chunk = append(chunk, b)
		if len(chunk) >= chunkSizeKB*1024 {
			chunkPath := filepath.Join(rootDir, fmt.Sprintf("chunk_%d", chunkIndex))
			os.WriteFile(chunkPath, chunk, 0644)
			chunkIndex++
			chunk = make([]byte, 0, chunkSizeKB*1024)
		}
	}

	if len(chunk) > 0 {
		chunkPath := filepath.Join(rootDir, fmt.Sprintf("chunk_%d", chunkIndex))
		os.WriteFile(chunkPath, chunk, 0644)
	}
}
