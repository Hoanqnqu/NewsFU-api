package service

import (
	"context"
	"mime/multipart"
	outport "news-api/application/port/out"
)

type UploadService struct {
	s3port outport.StoragePort
}

func NewUploadService(s3port outport.StoragePort) *UploadService {
	return &UploadService{
		s3port: s3port,
	}
}

func (u UploadService) PutObject(ctx context.Context, objectKey string, file multipart.File) (url string, err error) {
	url, err = u.s3port.PutObject(ctx, objectKey, file)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u UploadService) DeleteObject(ctx context.Context, objectKey string) error {
	err := u.s3port.DeleteObject(ctx, objectKey)
	if err != nil {
		return err
	}
	return nil
}
