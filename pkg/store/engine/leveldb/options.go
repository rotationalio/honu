package leveldb

import (
	"github.com/rotationalio/honu/pkg/config"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Returns the default configuration for the underlying storage engine configured
// specifically for Honu's use of LevelDB.
func Options(conf config.StoreConfig) *opt.Options {
	opts := &opt.Options{
		ReadOnly:       conf.ReadOnly,
		ErrorIfMissing: false,
		ErrorIfExist:   false,
		NoSync:         false,
		NoWriteMerge:   false,
		Strict:         opt.DefaultStrict,
	}
	return opts
}
