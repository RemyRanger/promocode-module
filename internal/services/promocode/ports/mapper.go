package ports

import (
	"APIs/internal/common/models"
	"APIs/internal/common/utils"
	"encoding/json"
)

func ModelToDto(model *models.Promocode) (*Promocode, error) {
	var dto Promocode

	id := (Id)(model.ID)
	dto.Id = &id
	dto.Name = model.Name
	dto.CreatedAt = model.CreatedAt
	dto.UpdatedAt = model.UpdatedAt
	dto.Advantage = Advantage{
		Percent: model.AdvantagePercent,
	}

	restrictionB, err := model.Restrictions.Value()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(restrictionB.([]byte), &dto.Restrictions); err != nil {
		return nil, err
	}

	return &dto, nil
}

func DtoToModel(dto *PromocodeIn) (*models.Promocode, error) {
	var model models.Promocode

	model.Name = dto.Name
	model.AdvantagePercent = dto.Advantage.Percent

	restrictionB, err := json.Marshal(&dto.Restrictions)
	if err != nil {
		return nil, err
	}
	model.Restrictions = utils.JSON(restrictionB)

	return &model, nil
}
