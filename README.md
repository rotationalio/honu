# Honu

[![Go Report Card](https://goreportcard.com/badge/github.com/rotationalio/honu)](https://goreportcard.com/report/github.com/rotationalio/honu)
![GitHub Actions CI](https://github.com/rotationalio/honu/actions/workflows/tests.yaml/badge.svg?branch=main)

The Honu Database is an eventually consistent replicated document database that intended for large systems that are distributed globally. Honu uses smart anti-entropy replication to quickly replicate collections across multiple nodes.

Smart anti-entropy uses reinforcement learning with multi-armed bandits to optimize replication. Adaptive consistency reduces costs (ingress and egress data transfer) as well as improves consistency by lowering the likelihood of stale reads or forked writes.

The goal of the database is to provide scalable data retrieval both in terms of number of nodes (e.g. scale to 100s of nodes) and amount of data (hundreds of terabytes). In addition to scale, this database provides data access controls, privacy and provenance, and other security related features. In short, HonuDB is a distributed data governance database for machine learning and artificial intelligence workloads.
