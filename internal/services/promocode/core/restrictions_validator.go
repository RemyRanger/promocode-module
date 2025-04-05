package core

import (
	"APIs/internal/common/entities"
	"APIs/internal/services/promocode/ports"
	weather_ports "APIs/internal/services/weather/ports"
	"context"
	"encoding/json"
	"time"
)

type RestrictionsValidator struct {
	ctx            context.Context
	weatherService weather_ports.Service
	age            int64
	town           string
}

func NewRestrictionsValidator(ctx context.Context, weatherService weather_ports.Service, age int64, town string) *RestrictionsValidator {
	return &RestrictionsValidator{
		ctx:            ctx,
		weatherService: weatherService,
		age:            age,
		town:           town,
	}
}

// validateRestrictions returns reasons as []string and error.
// The returned []string is empty if Restrictions are validated
// Input Restriction could be DateRestriction, OrCondition or AndCondition
func (r *RestrictionsValidator) validateRestrictions(restrictions []ports.Restriction) (reasons []string, err error) {
	for _, restriction := range restrictions {
		// Case Date Restriction
		if checked, valid, err := r.validateDateRestriction(restriction); err != nil {
			return nil, err
		} else if !valid && checked {
			resB, _ := json.Marshal(restriction)
			reasons = append(reasons, "unable to validate Date restriction : "+string(resB))
		} else if valid && checked {
			continue
		}

		// Case Or Condition
		if checked, orReasons, err := r.validateOrCondition(restriction); err != nil {
			return nil, err
		} else if checked && len(orReasons) > 0 {
			reasons = append(reasons, orReasons...)
		} else if checked {
			continue
		}

		// Case And Condition
		if checked, andReasons, err := r.validateAndCondition(restriction); err != nil {
			return nil, err
		} else if checked && len(andReasons) > 0 {
			reasons = append(reasons, andReasons...)
		} else if checked {
			continue
		}
	}

	return reasons, nil
}

// validateDateRestriction returns checked as boolean, valid as boolean and error.
// The returned checked is true if Restriction is DateRestriction.
// The returned valid is true if Restriction is evaluated as valid.
func (r *RestrictionsValidator) validateDateRestriction(restriction ports.Restriction) (checked bool, valid bool, err error) {
	dateRestriction, err := restriction.AsDateRestriction()
	if err != nil {
		return false, false, err
	}

	if (ports.DateRestriction{}) != dateRestriction {
		return true, dateRestriction.Date.After.Before(time.Now()) && dateRestriction.Date.Before.After(time.Now()), nil
	}

	return false, false, nil
}

func (r *RestrictionsValidator) validateOrCondition(restriction ports.Restriction) (checked bool, reasons []string, err error) {
	orCondition, err := restriction.AsOrCondition()
	if err != nil {
		return false, reasons, err
	}

	if len(orCondition.Or) > 0 {
		for _, or := range orCondition.Or {
			// Case Age Rule
			if checked, valid, err := r.validateAgeRule(or); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			} else if valid && checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}

			// Case Weather Rule
			if checked, valid, err := r.validateWeatherRule(or); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			} else if valid && checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}

			// Case Or Rule
			if checked, orReasons, err := r.validateOrConditionRule(or); err != nil {
				return true, reasons, err
			} else if checked && len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			} else if checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}

			// Case And Rule
			if checked, andReasons, err := r.validateAndConditionRule(or); err != nil {
				return true, reasons, err
			} else if checked && len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			} else if checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}
		}

		return true, reasons, nil
	}

	return false, reasons, nil
}

func (r *RestrictionsValidator) validateAndCondition(restriction ports.Restriction) (checked bool, reasons []string, err error) {
	andCondition, err := restriction.AsAndCondition()
	if err != nil {
		return false, reasons, err
	}

	if len(andCondition.And) > 0 {
		for _, and := range andCondition.And {
			// Case Age Rule
			if checked, valid, err := r.validateAgeRule(and); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			} else if valid && checked {
				continue
			}

			// Case Weather Rule
			if checked, valid, err := r.validateWeatherRule(and); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			} else if valid && checked {
				continue
			}

			// Case Or Rule
			if checked, orReasons, err := r.validateOrConditionRule(and); err != nil {
				return true, reasons, err
			} else if checked && len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			} else if checked {
				continue
			}

			// Case And Rule
			if checked, andReasons, err := r.validateAndConditionRule(and); err != nil {
				return true, reasons, err
			} else if checked && len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			} else if checked {
				continue
			}
		}

		return true, reasons, nil
	}

	return false, reasons, nil
}

func (r *RestrictionsValidator) validateAgeRule(rule ports.Rule) (checked bool, valid bool, err error) {
	ageRule, err := rule.AsAgeRule()
	if err != nil {
		return false, false, err
	}

	if (ports.AgeRule{}) != ageRule {
		if ageRule.Age.Eq != nil {
			return true, *ageRule.Age.Eq == r.age, nil
		}
		if ageRule.Age.Gt != nil && ageRule.Age.Lt != nil {
			return true, *ageRule.Age.Gt <= r.age && *ageRule.Age.Lt >= r.age, nil
		}
		if ageRule.Age.Gt != nil {
			return true, *ageRule.Age.Gt <= r.age, nil
		}
		if ageRule.Age.Lt != nil {
			return true, *ageRule.Age.Lt >= r.age, nil
		}
	}

	return false, false, nil
}

func (r *RestrictionsValidator) validateWeatherRule(rule ports.Rule) (checked bool, valid bool, err error) {
	weatherRule, err := rule.AsWeatherRule()
	if err != nil {
		return false, false, err
	}

	if (ports.WeatherRule{}) != weatherRule {
		valid, err := r.weatherService.ValidateWeather(r.ctx, entities.WeatherQuery{
			Town:    r.town,
			TempMin: weatherRule.Weather.Temp.Gt,
			Type:    entities.WeatherType(weatherRule.Weather.Is),
		})
		return true, valid, err
	}

	return false, false, nil
}

func (r *RestrictionsValidator) validateOrConditionRule(rule ports.Rule) (checked bool, reasons []string, err error) {
	orCondition, err := rule.AsOrCondition()
	if err != nil {
		return false, nil, err
	}

	if len(orCondition.Or) > 0 {
		for _, or := range orCondition.Or {
			// Case Age Rule
			if checked, valid, err := r.validateAgeRule(or); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			} else if valid && checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}

			// Case Weather Rule
			if checked, valid, err := r.validateWeatherRule(or); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			} else if valid && checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}

			// Case Or Rule
			if checked, orReasons, err := r.validateOrConditionRule(or); err != nil {
				return true, reasons, err
			} else if checked && len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			} else if checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}

			// Case And Rule
			if checked, andReasons, err := r.validateAndConditionRule(or); err != nil {
				return true, reasons, err
			} else if checked && len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			} else if checked {
				reasons = []string{} // erase reasons
				break                // or success, don't go further
			}
		}

		return true, reasons, nil
	}

	return false, reasons, nil
}

func (r *RestrictionsValidator) validateAndConditionRule(rule ports.Rule) (checked bool, reasons []string, err error) {
	andCondition, err := rule.AsAndCondition()
	if err != nil {
		return false, reasons, err
	}

	if len(andCondition.And) > 0 {
		for _, and := range andCondition.And {
			// Case Age Rule
			if checked, valid, err := r.validateAgeRule(and); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			} else if valid && checked {
				continue
			}

			// Case Weather Rule
			if checked, valid, err := r.validateWeatherRule(and); err != nil {
				return true, reasons, err
			} else if !valid && checked {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			} else if valid && checked {
				continue
			}

			// Case Or Rule
			if checked, orReasons, err := r.validateOrConditionRule(and); err != nil {
				return true, reasons, err
			} else if checked && len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			} else if checked {
				continue
			}

			// Case And Rule
			if checked, andReasons, err := r.validateAndConditionRule(and); err != nil {
				return true, reasons, err
			} else if checked && len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			} else if checked {
				continue
			}
		}

		return true, reasons, nil
	}

	return false, reasons, nil
}
