package storage

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath         string
	AbsoluteDestPath string
}

func NewLocalStorage(basePath string, absoluteDestPath string) *LocalStorage {
	_ = os.MkdirAll(basePath+absoluteDestPath, os.ModePerm)
	return &LocalStorage{
		BasePath:         basePath,
		AbsoluteDestPath: absoluteDestPath,
	}
}

func (l *LocalStorage) Save(c *gin.Context, file *multipart.FileHeader) (*FileMeta, error) {
	originalFileName := filepath.Base(file.Filename)
	extension := strings.Split(originalFileName, ".")[1]
	newName := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	fullPath := filepath.Join(l.BasePath, l.AbsoluteDestPath, newName)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		log.Printf("Failed to save file: %v", err)
		return nil, err
	}
	fileMeta := &FileMeta{
		OriginalName: originalFileName,
		SavedName:    newName,
		FullPath:     fullPath,
		AbsolutePath: filepath.Join("/", l.AbsoluteDestPath, newName),
	}
	log.Printf("File uploaded successfully: " + fileMeta.AbsolutePath)
	return fileMeta, nil
}

func (l *LocalStorage) Delete(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		log.Printf("Failed to delete file: %v", err)
		return err
	}
	log.Printf("File deleted successfully: " + filePath)
	return nil
}
