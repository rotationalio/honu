package options

import (
	"errors"
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
	optionSlice := []option{}
	optionPairs := strings.Split(optionString, ",")
	for _, optionPair := range optionPairs {
		//Handle whitespace after a comma by trimming.
		optionPair = strings.TrimSpace(optionPair)
		options := strings.Split(optionPair, " ")
		if len(options) != 2 {
			//TODO: Come up with a better error to return.
			return nil, errors.New("improperly formated option string")
		}
		optionSlice = append(optionSlice, option{field: options[0], value: options[1]})
	}
	return optionSlice, nil
}
