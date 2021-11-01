package options

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

//Empty struct used to call Read() and Write()
//to get LevelDB read/write options.
type LeveldbOptions struct{}

//Parses an option string and returns a LevelDB ReadOption
//struct with those options set.
func (ldb LeveldbOptions) Read(optionString string) (*opt.ReadOptions, error) {
	if len(strings.TrimSpace(optionString)) == 0 {
		return nil, nil
	}
	options, err := parse(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.ReadOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "DontFillCache":
			returnOption.DontFillCache, err = strconv.ParseBool(option.value)
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

//Parses an option string and returns a LevelDB WriteOption
//struct with those options set.
func (ldb LeveldbOptions) Write(optionString string) (*opt.WriteOptions, error) {
	if len(strings.TrimSpace(optionString)) == 0 {
		return nil, nil
	}
	options, err := parse(optionString)
	if err != nil {
		return nil, err
	}
	returnOption := opt.WriteOptions{}
	for _, option := range options {
		switch field := option.field; field {
		case "NoWriteMerge":
			returnOption.NoWriteMerge, err = strconv.ParseBool(option.value)
			if err != nil {
				return nil, err
			}
		case "Sync":
			returnOption.Sync, err = strconv.ParseBool(option.value)
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
