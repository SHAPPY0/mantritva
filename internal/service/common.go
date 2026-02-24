package service

import (
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
)

type CommonService interface {
	CreatePreSignedUrl(string) (string, error)
	DownloadFile(*gin.Context, string) (*[]byte, string, error)
	CreateEnquiry(*gin.Context, dto.Enquiry, []*multipart.FileHeader) (bool, error)
}

type commonService struct {
	Repository  repository.EventRepository
	S3Storage   *storage.S3Storage
	EnquiryRepo repository.EnquiryRepository
}

func NewCommonService(repos *repository.Repos) CommonService {
	cs := &commonService{
		Repository:  repos.Event,
		EnquiryRepo: repos.Enquiry,
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

func (s *commonService) CreateEnquiry(c *gin.Context, data dto.Enquiry, files []*multipart.FileHeader) (bool, error) {
	en := &model.Enquiry{
		Name:        data.Name,
		Email:       data.Email,
		EnquiryType: data.EnquiryType,
		Subject:     data.Subject,
		Description: data.Description,
		FileUrls:    make([]string, 0),
	}
	storage, err := storage.New(filepath.Join("enquiries"))
	if err != nil {
		return false, err
	}
	for _, file := range files {
		fileMeta, err := storage.Save(c, file)
		if err != nil {
			log.Errorf("Error while saving file %v", err)
		} else {
			en.FileUrls = append(en.FileUrls, fileMeta.DownloadURL)
		}
	}
	if err := s.EnquiryRepo.Create(en); err != nil {
		return false, err
	}
	return true, nil
}
