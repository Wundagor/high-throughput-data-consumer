package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SourceData struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type Repository interface {
	InsertData(ctx context.Context, data SourceData) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) InsertData(ctx context.Context, data SourceData) error {
	query := `INSERT INTO destination_data (name, description, created_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, data.Name, data.Description, data.CreatedAt)

	return err
}
