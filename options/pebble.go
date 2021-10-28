package options

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

type PebbleOptions struct{}

func (p PebbleOptions) Write(optionString *string) (*pebble.WriteOptions, error) {
	if optionString == nil {
		return nil, nil
	}
	options, err := parse(*optionString)
	if err != nil {
		return nil, err
	}
	returnOption := pebble.WriteOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "Sync":
			returnOption.Sync, err = set(option.value, option.field)
			if err != nil {
				return nil, err
			}
		default:
			errString := fmt.Sprintf("%s is not a valid pebble writeoption", option.field)
			err = errors.New(errString)
			return nil, err
		}
	}
	return &returnOption, nil
}
