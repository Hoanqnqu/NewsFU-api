package outport

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type News struct {
	ID          uuid.UUID
	Title       string
	Content     string
	Description string
	Author      string
	Url         string
	ImageUrl    string
	PublishAt   time.Time
	Categories  []uuid.UUID
}
type NewsWithCategory struct {
	ID          pgtype.UUID
	Author      pgtype.Text
	Title       pgtype.Text
	Description pgtype.Text
	Content     pgtype.Text
	Url         pgtype.Text
	ImageUrl    pgtype.Text
	Categories  []pgtype.UUID
	PublishAt   pgtype.Timestamp
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	DeletedAt   pgtype.Timestamp
	View        int
}

type NewsPort interface {
	GetAll() ([]NewsWithCategory, error)
	Insert(news News) error
	Update(news News) error
	GetNewsByID(newsID, userID string) (*NewsWithCategory, bool, bool, error)
	SearchNews(keyword string) ([]NewsWithCategory, error)
	GetNewsByIDs(ids []string) ([]NewsWithCategory, error)
	Delete(id string) error
}
