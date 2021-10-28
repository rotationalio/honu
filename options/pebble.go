package options

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

type pebbleOptions struct{}

func (p *pebbleOptions) Read(optionString *string) error {
	errorString := fmt.Sprintf("Pebble does not support readoptions")
	err := errors.New(errorString)
	return err
}

func (p *pebbleOptions) Write(optionString *string) (*pebble.WriteOptions, error) {
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
