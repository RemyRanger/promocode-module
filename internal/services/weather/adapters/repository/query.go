package repository

import (
	"APIs/internal/common/models"
	"context"

	"github.com/pkg/errors"
)

func (r *Repository) GetWeather(ctx context.Context, town string) (*models.Weather, error) {
	model := &models.Weather{}
	if err := r.BeginWithCtx(ctx).Where("town = ?", town).First(model).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return model, nil
}
