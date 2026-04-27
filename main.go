package main

import (
	server "file-uploader/internal/http"
	"file-uploader/internal/repository/sqlite"
	"file-uploader/internal/service"
	"log"
	"net/http"

	"fmt"
	"os"
)

func main() {
	dbConn, err := sqlite.NewDBConnection()
	if err != nil {
		fmt.Printf("Error during creating a database: %s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	fileService := service.NewFileService(dbConn)
	fileHandler := server.NewFileHandler(fileService)

	http.HandleFunc("GET /files", fileHandler.GetFiles)
	http.HandleFunc("GET /files/{id}", fileHandler.DownloadFile)
	http.HandleFunc("DELETE /files/{id}", fileHandler.DeleteFile)
	http.HandleFunc("POST /files", fileHandler.AddFile)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// TODO:
// better error handling
// graceful shutdown
// rollback on failure
// user authentication
// do chunking on a client side so i can upload large files

// FIXES:
// service layer: DELETE from database
