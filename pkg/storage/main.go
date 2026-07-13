package storage

import (
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
)

type FileMeta struct {
	OriginalName string
	SavedName    string
	FullPath     string
	AbsolutePath string
	Extenstion   string
	DownloadURL  string
	Type         string
}

type Storage interface {
	Save(*gin.Context, *multipart.FileHeader) (*FileMeta, error)
	Delete(string) error
	DeleteAll([]model.Document) error
	// Get(fileName string) (*FileMeta, error)
	// List() ([]FileMeta, error)
}

func New(destinationPath string) (Storage, error) {
	env := config.Env.GetString("ENV")
	fileBasePath := config.Env.GetString("FILE_BASE_PATH")
	switch env {
	case "DEV":
		return NewLocalStorage(fileBasePath, destinationPath)
	case "PROD":
		return NewS3Storage(fileBasePath, destinationPath)
	default:
		return nil, fmt.Errorf("No valid storage found")
	}
}
