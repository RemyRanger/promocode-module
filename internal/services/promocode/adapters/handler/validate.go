package handler

import (
	"APIs/internal/services/promocode/ports"

	"github.com/go-playground/validator/v10"
)

func validateRestrictions(restrictions []ports.Restriction) error {

	for _, restriction := range restrictions {
		if err := validateDateRestriction(restriction); err != nil {
			return err
		}
		if err := validateOrCondition(restriction); err != nil {
			return err
		}
		if err := validateAndCondition(restriction); err != nil {
			return err
		}
	}

	return nil
}

func validateDateRestriction(restriction ports.Restriction) error {
	dateRestriction, err := restriction.AsDateRestriction()
	if err != nil {
		return err
	}

	if (ports.DateRestriction{}) != dateRestriction {
		validate := validator.New()
		if err := validate.Struct(dateRestriction); err != nil {
			return err
		}
	}

	return nil
}

func validateOrCondition(restriction ports.Restriction) error {
	orCondition, err := restriction.AsOrCondition()
	if err != nil {
		return err
	}

	if len(orCondition.Or) > 0 {
		for _, or := range orCondition.Or {
			if err := validateAgeRule(or); err != nil {
				return err
			}
			if err := validateWeatherRule(or); err != nil {
				return err
			}
			if err := validateOrConditionRule(or); err != nil {
				return err
			}
			if err := validateAndConditionRule(or); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateAndCondition(restriction ports.Restriction) error {
	andCondition, err := restriction.AsAndCondition()
	if err != nil {
		return err
	}

	if len(andCondition.And) > 0 {
		for _, and := range andCondition.And {
			if err := validateAgeRule(and); err != nil {
				return err
			}
			if err := validateWeatherRule(and); err != nil {
				return err
			}
			if err := validateOrConditionRule(and); err != nil {
				return err
			}
			if err := validateAndConditionRule(and); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateAgeRule(rule ports.Rule) error {
	ageRule, err := rule.AsAgeRule()
	if err != nil {
		return err
	}

	if (ports.AgeRule{}) != ageRule {
		validate := validator.New()
		if err := validate.Struct(ageRule); err != nil {
			return err
		}
	}

	return nil
}

func validateWeatherRule(rule ports.Rule) error {
	weatherRule, err := rule.AsWeatherRule()
	if err != nil {
		return err
	}

	if (ports.WeatherRule{}) != weatherRule {
		validate := validator.New()
		if err := validate.Struct(weatherRule); err != nil {
			return err
		}
	}

	return nil
}

func validateOrConditionRule(rule ports.Rule) error {
	orCondition, err := rule.AsOrCondition()
	if err != nil {
		return err
	}

	if len(orCondition.Or) > 0 {
		for _, or := range orCondition.Or {
			if err := validateAgeRule(or); err != nil {
				return err
			}
			if err := validateWeatherRule(or); err != nil {
				return err
			}
			if err := validateOrConditionRule(or); err != nil {
				return err
			}
			if err := validateAndConditionRule(or); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateAndConditionRule(rule ports.Rule) error {
	andCondition, err := rule.AsAndCondition()
	if err != nil {
		return err
	}

	if len(andCondition.And) > 0 {
		for _, and := range andCondition.And {
			if err := validateAgeRule(and); err != nil {
				return err
			}
			if err := validateWeatherRule(and); err != nil {
				return err
			}
			if err := validateOrConditionRule(and); err != nil {
				return err
			}
			if err := validateAndConditionRule(and); err != nil {
				return err
			}
		}
	}

	return nil
}
