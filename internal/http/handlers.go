package server

import (
	"fmt"
	"net/http"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET ALL FILES")
}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET ONE FILES")
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET ONE FILES")
}

func CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "CREATE FILE")
}
