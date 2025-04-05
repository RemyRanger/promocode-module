package ports

import (
	"APIs/internal/common/entities"
	"context"
)

type Service interface {
	ValidateWeather(ctx context.Context, query entities.WeatherQuery) (bool, error)
}
