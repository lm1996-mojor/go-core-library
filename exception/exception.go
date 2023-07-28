package exception

import (
	"fmt"

	"mojor/go-core-library/log"

	"github.com/pkg/errors"
)

// Error contains code and message
type Error struct {
	Code int
	Msg  string
}

// Panic a customer error
func Panic(code int, msg string) {
	panic(Error{
		Code: code,
		Msg:  msg,
	})
}

// CheckErr if error is not nil, it will panic the code and msg
func CheckErr(err error, code int, msg string) {
	if err != nil {
		err = errors.WithStack(err)
		log.Error(fmt.Sprintf("%+v", err))
		text := msg
		if len(text) == 0 {
			text = err.Error()
		}
		Panic(code, text)
	}
}

// ParseErr extract the code and message from the error
func ParseErr(e interface{}) (code int, msg string) {
	if err, ok := e.(Error); ok {
		return err.Code, err.Msg
	}
	if err, ok := e.(error); ok {
		log.Error(fmt.Sprintf("%+v", err))

		return 500, err.Error()
	}
	if err, ok := e.(string); ok {
		log.Error(err)

		return 500, err
	}

	if err, ok := e.(fmt.Stringer); ok {
		log.Error(err.String())

		return 500, err.String()
	}

	return 500, "System error"
}
