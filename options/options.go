package options

import (
	"errors"
	"strings"
)

type option struct {
	field string
	value string
}

/* func CheckOptions(option Option) (*Option, error) {
	optType := reflect.ValueOf(option.optionObject)
	readOpts, err := parseOptionString(option)
	if err != nil {
		print(err)
		return nil, err
	}
	for _, possibleOpt := range readOpts {
		checkField := optType.Elem().FieldByName(possibleOpt)
		if !checkField.IsValid() {
			errorString := fmt.Sprintf("%s is not a valid option", possibleOpt)
			err = errors.New(errorString)
			return nil, err
		} else {
			optType
		}
	}
	return &option, nil
} */

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
