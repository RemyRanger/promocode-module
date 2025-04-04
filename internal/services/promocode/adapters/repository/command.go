package repository

import (
	"APIs/internal/common/models"
	"context"
)

func (r *Repository) CreatePromocode(ctx context.Context, record *models.Promocode) (*models.Promocode, error) {
	if err := r.BeginWithCtx(ctx).Create(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
