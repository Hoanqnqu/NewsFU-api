package outport

import (
	"news-api/internal/db"

	"github.com/google/uuid"
)

type Category struct {
	ID   uuid.UUID
	Name string
}

type Categories interface {
	GetAll() ([]db.Category, error)
	Insert(category Category) error
	Update(category Category) error
	Search(keyword string) ([]db.Category, error)
	Delete(id uuid.UUID) error
}
