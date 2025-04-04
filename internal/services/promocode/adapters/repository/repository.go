package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) BeginWithCtx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *Repository) BeginTransactionWithCtx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.Begin().WithContext(ctx)
	// check transaction started correctly
	if err := tx.Error; err != nil {
		return nil, err
	}
	return tx, nil
}
