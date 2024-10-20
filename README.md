# Honu

[![Go Report Card](https://goreportcard.com/badge/github.com/rotationalio/honu)](https://goreportcard.com/report/github.com/rotationalio/honu)
![GitHub Actions CI](https://github.com/rotationalio/honu/actions/workflows/tests.yaml/badge.svg?branch=main)

The Honu Database is an eventually consistent replicated document database that intended for large systems that are distributed globally. Honu uses smart anti-entropy replication to quickly replicate collections across multiple nodes.

Smart anti-entropy uses reinforcement learning with multi-armed bandits to optimize replication. Adaptive consistency reduces costs (ingress and egress data transfer) as well as improves consistency by lowering the likelihood of stale reads or forked writes.

The goal of the database is to provide scalable data retrieval both in terms of number of nodes (e.g. scale to 100s of nodes) and amount of data (hundreds of terabytes). In addition to scale, this database provides data access controls, privacy and provenance, and other security related features. In short, HonuDB is a distributed data governance database for machine learning and artificial intelligence workloads.


## Object Storage

Protocol Buffers are a compact, cross-language compatible data serialization format that facilitates compact network communications. However, in order to make them general purpose and flexible, they require a lot of reflection to work in Go. Since a database is a high performance application, I've implemented a data serialization format that uses no reflection and as a result performs far better at decoding than Protocol Buffers:

```
goos: darwin
goarch: arm64
pkg: github.com/rotationalio/honu/pkg/store
cpu: Apple M1 Max
BenchmarkSerialization/Small/Encode/Honu-10         	  287677	      4170 ns/op	      2323 bytes	   10089 B/op	      20 allocs/op
BenchmarkSerialization/Small/Encode/Protobuf-10     	  374221	      3105 ns/op	      5000 bytes	    3980 B/op	       1 allocs/op
BenchmarkSerialization/Small/Decode/Honu-10         	  419572	      3116 ns/op	    4635 B/op	      60 allocs/op
BenchmarkSerialization/Small/Decode/Protobuf-10     	  141960	      7732 ns/op	    8368 B/op	     128 allocs/op
PASS
ok  	github.com/rotationalio/honu/pkg/store	58.097s
```