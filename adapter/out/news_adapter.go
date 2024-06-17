package outAdapter

import (
	"context"
	"encoding/json"
	"fmt"
	outport "news-api/application/port/out"
	db "news-api/internal/db"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsAdapter struct {
	pool *pgxpool.Pool
}

func NewNewsAdapter(pool *pgxpool.Pool) *NewsAdapter {
	return &NewsAdapter{pool: pool}
}

func (u *NewsAdapter) GetAll() ([]outport.NewsWithCategory, error) {
	query := db.New(u.pool)
	news, err := query.GetAllNews(context.Background())
	sl := make([]outport.NewsWithCategory, len(news))
	if err != nil {
		return nil, err
	}
	for i, v := range news {
		var categoryIds []pgtype.UUID
		sl[i].Author = v.Author
		sl[i].Content = v.Content
		sl[i].Description = v.Description
		sl[i].Title = v.Title
		sl[i].Url = v.Url
		sl[i].ImageUrl = v.ImageUrl
		sl[i].PublishAt = v.PublishAt
		sl[i].ID = v.ID
		err = json.Unmarshal(v.CategoryIds, &categoryIds)
		fmt.Println("CategoryIds:", categoryIds)
		if err != nil {
			return nil, err
		}
		sl[i].Categories = categoryIds
		sl[i].View = convertViewType(v.ViewCount)
	}
	return sl, nil
}

func (u *NewsAdapter) SearchNews(keyword string) ([]outport.NewsWithCategory, error) {
	query := db.New(u.pool)
	news, err := query.SearchNews(context.Background(), pgtype.Text{
		String: keyword,
		Valid:  true,
	})
	sl := make([]outport.NewsWithCategory, len(news))
	if err != nil {
		return nil, err
	}
	for i, v := range news {
		var categoryIds []pgtype.UUID
		sl[i].Author = v.Author
		sl[i].Content = v.Content
		sl[i].Description = v.Description
		sl[i].Title = v.Title
		sl[i].Url = v.Url
		sl[i].ImageUrl = v.ImageUrl
		sl[i].PublishAt = v.PublishAt
		sl[i].ID = v.ID
		err = json.Unmarshal(v.CategoryIds, &categoryIds)
		fmt.Println("CategoryIds:", categoryIds)
		if err != nil {
			return nil, err
		}
		sl[i].Categories = categoryIds
		sl[i].View = convertViewType(v.ViewCount)
	}
	return sl, nil
}

func (u *NewsAdapter) Insert(news outport.News) error {
	query := db.New(u.pool)
	err := query.InsertNews(context.Background(), db.InsertNewsParams{
		ID: pgtype.UUID{
			Bytes: news.ID,
			Valid: true,
		},
		Author: pgtype.Text{
			String: news.Author,
			Valid:  true,
		},
		Title: pgtype.Text{
			String: news.Title,
			Valid:  true,
		},
		Content: pgtype.Text{
			String: news.Content,
			Valid:  true,
		},
		Description: pgtype.Text{
			String: news.Description,
			Valid:  true,
		},
		Url: pgtype.Text{
			String: news.Url,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: news.ImageUrl,
			Valid:  true,
		},
		PublishAt: pgtype.Timestamp{
			Time:  news.PublishAt,
			Valid: true,
		},
	})

	if err == nil {
		for _, v := range news.Categories {
			err = query.InsertHasCategory(context.Background(), db.InsertHasCategoryParams{
				NewsID: pgtype.UUID{
					Bytes: news.ID,
					Valid: true,
				},
				CategoryID: pgtype.UUID{
					Bytes: v,
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (u *NewsAdapter) Update(news outport.News) error {
	query := db.New(u.pool)

	err := query.UpdateNews(context.Background(), db.UpdateNewsParams{
		ID: pgtype.UUID{
			Bytes: news.ID,
			Valid: true,
		},
		Author: pgtype.Text{
			String: news.Author,
			Valid:  true,
		},
		Title: pgtype.Text{
			String: news.Title,
			Valid:  true,
		},
		Content: pgtype.Text{
			String: news.Content,
			Valid:  true,
		},
		Description: pgtype.Text{
			String: news.Description,
			Valid:  true,
		},
		Url: pgtype.Text{
			String: news.Url,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: news.ImageUrl,
			Valid:  true,
		},
		PublishAt: pgtype.Timestamp{
			Time:  news.PublishAt,
			Valid: true,
		},
	})
	if err == nil {
		err = query.DeleteHasCategory(context.Background(), pgtype.UUID{
			Bytes: news.ID,
			Valid: true,
		})
		if err == nil {
			for _, v := range news.Categories {
				err = query.InsertHasCategory(context.Background(), db.InsertHasCategoryParams{
					NewsID: pgtype.UUID{
						Bytes: news.ID,
						Valid: true,
					},
					CategoryID: pgtype.UUID{
						Bytes: v,
						Valid: true,
					},
				})
			}
		}

	}
	return err
}

// GetNewsByID retrieves a news item by its ID, including whether the user has liked or disliked it.
// It returns the news item, whether the user has liked it, and whether the user has disliked it, or an error if there was a problem retrieving the data.
func (u *NewsAdapter) GetNewsByID(newsID string, userID string) (news *outport.NewsWithCategory, isLiked bool, isDisliked bool, err error) {
	// Create a new database query object
	query := db.New(u.pool)
	// Get the news item from the database
	_news, err := query.GetNews(context.Background(), pgtype.UUID{
		Bytes: uuid.MustParse(newsID),
		Valid: true,
	})

	// If there was an error retrieving the news item, return the error
	if err != nil {
		return nil, false, false, err
	}

	// Unmarshal the category IDs from the database into a slice of UUID objects
	var category_ids []pgtype.UUID
	err = json.Unmarshal(_news.CategoryIds, &category_ids)
	fmt.Println("CategoryIds:", category_ids)

	if err != nil {
		return nil, false, false, err
	}

	// Populate the news object with the retrieved data
	news = &outport.NewsWithCategory{
		Author:      _news.Author,
		Content:     _news.Content,
		Description: _news.Description,
		Title:       _news.Title,
		Url:         _news.Url,
		ImageUrl:    _news.ImageUrl,
		PublishAt:   _news.PublishAt,
		ID:          _news.ID,
		Categories:  category_ids,
		View:        convertViewType(_news.ViewCount),
	}
	// Check whether the user has liked the news item
	_, _err := query.GetLike(context.Background(), db.GetLikeParams{
		NewsID: pgtype.UUID{
			Bytes: uuid.MustParse(newsID),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: uuid.MustParse(userID),
			Valid: true,
		},
	})
	if _err == nil {
		isLiked = true
	}

	// Check whether the user has disliked the news item
	_, _err = query.GetDislike(context.Background(), db.GetDislikeParams{
		NewsID: pgtype.UUID{
			Bytes: uuid.MustParse(newsID),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: uuid.MustParse(userID),
			Valid: true,
		},
	})
	if _err == nil {
		isDisliked = true
	}
	return
}

func (u *NewsAdapter) GetNewsByIDs(ids []string) ([]outport.NewsWithCategory, error) {
	query := db.New(u.pool)
	pgIDs := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		pgIDs[i] = pgtype.UUID{
			Bytes: uuid.MustParse(id),
			Valid: true,
		}
	}
	news, err := query.GetNewsByIds(context.Background(), pgIDs)
	sl := make([]outport.NewsWithCategory, len(news))
	if err != nil {
		return nil, err
	}
	for i, v := range news {
		sl[i].Author = v.Author
		sl[i].Content = v.Content
		sl[i].Description = v.Description
		sl[i].Title = v.Title
		sl[i].Url = v.Url
		sl[i].ImageUrl = v.ImageUrl
		sl[i].PublishAt = v.PublishAt
		sl[i].ID = v.ID
		sl[i].View = convertViewType(v.ViewCount)
	}
	return sl, nil
}

func (u *NewsAdapter) Delete(id string) error {
	query := db.New(u.pool)
	return query.DeleteNews(context.Background(), pgtype.UUID{
		Bytes: uuid.MustParse(id),
		Valid: true,
	})
}

func convertViewType(v interface{}) int {
	switch v := v.(type) {
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	case uint64:
		return int(v)
	case uint32:
		return int(v)
	case uint:
		return int(v)
	default:
		return 0
	}
}
