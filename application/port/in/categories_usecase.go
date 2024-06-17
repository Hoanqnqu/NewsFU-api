package inport

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateCategoryPayload struct {
	Name string `json:"name"`
}

type UpdateCategoryPayload struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CategoriesUseCase interface {
	GetAll() ([]*Category, error)
	Insert(category *CreateCategoryPayload) error
	Update(category *UpdateCategoryPayload) error
	Search(keyword string) ([]*Category, error)
	Delete(id string) error
}
