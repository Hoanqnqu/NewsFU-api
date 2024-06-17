package outAdapter

import (
	"context"
	"github.com/google/uuid"
	outport "news-api/application/port/out"
	db "news-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryAdapter struct {
	pool *pgxpool.Pool
}

func NewCategoryAdapter(pool *pgxpool.Pool) *CategoryAdapter {
	return &CategoryAdapter{pool: pool}
}

func (u *CategoryAdapter) GetAll() ([]db.Category, error) {
	query := db.New(u.pool)

	return query.GetAllCategories(context.Background())

}

func (u *CategoryAdapter) Insert(category outport.Category) error {
	query := db.New(u.pool)

	err := query.InsertCategory(context.Background(), db.InsertCategoryParams{
		ID: pgtype.UUID{
			Bytes: category.ID,
			Valid: true,
		},
		Name: pgtype.Text{
			String: category.Name,
			Valid:  true,
		},
	})
	return err
}

func (u *CategoryAdapter) Search(keyword string) ([]db.Category, error) {
	query := db.New(u.pool)
	return query.SearchCategories(context.Background(), pgtype.Text{
		String: keyword,
		Valid:  true,
	})
}

func (u *CategoryAdapter) Update(category outport.Category) error {
	query := db.New(u.pool)

	err := query.UpdateCategory(context.Background(), db.UpdateCategoryParams{
		ID: pgtype.UUID{
			Bytes: category.ID,
			Valid: true,
		},
		Name: pgtype.Text{
			String: category.Name,
			Valid:  true,
		},
	})
	return err
}

func (u *CategoryAdapter) Delete(id uuid.UUID) error {
	return nil
}
