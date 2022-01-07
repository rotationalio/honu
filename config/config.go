package config

import (
	ldbopt "github.com/syndtr/goleveldb/leveldb/opt"
)

// DefaultConfig is used if the user does not specify a configuration
var DefaultConfig = Config{
	Versions: ReplicaConfig{
		PID:    1,
		Region: "local",
		Name:   "localhost",
	},
}

// New creates a configuration with the required options and can also be used to specify
// optional configuration e.g. for engine-specific operations.
func New(options ...Option) Config {
	// Create the default configuration in editable mode
	conf := &Config{
		Versions: ReplicaConfig{
			PID:    DefaultConfig.Versions.PID,
			Region: DefaultConfig.Versions.Region,
			Name:   DefaultConfig.Versions.Name,
		},
	}

	// Apply all options to the configuration
	for _, opt := range options {
		opt(conf)
	}

	// Return the value of the configuration
	return *conf
}

// Config specifies the options necessary to open a Honu database.
type Config struct {
	Versions   ReplicaConfig
	LDBOptions *ldbopt.Options
}

// ReplicaConfig specifies the information needed for the Version manager to maintain
// global object versioning and provenance. This is closely tied to a Replica's
// configuration, e.g. the PID is the process ID of a running replica, the region is
// where the replica is running, and the name is usually the hostname of the replica.
type ReplicaConfig struct {
	PID    uint64 `split_words:"true" required:"false"`
	Region string `split_words:"true" required:"false"`
	Name   string `split_words:"true" required:"false"`
}

// Option modifies a configuration to add optional configuration items.
type Option func(*Config)

func WithReplica(pid uint64, region, name string) Option {
	return func(cfg *Config) {
		cfg.Versions = ReplicaConfig{
			PID:    pid,
			Region: region,
			Name:   name,
		}
	}
}

func WithLevelDB(opt *ldbopt.Options) Option {
	return func(cfg *Config) {
		cfg.LDBOptions = opt
	}
}
