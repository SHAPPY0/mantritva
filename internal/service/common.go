package service

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
)

type CommonService interface {
	CreatePreSignedUrl(string) (string, error)
	DownloadFile(*gin.Context, string) (*[]byte, string, error)
}

type commonService struct {
	Repository repository.EventRepository
	S3Storage  *storage.S3Storage
}

func NewCommonService(repos *repository.Repos) CommonService {
	cs := &commonService{
		Repository: repos.Event,
	}
	storage, err := storage.NewS3Storage(config.Env.GetString("FILE_BASE_PATH"), "")
	if err != nil {
		log.Errorf("Error initiating storage connection %v", err)
	} else {
		log.Infof("S3 Storage initialized successfully")
		cs.S3Storage = storage
	}
	return cs
}

func (s *commonService) CreatePreSignedUrl(filepath string) (string, error) {
	if s.S3Storage == nil {
		log.Errorf("S3 Storage not initialized")
		return "", nil
	}
	url, err := s.S3Storage.CreatePreSignedUrl(filepath)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *commonService) DownloadFile(c *gin.Context, url string) (*[]byte, string, error) {
	if s.S3Storage == nil {
		log.Errorf("S3 Storage not initialized")
		return nil, "", nil
	}
	return s.S3Storage.DownloadFile(c, url)
}
