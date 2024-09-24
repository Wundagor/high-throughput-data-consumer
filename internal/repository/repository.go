package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type SourceData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
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
