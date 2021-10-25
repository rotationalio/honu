package options

import (
	"errors"
	"fmt"
	"strings"
)

type option struct {
	field string
	value string
}

func parseOptionString(optionString string) ([]option, error) {
	parsedOptionString := strings.FieldsFunc(optionString, splitOptionString)
	err := len(parsedOptionString)%2 != 0
	if err {
		return nil, errors.New("improperly formated option string")
	}
	optionList := []option{}
	for i := 0; i < len(parsedOptionString)-1; i++ {
		field := parsedOptionString[i]
		value := parsedOptionString[i+1]
		parsedOption := option{field, value}
		optionList = append(optionList, parsedOption)
	}
	return optionList, nil
}

func splitOptionString(char rune) bool {
	return char == ',' || char == ' '
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
