package database

import (
	"database/sql"

	"github.com/Wundagor/high-throughput-data-consumer/internal/config"
	"github.com/Wundagor/high-throughput-data-consumer/models"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func Connect(config config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Save(data *models.SourceData) error {
	_, err := db.Exec(
		"INSERT INTO destination_data (id, name, description, created_at) VALUES (?, ?, ?, ?)",
		data.ID, data.Name, data.Description, data.CreatedAt,
	)

	return err
}
