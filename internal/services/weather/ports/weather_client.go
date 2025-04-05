package ports

import (
	"APIs/internal/common/entities"
	"context"
)

type WeatherClient interface {
	FetchWeather(ctx context.Context, town string) (*entities.Weather, error)
}
