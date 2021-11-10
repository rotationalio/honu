package options

import (
	"github.com/cockroachdb/pebble"
	ldb "github.com/syndtr/goleveldb/leveldb/opt"
)

type Options struct {
	LeveldbRead  *ldb.ReadOptions
	LeveldbWrite *ldb.WriteOptions
	PebbleWrite  *pebble.WriteOptions
}

type SetOptions func(cfg *Options) error

func WithLeveldbRead(leveldbRead *ldb.ReadOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LeveldbRead = leveldbRead
		return nil
	}
}

func WithLeveldbWrite(leveldbWrite *ldb.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LeveldbWrite = leveldbWrite
		return nil
	}
}

func WithPebbleWrite(pebbleWrite *pebble.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.PebbleWrite = pebbleWrite
		return nil
	}
}

//Used to generalize an option setting, used by
//parse() which returns a slice of option structs.
// type option struct {
// 	field string
// 	value string
// }
//
// //Parses option strings expected to be in the form of:
// //	option1 value, option2 value... etc.
// func parse(optionString string) ([]option, error) {
// 	optionSlice := []option{}
// 	optionPairs := strings.Split(optionString, ",")
// 	for _, optionPair := range optionPairs {
// 		//Handle whitespace around a comma by trimming.
// 		optionPair = strings.TrimSpace(optionPair)
// 		options := strings.Split(optionPair, " ")
// 		if len(options) != 2 {
// 			//TODO: Come up with a better error to return.
// 			return nil, errors.New("improperly formated option string")
// 		}
// 		optionSlice = append(optionSlice, option{field: options[0], value: options[1]})
// 	}
// 	return optionSlice, nil
// }
