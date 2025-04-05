package ports

import (
	"context"

	"APIs/internal/common/models"
)

type Repository interface {
	GetWeather(ctx context.Context, town string) (*models.Weather, error)

	CreateWeather(ctx context.Context, model *models.Weather) (*models.Weather, error)

	UpdateWeather(ctx context.Context, model *models.Weather) (*models.Weather, error)
}
