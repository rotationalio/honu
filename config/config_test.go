package config_test

import (
	"testing"

	"github.com/rotationalio/honu/config"
	"github.com/stretchr/testify/require"
	ldbopt "github.com/syndtr/goleveldb/leveldb/opt"
)

func TestConfig(t *testing.T) {
	// Test Default Config
	conf, err := config.New()
	require.NoError(t, err)
	require.Equal(t, config.DefaultConfig, conf)
	require.NotZero(t, conf.Versions.PID)
	require.NotEmpty(t, conf.Versions.Region)

	// Test WithVersions
	conf, err = config.New(config.WithReplica(config.ReplicaConfig{8, "us-antarctic-23", "research"}))
	require.NoError(t, err)
	require.NotEmpty(t, conf.Versions)
	require.Equal(t, uint64(8), conf.Versions.PID)
	require.Equal(t, "us-antarctic-23", conf.Versions.Region)
	require.Equal(t, "research", conf.Versions.Name)

	// Test WithLevelDB Options
	conf, err = config.New(config.WithLevelDB(&ldbopt.Options{Strict: ldbopt.StrictJournal}))
	require.NoError(t, err)
	require.Equal(t, config.DefaultConfig.Versions, conf.Versions)
	require.NotNil(t, conf.LDBOptions)
	require.Equal(t, conf.LDBOptions.Strict, ldbopt.StrictJournal)
}
