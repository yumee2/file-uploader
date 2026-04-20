package service

import (
	"file-uploader/internal/service/dto"
	storage "file-uploader/internal/storage/filesystem"
	"file-uploader/models"
	"io"
	"mime"
	"path/filepath"

	"github.com/google/uuid"
)

type DataBaseI interface {
	AddFile(file *models.File) error
	GetFile(id string) (*models.File, error)
	GetFiles() ([]*models.File, error)
	DeleteFile(id string) error
}

type FileService struct {
	db DataBaseI
}

func NewFileService(db DataBaseI) *FileService {
	return &FileService{db: db}
}

func (fs *FileService) AddFile(fileDTO *dto.UploadDTO) (string, error) {
	fileUUID := uuid.NewString()
	size, err := storage.SaveFile(fileUUID, fileDTO.Body)
	if err != nil {
		return "", err
	}

	file := &models.File{
		ID:           fileUUID,
		OriginalName: filepath.Base(fileDTO.Name),
		MimeType:     mime.TypeByExtension(filepath.Ext(fileDTO.Name)),
		Size:         size,
	}

	err = fs.db.AddFile(file)
	if err != nil {
		return "", err
	}

	return fileUUID, nil
}

func (fs *FileService) DownloadFile(id string, w io.Writer) (*models.File, func(io.Writer) error, error) {
	file, err := fs.db.GetFile(id)
	if err != nil {
		return nil, nil, err
	}

	streamFunc := func(w io.Writer) error {
		return storage.WriteFileTo(id, w)
	}
	return file, streamFunc, nil
}

func (fs *FileService) GetFiles() ([]*models.File, error) {
	files, err := fs.db.GetFiles()
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fs *FileService) DeleteFile(id string) error {
	err := storage.DeleteFile(id)
	if err != nil {
		return err
	}
	return fs.db.DeleteFile(id)
}
