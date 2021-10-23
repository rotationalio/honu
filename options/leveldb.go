package options

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

func CreateLeveldbReadOption(optionString string) (*opt.ReadOptions, error) {
	options, err := parseOptionString(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.ReadOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "DontFillCache":
			returnOption.DontFillCache, err = setBoolOption(option.value, option.field)
			if err != nil {
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

func CreateLeveldbWriteOption(optionString string) (*opt.WriteOptions, error) {
	options, err := parseOptionString(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.WriteOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "NoWriteMerge":
			returnOption.NoWriteMerge, err = setBoolOption(option.value, option.field)
			if err != nil {
				return nil, err
			}
		case "Sync":
			returnOption.Sync, err = setBoolOption(option.value, option.field)
			if err != nil {
				return nil, err
			}
		default:
			errString := fmt.Sprintf("%s is not a valid leveldb writeoption", option.field)
			err = errors.New(errString)
			return nil, err
		}
	}
	return &returnOption, nil
}

func setBoolOption(boolString string, optionString string) (bool, error) {
	value := strings.ToLower(boolString)
	if value == "true" {
		return true, nil
	} else if value == "false" {
		return false, nil
	} else {
		errString := fmt.Sprintf("%s is not a valid setting for the %s option", value, optionString)
		err := errors.New(errString)
		return false, err
	}
}
