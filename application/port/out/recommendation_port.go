package outport

import (
	"context"
	"github.com/google/uuid"
)

type RecommendationSystem interface {
	InsertUser(ctx context.Context, id uuid.UUID) error
	InsertNews(ctx context.Context, id uuid.UUID, categories []uuid.UUID) error
	SendLike(ctx context.Context, userID string, newsID string) error
	SendSave(ctx context.Context, userID string, newsID string) error
	SendDislike(ctx context.Context, userID string, newsID string) error
	GetLatestNews(ctx context.Context, count int, offset int) ([]string, error)
	GetPopularByCategory(ctx context.Context, categoryID string, count int, offset int) ([]string, error)
	GetRecommendForUser(ctx context.Context, userID string, count int, offset int) ([]string, error)
	DeleteUser(ctx context.Context, userID string) error
	DeleteNews(ctx context.Context, newsID string) error
}
