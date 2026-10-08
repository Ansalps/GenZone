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

// UploadToS3 uploads a file to S3 and returns the relative S3 object Key.
func ResolveObjectURL(ctx context.Context, bucketEnvKey, key string) (string, error) {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return "", nil
	}
	if strings.HasPrefix(trimmedKey, "http://") || strings.HasPrefix(trimmedKey, "https://") {
		return trimmedKey, nil
	}

	bucketName := resolveBucketName(bucketEnvKey, "S3_BUCKET_NAME", "PRIVATE_S3_BUCKET_NAME", "PRIVATE_SE_BUCKET_NAME", "AWS_BUCKET_NAME")
	if bucketName == "" {
		return "", errors.New("s3 bucket name not configured")
	}

	awsRegion := resolveAWSRegion()
	if awsRegion == "" {
		return "", errors.New("aws region not configured")
	}

	isPrivate := strings.Contains(strings.ToLower(bucketEnvKey), "private") || strings.Contains(strings.ToLower(bucketName), "private")
	if isPrivate {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
		if err != nil {
			return "", fmt.Errorf("unable to load AWS SDK config: %w", err)
		}
		client := s3.NewFromConfig(cfg)
		return GeneratePresignedURL(ctx, client, bucketName, trimmedKey, 24*time.Hour)
	}

	return GetPublicURL(bucketName, awsRegion, trimmedKey), nil
}

func UploadToS3(ctx context.Context, fileHeader *multipart.FileHeader, folderName string, bucketEnvKey string) (string, error) {
	s3BucketName := resolveBucketName(bucketEnvKey, "S3_BUCKET_NAME", "PRIVATE_S3_BUCKET_NAME", "PRIVATE_SE_BUCKET_NAME", "AWS_BUCKET_NAME")
	if s3BucketName == "" {
		return "", errors.New("s3 bucket name not configured")
	}

	awsRegion := resolveAWSRegion()
	if awsRegion == "" {
		return "", errors.New("aws region not configured")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return "", fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// Clean key generation
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	cleanFolder := strings.Trim(folderName, "/")
	key := filename
	if cleanFolder != "" {
		key = path.Join(cleanFolder, filename)
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}

	// ALWAYS return the relative key (e.g. "profile-pictures/uuid.avif") to store in DB
	return key, nil
}

// GetPublicURL returns a static public S3 or CloudFront URL.
func GetPublicURL(bucketName, region, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key)
}

// GeneratePresignedURL generates a temporary read URL for private bucket objects.
func GeneratePresignedURL(ctx context.Context, client *s3.Client, bucketName, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(client, s3.WithPresignExpires(expiration))
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}
	return req.URL, nil
}
