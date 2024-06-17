package outAdapter

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/zhenghaoz/gorse/client"
	"time"
)

type GorseAdapter struct {
	gorse *client.GorseClient
}

func NewGorseAdapter(gorse *client.GorseClient) *GorseAdapter {
	return &GorseAdapter{gorse: gorse}
}

func (g *GorseAdapter) InsertUser(ctx context.Context, id uuid.UUID) error {
	rowAffected, err := g.gorse.InsertUser(ctx, client.User{
		UserId: id.String(),
	})
	if err != nil {
		return err
	}
	if rowAffected.RowAffected != 1 {
		return fmt.Errorf("can not insert user")
	}
	return nil
}

func (g *GorseAdapter) InsertNews(ctx context.Context, id uuid.UUID, categories []uuid.UUID) error {
	rowAffected, err := g.gorse.InsertItem(ctx, client.Item{
		ItemId:   id.String(),
		IsHidden: false,
		Categories: func() []string {
			result := make([]string, len(categories))
			for i, v := range categories {
				result[i] = v.String()
			}
			return result
		}(),
		Timestamp: time.Now().String(),
	})
	if err != nil {
		return err
	}
	if rowAffected.RowAffected != 1 {
		return fmt.Errorf("can not insert news")
	}
	return nil
}

func (g *GorseAdapter) SendLike(ctx context.Context, userID string, newsID string) error {
	rowAffected, err := g.gorse.InsertFeedback(ctx, []client.Feedback{
		{
			FeedbackType: "like",
			UserId:       userID,
			ItemId:       newsID,
			Timestamp:    time.Now().String(),
		},
	})
	if err != nil {
		return err
	}
	if rowAffected.RowAffected != 1 {
		return fmt.Errorf("can not send like feedback")
	}
	return nil
}

func (g *GorseAdapter) SendSave(ctx context.Context, userID string, newsID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GorseAdapter) SendDislike(ctx context.Context, userID string, newsID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GorseAdapter) GetLatestNews(ctx context.Context, count int, offset int) ([]string, error) {
	scores, err := g.gorse.GetItemLatest(ctx, "", count, offset)
	if err != nil {
		return nil, err
	}
	rs := make([]string, len(scores))
	for i, score := range scores {
		rs[i] = score.Id
	}
	return rs, nil
}

func (g *GorseAdapter) GetPopularByCategory(ctx context.Context, categoryID string, count int, offset int) ([]string, error) {
	var scores []client.Score
	var err error
	if categoryID == "" {
		scores, err = g.gorse.GetItemPopular(ctx, "", count, offset)
	}
	scores, err = g.gorse.GetItemPopularWithCategory(ctx, "", categoryID, count, offset)
	if err != nil {
		return nil, err
	}
	rs := make([]string, len(scores))
	for i, score := range scores {
		rs[i] = score.Id
	}
	return rs, nil
}

func (g *GorseAdapter) GetRecommendForUser(ctx context.Context, userID string, count int, offset int) ([]string, error) {
	return g.gorse.GetItemRecommendWithCategory(context.Background(), userID, "", "read", "300s", count, offset)
}

func (g *GorseAdapter) DeleteUser(ctx context.Context, userID string) error {
	rowAffected, err := g.gorse.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	if rowAffected.RowAffected != 1 {
		return fmt.Errorf("can not delete user")
	}
	return nil
}

func (g *GorseAdapter) DeleteNews(ctx context.Context, newsID string) error {
	rowAffected, err := g.gorse.DeleteItem(ctx, newsID)
	if err != nil {
		return err
	}
	if rowAffected.RowAffected != 1 {
		return fmt.Errorf("can not delete news")
	}
	return nil
}
