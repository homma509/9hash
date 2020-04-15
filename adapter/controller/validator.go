package controller

import (
	"encoding/json"
	"reflect"

	"gopkg.in/validator.v2"
)

type ValidatorSetting struct {
	ArgName      string
	ValidateTags string
}

func init() {
	validator.SetValidationFunc("required", requiredValidator)
}

func validate(params map[string]interface{}, settings []*ValidatorSetting) map[string]error {
	errs := map[string]error{}
	for _, setting := range settings {
		err := validator.Valid(params[setting.ArgName], setting.ValidateTags)
		if err != nil {
			errs[setting.ArgName] = err.(validator.ErrorArray)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func ValidateBody(body string, settings []*ValidatorSetting) map[string]error {
	var b map[string]interface{}
	err := json.Unmarshal([]byte(body), &b)
	if err != nil {
		return map[string]error{}
	}
	return validate(b, settings)
}

func requiredValidator(v interface{}, param string) error {
	if v == nil {
		return validator.ErrZeroValue
	}

	s := reflect.ValueOf(v)

	if s.String() == "" {
		return validator.ErrZeroValue
	}

	return nil
}
