---
title: Data Model
weight: 10
resources:
- name: keylayout
  src: figures/keylayout.svg
  title: Byte Layout of Keys
---

The database engine manages the data as key/value pairs on disk such that the keys are ordered in byte-sort order for fast iteration. Generally speaking the engine uses an LSM-Tree (log structured merge tree) or similar structure for fast appends to the database and routine compaction and merging.

The layout of the `keys` and the data objects are important to understand.

{{< img name="keylayout" size="large" lazy=false >}}