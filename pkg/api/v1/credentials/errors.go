package credentials

import "errors"

var (
	ErrInvalidCredentials = errors.New("missing, invalid or expired credentials")
)
