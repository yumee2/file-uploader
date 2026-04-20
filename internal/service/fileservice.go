package service

import (
	"file-uploader/internal/service/dto"
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

type FileStorageI interface {
	SaveFile(uuid string, data io.Reader) (int64, error)
	WriteFileTo(fileUUID string, w io.Writer) error
	DeleteFile(fileUUID string) error
}

type FileService struct {
	db      DataBaseI
	storage FileStorageI
}

func NewFileService(db DataBaseI, storage FileStorageI) *FileService {
	return &FileService{db: db, storage: storage}
}

func (fs *FileService) AddFile(fileDTO *dto.UploadDTO) (string, error) {
	fileUUID := uuid.NewString()
	size, err := fs.storage.SaveFile(fileUUID, fileDTO.Body)
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

func (fs *FileService) DownloadFile(id string, w io.Writer) (*models.File, error) {
	file, err := fs.db.GetFile(id)
	if err != nil {
		return nil, err
	}

	err = fs.storage.WriteFileTo(id, w)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fs *FileService) GetFiles() ([]*models.File, error) {
	files, err := fs.db.GetFiles()
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fs *FileService) DeleteFile(id string) error {
	err := fs.storage.DeleteFile(id)
	if err != nil {
		return err
	}
	return fs.db.DeleteFile(id)
}
