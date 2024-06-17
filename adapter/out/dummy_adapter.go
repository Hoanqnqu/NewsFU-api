package outAdapter

import "github.com/jackc/pgx/v5/pgxpool"

type DummyAdapter struct {
	pool *pgxpool.Pool
}

func NewDummyAdapter(pool *pgxpool.Pool) *DummyAdapter {
	return &DummyAdapter{pool: pool}
}

func (d *DummyAdapter) GetDummy() string {
	return "Hello Hoang Ngu"
}
