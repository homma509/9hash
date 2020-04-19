package controller

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/golang/glog"
	"gopkg.in/validator.v2"
)

var (
	ErrUint = validator.TextErr{Err: errors.New("invalid uint")}
)

type ValidatorSetting struct {
	ArgName      string
	ValidateTags string
}

func init() {
	validator.SetValidationFunc("required", requiredValidator)
	validator.SetValidationFunc("uint", uintValidator)
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

func uintValidator(v interface{}, param string) error {
	if v == nil {
		return nil
	}

	s := reflect.ValueOf(v)

	if s.String() == "" {
		return nil
	}

	var n int

	switch s.Kind() {
	case reflect.String:
		n64, err := strconv.ParseInt(s.String(), 10, 64)
		if err != nil {
			glog.Warningf("%s:%s", param, err.Error())
			return validator.ErrUnsupported
		}
		n = int(n64)
	case reflect.Int:
		n = v.(int)
	case reflect.Float64:
		n = int(v.(float64))
	default:
		glog.Warningf("%s:%s", param, s.Kind())
		return validator.ErrUnsupported
	}

	if n < 0 {
		return ErrUint
	}

	return nil
}
