package controller

import (
	"fmt"

	"gopkg.in/validator.v2"
)

var messages = map[error]string{
	validator.ErrZeroValue: "%sを入力してください。",
}

var names = map[string]string{
	"value":  "value",
	"values": "values",
}

func ToMessages(errs map[string]error) map[string]string {
	msgs := map[string]string{}
	for k, v := range errs {
		name := names[k]
		message, ok := messages[v]
		if ok {
			msgs[name] = fmt.Sprintf(message, name)
		} else {
			msgs[name] = v.Error()
		}
	}
	return msgs
}
