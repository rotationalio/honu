package options

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

func CreateReadOption(optionString string) (*opt.ReadOptions, error) {
	options, err := parseOptionString(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.ReadOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "DontFillCache":
			value := strings.ToLower(option.value)
			if value == "true" {
				returnOption.DontFillCache = true
			} else if value == "false" {
				returnOption.DontFillCache = false
				return nil, errors.New("")
			} else {
				errString := fmt.Sprintf("%s is not a valid setting for the DontFillCache option", value)
				err = errors.New(errString)
				return nil, err
			}
		case "strict":
			strictInt, err := strconv.Atoi(option.value)
			if err != nil {
				return nil, err
			}
			strict := opt.Strict(strictInt)
			returnOption.Strict = strict
		default:
			errString := fmt.Sprintf("%s is not a valid leveldb readoption", option.field)
			err = errors.New(errString)
			return nil, err
		}
	}
	return &returnOption, nil
}

func CreateWriteOption(optionString string) (*opt.WriteOptions, error) {
	return nil, nil
}
