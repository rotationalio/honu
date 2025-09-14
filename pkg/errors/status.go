package errors

import (
	"fmt"

	"go.rtnl.ai/honu/pkg/api/v1"
)

func Status(code int, err any) *api.StatusError {
	return &api.StatusError{
		Code:  code,
		Reply: api.Error(err),
	}
}

func Statusf(code int, format string, args ...any) *api.StatusError {
	return &api.StatusError{
		Code:  code,
		Reply: api.Error(fmt.Errorf(format, args...)),
	}
}
