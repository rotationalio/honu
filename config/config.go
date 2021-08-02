package config

import "time"

type ReplicaConfig struct {
	Enabled        bool          `split_words:"true" default:"true"`
	BindAddr       string        `split_words:"true" default:":4435"`
	PID            uint64        `split_words:"true" required:"false"`
	Region         string        `split_words:"true" required:"false"`
	Name           string        `split_words:"true" required:"false"`
	GossipInterval time.Duration `split_words:"true" default:"1m"`
	GossipSigma    time.Duration `split_words:"true" default:"5s"`
}
