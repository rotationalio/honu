package options

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

func CreatePebbleWriteOptions(optionString *string) (*pebble.WriteOptions, error) {
	if optionString == nil {
		return nil, nil
	}
	options, err := parseOptionString(*optionString)
	if err != nil {
		return nil, err
	}
	returnOption := pebble.WriteOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "Sync":
			returnOption.Sync, err = setBoolOption(option.value, option.field)
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

func CreatePebbleIterOptions(optionString *string) (*pebble.IterOptions, error) {
	if optionString == nil {
		return nil, nil
	}
	options, err := parseOptionString(*optionString)
	if err != nil {
		return nil, err
	}
	returnOption := pebble.IterOptions{}
	for _, option := range options {
		switch field := option.field; field {
		// TODO: Impliment iteration options
		}
	}
	return &returnOption, nil
}
