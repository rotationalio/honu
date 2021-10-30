package options

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LeveldbOptions struct{}

func (ldb LeveldbOptions) Read(optionString *string) (*opt.ReadOptions, error) {
	if optionString == nil {
		return nil, nil
	}
	options, err := parse(*optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.ReadOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "DontFillCache":
			returnOption.DontFillCache, err = set(option.value, option.field)
			if err != nil {
				return nil, err
			}
		case "Strict":
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

func (ldb LeveldbOptions) Write(optionString *string) (*opt.WriteOptions, error) {
	if optionString == nil {
		return nil, nil
	}
	options, err := parse(*optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.WriteOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "NoWriteMerge":
			returnOption.NoWriteMerge, err = set(option.value, option.field)
			if err != nil {
				return nil, err
			}
		case "Sync":
			returnOption.Sync, err = set(option.value, option.field)
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
