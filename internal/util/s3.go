package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucket   = "s3-image-processing-saker233"
	filePath = "image/image.jpg"
	dir      = "image/"
)

var s3Client *s3.Client

func InitS3() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	s3Client = s3.NewFromConfig(cfg)

	log.Println("S3 client initialized")
}

func UploadS3() (string, error) {
	ctx := context.Background()

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %w", filePath, err)
	}
	defer file.Close()

	// Generate a unique S3 key
	key := "images/" + uuid.New().String() + ".jpg"

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return "", fmt.Errorf("unable to upload %q to %q: %w", filePath, bucket, err)
	}

	log.Printf("Successfully uploaded with key: %s", key)

	return key, nil
}

func DownloadS3(obj string) error {
	ctx := context.Background()

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(obj),
	})
	if err != nil {
		return fmt.Errorf("unable to download %q: %w", obj, err)
	}
	defer result.Body.Close()

	filename := filepath.Base(obj)
	localPath := filepath.Join(dir, filename)

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("unable to create file %q: %w", localPath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		return fmt.Errorf("unable to save downloaded file: %w", err)
	}

	log.Printf("Successfully downloaded %q to %q", obj, localPath)

	return nil
}
