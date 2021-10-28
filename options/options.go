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

func parse(optionString string) ([]option, error) {
	parsed := strings.FieldsFunc(optionString, split)
	err := len(parsed)%2 != 0
	if err {
		return nil, errors.New("improperly formated option string")
	}
	optionList := []option{}
	for i := 0; i < len(parsed)-1; i++ {
		field := parsed[i]
		value := parsed[i+1]
		parsedOption := option{field, value}
		optionList = append(optionList, parsedOption)
	}
	return optionList, nil
}

func split(char rune) bool {
	return char == ',' || char == ' '
}

func set(str string, option string) (bool, error) {
	value := strings.ToLower(str)
	if value == "true" {
		return true, nil
	} else if value == "false" {
		return false, nil
	} else {
		errString := fmt.Sprintf("%s option must be a bool", option)
		err := errors.New(errString)
		return false, err
	}
}
