package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shappy0/mantritva/internal/model"
	conf "github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
)

type S3Storage struct {
	BasePath         string
	AbsoluteDestPath string
	Client           *s3.Client
	Bucket           string
}

func NewS3Storage(basePath string, absoluteDestPath string) (*S3Storage, error) {
	endpoint := conf.GetConfig().GetString("S3_ENDPOINT")
	bucket := conf.GetConfig().GetString("S3_BUCKET_NAME")
	accessKey := conf.GetConfig().GetString("S3_ACCESS_KEY")
	secretKey := conf.GetConfig().GetString("S3_SECRET_KEY")
	region := conf.GetConfig().GetString("S3_REGION")
	if endpoint == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("Missing S3 environment variables")
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // REQUIRED for Railway
	})

	return &S3Storage{
		BasePath:         basePath,
		AbsoluteDestPath: absoluteDestPath,
		Client:           client,
		Bucket:           bucket,
	}, nil
}

func (r *S3Storage) Save(c *gin.Context, fileObj *multipart.FileHeader) (*FileMeta, error) {
	originalFileName := filepath.Base(fileObj.Filename)
	extension := strings.Split(originalFileName, ".")[1]
	newName := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	fullPath := path.Join(r.Bucket, r.AbsoluteDestPath, newName)
	log.Info("Uploading file to S3 at path: " + fullPath)
	context := c.Request.Context()
	objectKey := path.Join(r.AbsoluteDestPath, newName)

	file, err := fileObj.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = r.UploadFile(context, objectKey, file, fileObj.Header.Get("Content-Type"), true)
	if err != nil {
		log.Error("Failed to upload file: " + err.Error())
		return nil, err
	}
	fileMeta := &FileMeta{
		OriginalName: originalFileName,
		SavedName:    newName,
		FullPath:     fullPath,
		AbsolutePath: objectKey,
		Extenstion:   extension,
		// DownloadURL:  conf.GetConfig().GetString("S3_ENDPOINT") + "/" + r.Bucket + "/" + objectKey,
		DownloadURL: getDownloadableUrl(objectKey),
	}
	log.Info("File uploaded successfully: " + fileMeta.DownloadURL)
	return fileMeta, nil
}

// UploadFile uploads data to Railway bucket
func (r *S3Storage) UploadFile(ctx context.Context, key string, body io.Reader, contentType string, public bool) error {

	input := &s3.PutObjectInput{
		Bucket:      aws.String(r.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}

	if public {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	_, err := r.Client.PutObject(ctx, input)
	return err
}

func (r *S3Storage) Delete(filePath string) error {
	if filePath != "" {
		return fmt.Errorf("Invalid filepath to delete")
	}
	_, err := r.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(strings.TrimPrefix(filePath, r.Bucket+"/")),
	})
	if err != nil {
		log.Error("Failed to delete file: " + err.Error())
		return err
	}
	log.Info("File deleted successfully: " + filePath)
	return nil
}

func (r *S3Storage) DeleteAll(documents []model.Document) error {
	for _, doc := range documents {
		err := r.Delete(doc.FileName)
		if err != nil {
			log.Warnf("Failed to delete file %s: %v", doc.FileName, err)
		}
	}
	return nil
}

func getDownloadableUrl(objectKey string) string {
	service_host := conf.GetConfig().GetString("SERVICE_HOST")
	if service_host == "" {
		service_host = conf.GetConfig().GetString("S3_ENDPOINT")
	}
	if service_host != "" {
		url := fmt.Sprintf("%s/api/v1/cmn/serve-static/%s", service_host, objectKey)
		return url
	}
	return ""
}

func (r *S3Storage) CreatePreSignedUrl(filePath string) (string, error) {
	presignClient := s3.NewPresignClient(r.Client)
	params := &s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(filePath),
	}
	presignedURL, err := presignClient.PresignGetObject(context.TODO(), params)
	if err != nil {
		log.Error("Failed to create presigned URL: " + err.Error())
		return "", err
	}
	return presignedURL.URL, nil
}

func (r *S3Storage) DownloadFile(c *gin.Context, url string) (*[]byte, string, error) {
	// fetch the file bytes from the presigned url
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Failed to create request for file: %v", err)
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Failed to fetch file from url: %v", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Errorf("Error fetching file, status: %d, body: %s", resp.StatusCode, string(body))
		return nil, "", err
	}

	// read content
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read file content: %v", err)
		return nil, "", err
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	return &data, contentType, nil
}
