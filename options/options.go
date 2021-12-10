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
func New(options ...SetOptions) (cfg *Options, err error) {
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
	LevelDBRead  *ldb.ReadOptions
	LevelDBWrite *ldb.WriteOptions
	PebbleWrite  *pebble.WriteOptions
	Namespace    string
	Destroy      bool
}

// Defines the signature of functions accepted as parameters by Honu methods.
type SetOptions func(cfg *Options) error

// WithNamespace returns a closure that sets a namespace other than the default.
func WithNamespace(namespace string) SetOptions {
	return func(cfg *Options) error {
		// If namespace is empty, keep default namespace
		if namespace != "" {
			cfg.Namespace = namespace
		}
		return nil
	}
}

// WithDestroy returns a closure that sets an option so that instead of creating a
// tombstone, the data is permanently deleted from the database, meaning that a
// subsequent Put will restart the versioning. Use Destroy with care in a distributed
// system, if you destroy a key in an anti-entropy environment it will simply be
// repaired by the system to the latest version before the Delete.
func WithDestroy() SetOptions {
	return func(cfg *Options) error {
		cfg.Destroy = true
		return nil
	}
}

//Closure returning a function that adds the leveldbRead
//parameter to an Options struct's LeveldbRead field.
func WithLeveldbRead(opts *ldb.ReadOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LevelDBRead = opts
		return nil
	}
}

//Closure returning a function that adds the leveldbWrite
//parameter to an Options struct's LeveldbWrite field.
func WithLeveldbWrite(opts *ldb.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LevelDBWrite = opts
		return nil
	}
}

//Closure returning a function that adds the pebbleWrite
//parameter to an Options struct's PebbleWrite field.
func WithPebbleWrite(opts *pebble.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.PebbleWrite = opts
		return nil
	}
}
