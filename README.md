# Honu

[![Go Report Card](https://goreportcard.com/badge/go.rtnl.ai/honu)](https://goreportcard.com/report/go.rtnl.ai/honu)
[![GitHub Actions CI](https://github.com/rotationalio/honu/actions/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/rotationalio/honu/actions/workflows/tests.yaml)

The Honu Database is an eventually consistent replicated document database that intended for large systems that are distributed globally. Honu uses smart anti-entropy replication to quickly replicate collections across multiple nodes.

Smart anti-entropy uses reinforcement learning with multi-armed bandits to optimize replication. Adaptive consistency reduces costs (ingress and egress data transfer) as well as improves consistency by lowering the likelihood of stale reads or forked writes.

The goal of the database is to provide scalable data retrieval both in terms of number of nodes (e.g. scale to 100s of nodes) and amount of data (hundreds of terabytes). In addition to scale, this database provides data access controls, privacy and provenance, and other security related features. In short, HonuDB is a distributed data governance database for machine learning and artificial intelligence workloads.


## Object Storage

Protocol Buffers are a compact, cross-language compatible data serialization format that facilitates compact network communications. However, in order to make them general purpose and flexible, they require a lot of reflection to work in Go. Since a database is a high performance application, I've implemented a data serialization format that uses no reflection and as a result performs far better at decoding than Protocol Buffers:

```
goos: darwin
goarch: arm64
pkg: go.rtnl.ai/honu/pkg/store/object
cpu: Apple M1 Max
BenchmarkSerialization/Small/Encode-10  	  578818	      1768 ns/op	      4520 bytes	    4487 B/op	       2 allocs/op
BenchmarkSerialization/Small/Decode-10  	  402945	      2686 ns/op	    2341 B/op	      62 allocs/op
BenchmarkSerialization/Medium/Encode-10 	  363483	      3308 ns/op	     10128 bytes	   28081 B/op	       2 allocs/op
BenchmarkSerialization/Medium/Decode-10 	  471076	      2562 ns/op	    2322 B/op	      61 allocs/op
BenchmarkSerialization/Large/Encode-10  	   93942	     12124 ns/op	    303630 bytes	  207933 B/op	       2 allocs/op
BenchmarkSerialization/Large/Decode-10  	  467475	      2736 ns/op	    2334 B/op	      62 allocs/op
BenchmarkSerialization/XLarge/Encode-10 	    7250	    138013 ns/op	   4926099 bytes	 3247592 B/op	       2 allocs/op
BenchmarkSerialization/XLarge/Decode-10 	  407468	      2749 ns/op	    2333 B/op	      62 allocs/op
```