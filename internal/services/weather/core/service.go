package core

import (
	"APIs/internal/common/entities"
	"APIs/internal/common/models"
	"APIs/internal/services/weather/ports"
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const weather_ttl = 20 // 20 minutes

type Service struct {
	repository    ports.Repository
	weatherClient ports.WeatherClient
}

func NewService(repository ports.Repository, weatherClient ports.WeatherClient) *Service {
	return &Service{
		repository:    repository,
		weatherClient: weatherClient,
	}
}

func (s *Service) ValidateWeather(ctx context.Context, query entities.WeatherQuery) (bool, error) {
	weatherModel, err := s.repository.GetWeather(ctx, query.Town)
	if err != nil {
		log.Info().Str("details", err.Error()).Msg("unable to find Weather in db")
	}

	if weatherModel != nil && weatherModel.UpdatedAt.After(time.Now().Add(-weather_ttl*time.Minute)) { // if weather is found and not too old, reuse it
		return weatherModel.Type == query.Type && weatherModel.Temp > float32(query.TempMin), nil
	}

	// else call openweather API and save result
	weatherResponse, err := s.weatherClient.FetchWeather(ctx, query.Town)
	if err != nil {
		log.Error().Err(err).Msg("error fetching weather from openweather API")
		return false, err
	}

	// upsert weather in db with new values
	weatherToUpsert := &models.Weather{
		Town: query.Town,
		Temp: weatherResponse.Temp,
		Type: weatherResponse.Type,
	}
	if weatherModel != nil {
		weatherToUpsert.ID = weatherModel.ID
		weatherModel, err = s.repository.UpdateWeather(ctx, weatherToUpsert)
		if err != nil {
			log.Error().Err(err).Msg("error updating weather in db")
			return false, err
		}
	} else {
		weatherModel, err = s.repository.CreateWeather(ctx, weatherToUpsert)
		if err != nil {
			log.Error().Err(err).Msg("error upserting weather in db")
			return false, err
		}
	}

	return weatherModel.Type == query.Type && weatherModel.Temp > float32(query.TempMin), nil
}
