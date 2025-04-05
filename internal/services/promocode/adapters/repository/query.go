package repository

import (
	"APIs/internal/common/models"
	"context"

	"github.com/pkg/errors"
)

func (r *Repository) GetPromocode(ctx context.Context, promocodeName string) (*models.Promocode, error) {
	model := &models.Promocode{}
	if err := r.BeginWithCtx(ctx).Where("name = ?", promocodeName).First(model).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return model, nil
}
