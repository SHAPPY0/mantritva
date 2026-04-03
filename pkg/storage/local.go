package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
)

type LocalStorage struct {
	BasePath         string
	AbsoluteDestPath string
}

func NewLocalStorage(basePath string, absoluteDestPath string) (*LocalStorage, error) {
	_ = os.MkdirAll(basePath+absoluteDestPath, os.ModePerm)
	return &LocalStorage{
		BasePath:         basePath,
		AbsoluteDestPath: absoluteDestPath,
	}, nil
}

func (l *LocalStorage) Save(c *gin.Context, file *multipart.FileHeader) (*FileMeta, error) {
	originalFileName := filepath.Base(file.Filename)
	extension := strings.Split(originalFileName, ".")[1]
	newName := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	fullPath := filepath.Join(l.BasePath, l.AbsoluteDestPath, newName)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		log.Errorf("Failed to save file: %v", err)
		return nil, err
	}
	fileMeta := &FileMeta{
		OriginalName: originalFileName,
		SavedName:    newName,
		FullPath:     fullPath,
		AbsolutePath: filepath.Join("/", l.AbsoluteDestPath, newName),
		Extenstion:   extension,
		DownloadURL:  "http://localhost:" + config.GetConfig().GetString("PORT") + "/static/" + l.AbsoluteDestPath + "/" + newName,
	}
	log.Infof("File uploaded successfully: " + fileMeta.AbsolutePath)
	return fileMeta, nil
}

func (l *LocalStorage) Delete(filePath string) error {
	if filePath != "" {
		return fmt.Errorf("Invalid filepath to delete")
	}
	if err := os.Remove(filePath); err != nil {
		log.Errorf("Failed to delete file: %v", err)
		return err
	}
	log.Infof("File deleted successfully: " + filePath)
	return nil
}

func (l *LocalStorage) DeleteAll(documents []model.Document) error {
	for _, doc := range documents {
		if err := l.Delete(doc.FileName); err != nil {
			log.Errorf("Failed to delete file: %v", err)
			continue
		}
	}
	return nil
}
