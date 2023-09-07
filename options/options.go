package options

import (
	"github.com/cockroachdb/pebble"
	ldb "github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	NamespaceDefault = "default"
)

// New creates a per-call Options object based on the variadic SetOptions closures
// supplied by the user. New also sets sensible defaults for various options.
func New(options ...Option) (cfg *Options, err error) {
	cfg = &Options{Namespace: NamespaceDefault}
	for _, option := range options {
		if err = option(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// Contains all available read/write options for the supported engines. Fields are set
// by closures implementing the SetOptions signature.
type Options struct {
	LevelDBRead      *ldb.ReadOptions
	LevelDBWrite     *ldb.WriteOptions
	PebbleWrite      *pebble.WriteOptions
	Namespace        string
	Force            bool
	Tombstones       bool
	RequireExists    bool
	RequireNotExists bool
}

// Defines the signature of functions accepted as parameters by Honu methods.
type Option func(cfg *Options) error

// WithNamespace returns a closure that sets a namespace other than the default.
func WithNamespace(namespace string) Option {
	return func(cfg *Options) error {
		// If namespace is empty, keep default namespace
		if namespace != "" {
			cfg.Namespace = namespace
		}
		return nil
	}
}

// WithForce prevents validation checks from returning an error during accesses.
func WithForce() Option {
	return func(cfg *Options) error {
		cfg.Force = true
		return nil
	}
}

// WithTombstones causes the iterator to include tombstones as its iterating.
func WithTombstones() Option {
	return func(cfg *Options) error {
		cfg.Tombstones = true
		return nil
	}
}

// Closure returning a function that adds the leveldbRead
// parameter to an Options struct's LeveldbRead field.
func WithLevelDBRead(opts *ldb.ReadOptions) Option {
	return func(cfg *Options) error {
		cfg.LevelDBRead = opts
		return nil
	}
}

// Closure returning a function that adds the leveldbWrite
// parameter to an Options struct's LeveldbWrite field.
func WithLevelDBWrite(opts *ldb.WriteOptions) Option {
	return func(cfg *Options) error {
		cfg.LevelDBWrite = opts
		return nil
	}
}

// Closure returning a function that adds the pebbleWrite
// parameter to an Options struct's PebbleWrite field.
func WithPebbleWrite(opts *pebble.WriteOptions) Option {
	return func(cfg *Options) error {
		cfg.PebbleWrite = opts
		return nil
	}
}
