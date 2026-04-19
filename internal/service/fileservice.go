package service

import "file-uploader/models"

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

func (fs *FileService) AddFile(file *models.File) error {
	return fs.db.AddFile(file)
}

func (fs *FileService) GetFile(id string) (*models.File, error) {
	return fs.db.GetFile(id)
}

func (fs *FileService) GetFiles() ([]*models.File, error) {
	return fs.db.GetFiles()
}

func (fs *FileService) DeleteFile(id string) error {
	return fs.db.DeleteFile(id)
}
