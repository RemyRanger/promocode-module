package ports

import (
	"APIs/internal/common/models"
	"context"
)

type Service interface {
	SavePromocode(ctx context.Context, model *models.Promocode) (*models.Promocode, error)

	ValidatePromocode(ctx context.Context, promocodeName string, age int64, town string) ([]string, error)
}
