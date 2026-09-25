package util

import (
	"context"
	"io"
	"log"
	"os"

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

func UploadS3() {
	ctx := context.Background()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Unable to open file %q: %v", filePath, err)
		return
	}
	defer file.Close()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
		Body:   file,
	})

	if err != nil {
		log.Printf("Unable to upload %q to %q: %v", filePath, bucket, err)
		return
	}

	log.Println("Successfully uploaded!")
}

func DownloadS3(obj string) {
	ctx := context.Background()

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(obj),
	})
	if err != nil {
		log.Printf("Unable to download %q: %v", obj, err)
		return
	}
	defer result.Body.Close()

	file, err := os.Create(dir + obj)
	if err != nil {
		log.Printf("Unable to create file: %v", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		log.Printf("Unable to save downloaded file: %v", err)
		return
	}

	log.Printf("Successfully downloaded %q", obj)
}
