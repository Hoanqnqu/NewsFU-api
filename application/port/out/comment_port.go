package outport

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Text        string
	PublishedAt time.Time
	UserName    string
	UserAvatar  string
}

type CommentPort interface {
	InsertComment(newsID, userID uuid.UUID, comment string) error
	GetCommentsByNews(newsID uuid.UUID) ([]*Comment, error)
}
