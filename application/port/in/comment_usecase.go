package inport

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Text        string    `json:"text"`
	PublishedAt time.Time `json:"published_at"`
	UserName    string    `json:"user_name"`
	UserAvatar  string    `json:"user_avatar"`
}

type CommentUseCase interface {
	InsertComment(newsID, userID uuid.UUID, comment string) error
	GetCommentsByNewsID(newsID uuid.UUID) ([]*Comment, error)
}
