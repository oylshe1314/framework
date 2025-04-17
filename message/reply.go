package message

import "github.com/oylshe1314/framework/errors"

type Reply struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func newReply(status int, message string, data interface{}) (msg *Reply) {
	return &Reply{Status: status, Message: message, Data: data}
}

func NewReply(v interface{}) *Reply {
	switch t := v.(type) {
	case nil:
		return newReply(0, "success", nil)
	case *Reply:
		return t
	case errors.StatusError:
		return newReply(t.Status(), t.Error(), nil)
	case error:
		return newReply(errors.StatusUnknown, t.Error(), nil)
	default:
		return newReply(0, "success", t)
	}
}
