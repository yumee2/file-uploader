package http

import (
	"file-uploader/models"
	"fmt"
	"io"
	"net/http"
)

type FileServiceI interface {
	AddFile(file *models.File) error
	DownloadFile(id string, w io.Writer) (*models.File, error)
	GetFiles() ([]*models.File, error)
	DeleteFile(id string) error
}

type FileHandler struct {
	service FileServiceI
}

func NewFileHandler(service FileServiceI) *FileHandler {
	return &FileHandler{service: service}
}

func (h *FileHandler) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET ALL FILES")
}

func (h *FileHandler) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET ONE FILES")
}

func (h *FileHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "DELETE FILE")
}

func (h *FileHandler) CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "CREATE FILE")
}
