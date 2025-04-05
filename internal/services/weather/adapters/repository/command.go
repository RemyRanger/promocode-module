package repository

import (
	"APIs/internal/common/models"
	"context"
)

func (r *Repository) CreateWeather(ctx context.Context, record *models.Weather) (*models.Weather, error) {
	if err := r.BeginWithCtx(ctx).Create(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (r *Repository) UpdateWeather(ctx context.Context, record *models.Weather) (*models.Weather, error) {
	if err := r.BeginWithCtx(ctx).Model(record).Updates(record).Take(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
