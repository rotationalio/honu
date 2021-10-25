package options

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

func CreatePebbleOptions(optionString string) (*pebble.Options, error) {
	options, err := parseOptionString(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := pebble.Options{}
	for _, option := range options {
		switch field := option.field; field {
		// TODO: Impliment generic options
		}
	}
	return &returnOption, nil
}

func CreatePebbleWriteOptions(optionString string) (*pebble.WriteOptions, error) {
	options, err := parseOptionString(optionString)
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

func CreatePebbleIterOptions(optionString string) (*pebble.WriteOptions, error) {
	options, err := parseOptionString(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := pebble.WriteOptions{}
	for _, option := range options {
		switch field := option.field; field {
		// TODO: Impliment iteration options
		}
	}
	return &returnOption, nil
}
