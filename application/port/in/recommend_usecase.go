package inport

import (
	"context"
	"github.com/google/uuid"
)

type RecommendUseCase interface {
	InsertUser(ctx context.Context, id uuid.UUID) error
}
