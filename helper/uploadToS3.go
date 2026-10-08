package helper

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// const (
// 	S3BucketName = "genzone-public-assets" // Replace or load from os.Getenv("S3_BUCKET_NAME")
// 	AWSRegion    = "ap-south-1"          // Replace with your bucket region
// )

func resolveBucketName(keys ...string) string {
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}
	return ""
}

func resolveAWSRegion() string {
	for _, key := range []string{"AWS_REGION", "AWS_REGION_NAME", "AWS_DEFAULT_REGION", "AWSRegion"} {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}
	return ""
}

func UploadToS3(fileHeader *multipart.FileHeader, folderName string, bucketEnvKey string) (string, error) {
	s3BucketName := resolveBucketName(bucketEnvKey, "S3_BUCKET_NAME", "PRIVATE_S3_BUCKET_NAME", "PRIVATE_SE_BUCKET_NAME", "AWS_BUCKET_NAME")
	if s3BucketName == "" {
		return "", errors.New("failed to fetch s3BucketName")
	}

	aWSRegion := resolveAWSRegion()
	if aWSRegion == "" {
		return "", errors.New("failed to fetch awsRegiion")
	}

	// 1. Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// 2. Load AWS SDK config (reads AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY automatically from environment)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(aWSRegion))
	if err != nil {
		return "", fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// 3. Construct a clean S3 key using path.Join
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)

	// Safety: Trim leading/trailing slashes and join properly
	cleanFolder := strings.Trim(folderName, "/")
	var uniqueFileName string
	if cleanFolder != "" {
		uniqueFileName = path.Join(cleanFolder, filename)
	} else {
		uniqueFileName = filename
	}

	// 4. Upload file to S3
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(uniqueFileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}

	// 5. Return a presigned URL so private buckets can still be displayed in the frontend.
	presignClient := s3.NewPresignClient(client, s3.WithPresignExpires(24*time.Hour))
	presignedRequest, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(uniqueFileName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url for uploaded object: %w", err)
	}

	return presignedRequest.URL, nil
}
