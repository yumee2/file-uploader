package main

import (
	server "file-uploader/internal/http"
	db "file-uploader/internal/repository/sqlite"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	dbConn, err := db.NewDBConnection()
	if err != nil {
		fmt.Printf("Error during creating a database: %s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	http.HandleFunc("GET /files", server.GetFilesHandler)
	http.HandleFunc("GET /files/{id}", server.GetFileHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
