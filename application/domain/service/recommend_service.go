package service

import (
	"context"
	"github.com/google/uuid"
	outport "news-api/application/port/out"
)

type RecommendService struct {
	recommendSystem outport.RecommendationSystem
}

func NewRecommendService(recommendSystem outport.RecommendationSystem) *RecommendService {
	return &RecommendService{recommendSystem: recommendSystem}
}

func (r *RecommendService) InsertUser(ctx context.Context, id uuid.UUID) error {
	return r.recommendSystem.InsertUser(ctx, id)
}
