package outAdapter

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	outport "news-api/application/port/out"
	"news-api/internal/db"
)

type CommentAdapter struct {
	pool *pgxpool.Pool
}

func NewCommentAdapter(pool *pgxpool.Pool) *CommentAdapter {
	return &CommentAdapter{pool: pool}
}

func (c *CommentAdapter) InsertComment(newsID, userID uuid.UUID, comment string) error {
	query := db.New(c.pool)
	return query.InsertComment(context.Background(), db.InsertCommentParams{
		NewsID: pgtype.UUID{
			Bytes: newsID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
		Text: pgtype.Text{
			String: comment,
			Valid:  true,
		},
	})
}

func (c *CommentAdapter) GetCommentsByNews(newsID uuid.UUID) ([]*outport.Comment, error) {
	query := db.New(c.pool)
	comments, err := query.QueryCommentByNews(context.Background(), pgtype.UUID{
		Bytes: newsID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}
	rs := make([]*outport.Comment, len(comments))
	for i, comment := range comments {
		rs[i] = &outport.Comment{
			Text:        comment.Text.String,
			PublishedAt: comment.PublishedAt.Time,
			UserName:    comment.Name.String,
			UserAvatar:  comment.ImageUrl.String,
		}
	}
	return rs, nil
}
