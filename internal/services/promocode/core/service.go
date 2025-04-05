package core

import (
	"APIs/internal/common/entities/custom_errors"
	"APIs/internal/common/models"
	"APIs/internal/services/promocode/ports"
	weather_ports "APIs/internal/services/weather/ports"
	"context"
)

type Service struct {
	repository     ports.Repository
	weatherService weather_ports.Service
}

func NewService(repository ports.Repository, weatherService weather_ports.Service) *Service {
	return &Service{
		repository:     repository,
		weatherService: weatherService,
	}
}

func (s *Service) SavePromocode(ctx context.Context, model *models.Promocode) (*models.Promocode, error) {
	_, err := s.repository.GetPromocode(ctx, model.Name)
	if err == nil {
		return nil, custom_errors.New(custom_errors.ErrPromocodeExist)
	}

	return s.repository.CreatePromocode(ctx, model)
}

func (s *Service) ValidatePromocode(ctx context.Context, promocodeName string, age int64, town string) ([]string, error) {
	// Fetch promocode from DB
	promocode, err := s.repository.GetPromocode(ctx, promocodeName)
	if err != nil {
		return nil, err
	}

	// Map promocode to dto for processing on the DTO json struct
	promocodeDto, err := ports.ModelToDto(promocode)
	if err != nil {
		return nil, err
	}

	// Validate Promocode restriction according to the inputs
	restrictionValidator := NewRestrictionsValidator(ctx, s.weatherService, age, town)
	reasons, err := restrictionValidator.validateRestrictions(promocodeDto.Restrictions)
	if err != nil {
		return nil, err
	}

	return reasons, nil
}
