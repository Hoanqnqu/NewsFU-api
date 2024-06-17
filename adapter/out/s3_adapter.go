package outAdapter

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
)

type S3Adapter struct {
	s3Client   *s3.Client
	bucketName string
	region     string
}

func NewS3Adapter(accessKey string, secretKey string, s3BucketRegion string, bucketName string) *S3Adapter {
	options := s3.Options{
		Region:      s3BucketRegion,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	}

	client := s3.New(options, func(o *s3.Options) {
		o.Region = s3BucketRegion
		o.UseAccelerate = false
	})

	return &S3Adapter{
		s3Client:   client,
		bucketName: bucketName,
		region:     s3BucketRegion,
	}
}

func (repo *S3Adapter) PutObject(ctx context.Context, objectKey string, file multipart.File) (url string, err error) {
	_, err = repo.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(repo.bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", repo.bucketName, repo.region, objectKey)
	return url, nil
}

func (repo *S3Adapter) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := repo.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(repo.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}
	return nil
}
