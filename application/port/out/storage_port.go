package outport

import (
	"context"
	"mime/multipart"
)

type StoragePort interface {
	PutObject(ctx context.Context, objectKey string, file multipart.File) (url string, err error)
	DeleteObject(ctx context.Context, objectKey string) error
}
