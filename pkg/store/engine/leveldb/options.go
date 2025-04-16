package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"go.rtnl.ai/honu/pkg/config"
)

// Returns the default configuration for the underlying storage engine configured
// specifically for Honu's use of LevelDB.
func Options(conf config.StoreConfig) *opt.Options {
	opts := &opt.Options{
		Compression:    opt.SnappyCompression,
		ErrorIfMissing: false,
		ErrorIfExist:   false,
		Filter:         filter.NewBloomFilter(100),
		NoSync:         false,
		NoWriteMerge:   false,
		ReadOnly:       conf.ReadOnly,
		Strict:         opt.DefaultStrict,
	}
	return opts
}
