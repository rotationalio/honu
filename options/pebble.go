package options

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cockroachdb/pebble"
)

//Empty struct used to call Read() and Write()
//to get Pebble read/write options.
type PebbleOptions struct{}

//Parse an option string and returns a Pebble WriteOption object
//with those options set (Pebble does not use read options)
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
			returnOption.Sync, err = strconv.ParseBool(option.value)
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
