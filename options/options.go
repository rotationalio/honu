package options

import (
	"errors"
	"fmt"
	"strings"
)

//Used to generalize an option setting, used by
//parse() which returns a slice of option structs.
type option struct {
	field string
	value string
}

//Parses option strings expected to be in the form of:
//	option1 value, option2 value... etc.
func parse(optionString string) ([]option, error) {
	parsed := strings.FieldsFunc(optionString, split)
	fmt.Println(len(parsed))
	fmt.Println(parsed)
	err := len(parsed)%2 != 0
	fmt.Println(err)
	if err {
		//TODO: think of a better error string to return
		return nil, errors.New("improperly formated option string")
	}

	optionList := []option{}
	for i := 0; i < len(parsed)-1; i = i + 2 {
		field := parsed[i]
		value := parsed[i+1]
		parsedOption := option{field, value}
		optionList = append(optionList, parsedOption)
	}
	return optionList, nil
}

//Function used by parse(), passed to
//strings.FieldsFunc() as the f parameter,
//which controls a string is split.
func split(char rune) bool {
	return char == ',' || char == ' '
}
