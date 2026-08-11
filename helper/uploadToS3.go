package helper

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"
)

// const (
// 	S3BucketName = "genzone-public-assets" // Replace or load from os.Getenv("S3_BUCKET_NAME")
// 	AWSRegion    = "ap-south-1"          // Replace with your bucket region
// )

func UploadToS3(fileHeader *multipart.FileHeader) (string, error) {
	s3BucketName := os.Getenv("S3_BUCKET_NAME")
	fmt.Println("s3BucketName")
	aWSRegion := os.Getenv("AWSRegion")
	fmt.Println("", s3BucketName)
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
	fmt.Println("", cfg.AccountIDEndpointMode, cfg.AppID)
	cred, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		fmt.Println("pring", err)
	}
	fmt.Println("cred", cred.AccessKeyID, cred.AccountID, cred.SecretAccessKey, cred.CanExpire, cred.Expires, cred.Expired())

	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Println("STS Error:", err)
	} else {
		fmt.Println("CURRENT IAM ARN:", *identity.Arn)
	}

	client := s3.NewFromConfig(cfg)

	// 3. Generate a unique filename using UUID to prevent overwrite collisions
	ext := filepath.Ext(fileHeader.Filename)
	uniqueFileName := fmt.Sprintf("categories/%s%s", uuid.New().String(), ext)

	// 4. Upload file to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(uniqueFileName),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}

	// 5. Return the public S3 URL
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3BucketName, aWSRegion, uniqueFileName)
	return s3URL, nil
}
