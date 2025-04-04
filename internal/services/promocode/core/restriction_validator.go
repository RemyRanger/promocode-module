package core

import (
	"APIs/internal/services/promocode/ports"
	"encoding/json"
	"time"
)

type RestrictionsValidator struct {
	age  int64
	town string
}

func NewRestrictionsValidator(age int64, town string) *RestrictionsValidator {
	return &RestrictionsValidator{
		age:  age,
		town: town,
	}
}

func (r *RestrictionsValidator) validateRestrictions(restrictions []ports.Restriction) []string {
	var reasons []string

	for _, restriction := range restrictions {
		if !r.validateDateRestriction(restriction) {
			resB, _ := json.Marshal(restriction)
			reasons = append(reasons, "unable to validate Date restriction : "+string(resB))
		}
		orReasons := r.validateOrCondition(restriction)
		if len(orReasons) > 0 {
			reasons = append(reasons, orReasons...)
		}
		andReasons := r.validateAndCondition(restriction)
		if len(orReasons) > 0 {
			reasons = append(reasons, andReasons...)
		}
	}

	return reasons
}

func (r *RestrictionsValidator) validateDateRestriction(restriction ports.Restriction) bool {
	dateRestriction, _ := restriction.AsDateRestriction()

	if (ports.DateRestriction{}) != dateRestriction {
		return dateRestriction.Date.After.Before(time.Now()) && dateRestriction.Date.Before.After(time.Now())
	}

	return true
}

func (r *RestrictionsValidator) validateOrCondition(restriction ports.Restriction) []string {
	orCondition, _ := restriction.AsOrCondition()

	var reasons []string
	if len(orCondition.Or) > 0 {
		for _, or := range orCondition.Or {
			if !r.validateAgeRule(or) {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			}
			if !r.validateWeatherRule(or) {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			}
			orReasons := r.validateOrConditionRule(or)
			if len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			}
			andReasons := r.validateAndConditionRule(or)
			if len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			}
		}

		if len(reasons) < len(orCondition.Or) { // at least one cdt is validate so the or is validate, clear reasons
			reasons = []string{}
		}
	}

	return reasons
}

func (r *RestrictionsValidator) validateAndCondition(restriction ports.Restriction) []string {
	andCondition, _ := restriction.AsAndCondition()

	var reasons []string
	if len(andCondition.And) > 0 {
		for _, and := range andCondition.And {
			if !r.validateAgeRule(and) {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			}
			if !r.validateWeatherRule(and) {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			}
			orReasons := r.validateOrConditionRule(and)
			if len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			}
			andReasons := r.validateAndConditionRule(and)
			if len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			}
		}
	}

	return reasons
}

func (r *RestrictionsValidator) validateAgeRule(rule ports.Rule) bool {
	ageRule, _ := rule.AsAgeRule()

	if (ports.AgeRule{}) != ageRule {
		if ageRule.Age.Eq != nil {
			return *ageRule.Age.Eq == r.age
		}
		if ageRule.Age.Gt != nil && ageRule.Age.Lt != nil {
			return *ageRule.Age.Gt <= r.age && *ageRule.Age.Lt >= r.age
		}
	}

	return true
}

func (r *RestrictionsValidator) validateWeatherRule(rule ports.Rule) bool {
	weatherRule, _ := rule.AsWeatherRule()

	if (ports.WeatherRule{}) != weatherRule {
		return true
	}

	return true
}

func (r *RestrictionsValidator) validateOrConditionRule(rule ports.Rule) []string {
	orCondition, _ := rule.AsOrCondition()

	var reasons []string
	if len(orCondition.Or) > 0 {
		for _, or := range orCondition.Or {
			if !r.validateAgeRule(or) {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			}
			if !r.validateWeatherRule(or) {
				resB, _ := json.Marshal(or)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			}
			orReasons := r.validateOrConditionRule(or)
			if len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			}
			andReasons := r.validateAndConditionRule(or)
			if len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			}
		}

		if len(reasons) < len(orCondition.Or) { // at least one cdt is validate so the or is validate, clear reasons
			reasons = []string{}
		}
	}

	return reasons
}

func (r *RestrictionsValidator) validateAndConditionRule(rule ports.Rule) []string {
	andCondition, _ := rule.AsAndCondition()

	var reasons []string
	if len(andCondition.And) > 0 {
		for _, and := range andCondition.And {
			if !r.validateAgeRule(and) {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Age rule : "+string(resB))
			}
			if !r.validateWeatherRule(and) {
				resB, _ := json.Marshal(and)
				reasons = append(reasons, "unable to validate Weather rule : "+string(resB))
			}
			orReasons := r.validateOrConditionRule(and)
			if len(orReasons) > 0 {
				reasons = append(reasons, orReasons...)
			}
			andReasons := r.validateAndConditionRule(and)
			if len(andReasons) > 0 {
				reasons = append(reasons, andReasons...)
			}
		}
	}

	return reasons
}
