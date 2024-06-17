package service

import (
	inport "news-api/application/port/in"
	outport "news-api/application/port/out"
	"news-api/internal/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func MapUser(user db.User) *inport.User {
	return &inport.User{
		ID:        ToUUID(user.ID).String(),
		AuthID:    user.AuthID,
		Email:     user.Email.String,
		Name:      user.Name.String,
		Role:      user.Role.String,
		ImageUrl:  user.ImageUrl.String,
		CreatedAt: PGTimestampToTime(user.CreatedAt),
		UpdatedAt: PGTimestampToTime(user.UpdatedAt),
		DeletedAt: PGTimestampToTime(user.DeletedAt),
	}
}
func MapNews(news outport.NewsWithCategory) *inport.News {
	return &inport.News{
		ID:          ToUUID(news.ID).String(),
		Author:      news.Author.String,
		Title:       news.Title.String,
		Content:     news.Content.String,
		Description: news.Description.String,
		URL:         news.Url.String,
		ImageURL:    news.ImageUrl.String,
		PublishAt:   news.PublishAt.Time,
		Categories: func() []uuid.UUID {
			if len(news.Categories) == 0 {
				return nil
			}
			ids := make([]uuid.UUID, len(news.Categories))
			for i, v := range news.Categories {
				ids[i] = v.Bytes
			}
			return ids
		}(),
		CreatedAt: PGTimestampToTime(news.CreatedAt),
		UpdatedAt: PGTimestampToTime(news.UpdatedAt),
		DeletedAt: PGTimestampToTime(news.DeletedAt),
		View:      news.View,
	}
}

func MapCategory(category db.Category) *inport.Category {
	return &inport.Category{
		ID:        ToUUID(category.ID).String(),
		Name:      category.Name.String,
		CreatedAt: PGTimestampToTime(category.CreatedAt),
		UpdatedAt: PGTimestampToTime(category.UpdatedAt),
		DeletedAt: PGTimestampToTime(category.DeletedAt),
	}
}

func ToUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return id.Bytes
}

func PGTimestampToTime(timestamp pgtype.Timestamp) *time.Time {
	if !timestamp.Valid {
		return nil
	}
	return &timestamp.Time
}
