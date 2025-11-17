package s3

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	literals "recorder-service/internal/config"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service interface {
	PutObject(ctx context.Context, pathToFiles string) ([]string, error)
}

type s3Service struct {
	bucketName string
	uploader   *manager.Uploader
	downloader *manager.Downloader
}

// GetObject implements S3Service.
func (s *s3Service) GetObject(ctx context.Context, key string) (string, error) {
	return "", nil
}

// PutObject implements S3Service.
func (s *s3Service) PutObject(ctx context.Context, pathToFiles string) ([]string, error) {
	entries, err := os.ReadDir(pathToFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %s", pathToFiles)
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(entries))
	keys := make(chan string, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			wg.Go(func() {
				dest := path.Join(pathToFiles, entry.Name())

				file, err := os.Open(dest)
				if err != nil {
					errors <- fmt.Errorf("failed to open file: %s", dest)
					return
				}
				defer file.Close()

				result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
					Bucket: aws.String(s.bucketName),
					Key:    aws.String(dest),
					Body:   file,
				})

				if err != nil {
					errors <- err
					return
				}

				keys <- dest

				log.Printf("Successfully uploaded: %s in location: %s", dest, result.Location)
			})
		}
	}

	wg.Wait()
	close(errors)
	close(keys)

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	res := []string{}
	for url := range keys {
		res = append(res, url)
	}

	return res, nil
}

func NewS3Service(ctx context.Context, bucketName string) (S3Service, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("no bucket name")
	}

	var cfg aws.Config
	var err error

	cfg, err = config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv(literals.AWS_ACCESS_KEY_ID),
			os.Getenv(literals.AWS_SECRET_ACCESS_KEY),
			os.Getenv(literals.AWS_SESSION_TOKEN),
		)),
		config.WithRegion(os.Getenv(literals.AWS_REGION)),
	)

	if err != nil {
		log.Printf("Failed to set config: %v", err)
		return nil, fmt.Errorf("failed to set config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)

	return &s3Service{uploader: uploader, downloader: downloader, bucketName: bucketName}, nil
}
