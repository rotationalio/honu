package credentials

type Credentials interface {
	AccessToken() (string, error)
}

type Token string

func (t Token) AccessToken() (string, error) {
	if string(t) == "" {
		return "", ErrInvalidCredentials
	}
	return string(t), nil
}
