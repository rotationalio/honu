package config

import (
	"github.com/rotationalio/confire"
	"github.com/rotationalio/honu/pkg/logger"
	"github.com/rs/zerolog"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be HONU_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const Prefix = "honu"

// Config contains all of the configuration parameters for an honudb instance and is
// loaded from the environment or a configuration file with reasonable defaults for
// values that are omitted. The Config should be validated in preparation for running
// the honudb instance to ensure that all server operations work as expected.
type Config struct {
	Maintenance bool                `default:"false" desc:"if true, the replica will start in maintenance mode"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info" desc:"specify the verbosity of logging (trace, debug, info, warn, error, fatal panic)"`
	ConsoleLog  bool                `split_words:"true" default:"false" desc:"if true logs colorized human readable output instead of json"`
	processed   bool
}

func New() (conf Config, err error) {
	if err = confire.Process(Prefix, &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Returns true if the config has not been correctly processed from the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Custom validations are added here, particularly validations that require one or more
// fields to be processed before the validation occurs.
func (c Config) Validate() (err error) {
	return err
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}
