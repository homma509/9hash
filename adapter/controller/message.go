package controller

import (
	"fmt"

	"gopkg.in/validator.v2"
)

var messages = map[error]string{
	validator.ErrZeroValue: "%sを入力してください。",
}

var names = map[string]string{
	"value":  "値",
	"values": "値",
}

func ToMessages(errs map[string]error) map[string]string {
	msgs := map[string]string{}
	for k, v := range errs {
		name := names[k]
		message, ok := messages[v]
		if ok {
			msgs[k] = fmt.Sprintf(message, name)
		} else {
			msgs[k] = v.Error()
		}
	}
	return msgs
}
