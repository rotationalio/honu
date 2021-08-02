package config

import "time"

// ReplicaConfig specifies the information needed for a Replica and Version manager to
// maintain global object versioning and provenance.
// TODO: this configuration is pulled from a service context not a library context.
type ReplicaConfig struct {
	Enabled        bool          `split_words:"true" default:"true"`
	BindAddr       string        `split_words:"true" default:":4435"`
	PID            uint64        `split_words:"true" required:"false"`
	Region         string        `split_words:"true" required:"false"`
	Name           string        `split_words:"true" required:"false"`
	GossipInterval time.Duration `split_words:"true" default:"1m"`
	GossipSigma    time.Duration `split_words:"true" default:"5s"`
}
