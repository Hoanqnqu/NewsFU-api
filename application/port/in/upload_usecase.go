package inport

import (
	"context"
	"mime/multipart"
)

type UploadUseCase interface {
	PutObject(ctx context.Context, objectKey string, file multipart.File) (string, error)
	DeleteObject(ctx context.Context, objectKey string) error
}
