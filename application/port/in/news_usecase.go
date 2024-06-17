package inport

import (
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID          string      `json:"id"`
	Author      string      `json:"author"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Description string      `json:"description"`
	URL         string      `json:"url"`
	ImageURL    string      `json:"image_url"`
	PublishAt   time.Time   `json:"publish_at"`
	Categories  []uuid.UUID `json:"categories"`
	CreatedAt   *time.Time  `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`
	DeletedAt   *time.Time  `json:"deleted_at"`
	IsLiked     bool        `json:"isLiked"`
	IsDisliked  bool        `json:"isDisliked"`
	View        int         `json:"view"`
	RelatedNews []*News     `json:"related_news"`
}
type CreateNewsPayload struct {
	Author      string      `json:"author"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Description string      `json:"description"`
	URL         string      `json:"url"`
	ImageURL    string      `json:"image_url"`
	PublishAt   time.Time   `json:"publish_at"`
	Categories  []uuid.UUID `json:"categories"`
}

type UpdateNewsPayload struct {
	ID          uuid.UUID   `json:"id"`
	Author      string      `json:"author"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Description string      `json:"description"`
	URL         string      `json:"url"`
	ImageURL    string      `json:"image_url"`
	PublishAt   time.Time   `json:"publish_at"`
	Categories  []uuid.UUID `json:"categories"`
}

type NewsUseCase interface {
	GetAll() ([]*News, error)
	Insert(news *CreateNewsPayload) error
	Update(news *UpdateNewsPayload) error
	GetNewsByID(newsID, userID string) (*News, error)
	SearchNews(keyword string) ([]*News, error)
	GetLatestNews(count int, offset int) ([]*News, error)
	GetPopular(categoryID string, count int, offset int) ([]*News, error)
	GetRecommend(userID string, count int, offset int) ([]*News, error)
	Delete(id string) error
}
