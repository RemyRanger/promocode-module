package ports

import (
	"context"

	"APIs/internal/common/models"
)

type Repository interface {
	GetPromocode(ctx context.Context, promocodeName string) (*models.Promocode, error)

	CreatePromocode(ctx context.Context, model *models.Promocode) (*models.Promocode, error)
}
