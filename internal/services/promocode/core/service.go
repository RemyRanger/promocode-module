package core

import (
	"APIs/internal/common/models"
	"APIs/internal/services/promocode/ports"
	"context"
)

type Service struct {
	repository ports.Repository
}

func New(repository ports.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) SavePromocode(ctx context.Context, model *models.Promocode) (*models.Promocode, error) {
	return s.repository.CreatePromocode(ctx, model)
}

func (s *Service) ValidatePromocode(ctx context.Context, promocodeName string, age int64, town string) ([]string, error) {
	promocode, err := s.repository.GetPromocode(ctx, promocodeName)
	if err != nil {
		return nil, err
	}

	promocodeDto, err := ports.ModelToDto(promocode)
	if err != nil {
		return nil, err
	}

	restrictionValidator := NewRestrictionsValidator(age, town)

	return restrictionValidator.validateRestrictions(promocodeDto.Restrictions), nil
}
