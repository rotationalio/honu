---
title: HonuDB Documentation
weight: 0
geekdocNav: false
geekdocAlign: center
geekdocAnchor: false
geekdocBreadcrumb: true
---

<!-- markdownlint-capture -->
<!-- markdownlint-disable MD033 -->

<span class="badge-placeholder">[![Build Status](https://github.com/rotationalio/honu/actions/workflows/tests.yaml/badge.svg)](https://github.com/rotationalio/honu/actions/workflows/tests.yaml)</span>
<span class="badge-placeholder">[![GitHub Release](https://img.shields.io/github/v/release/rotationalio/honu)](https://github.com/rotationalio/honu/releases/latest)</span>
<span class="badge-placeholder">[![GitHub Contributors](https://img.shields.io/github/contributors/rotationalio/honu)](https://github.com/rotationalio/honu/graphs/contributors)</span>
<span class="badge-placeholder">[![License: BSD3](https://img.shields.io/github/license/rotationalio/honu)](https://github.com/rotationalio/honu/blob/main/LICENSE)</span>

<!-- markdownlint-restore -->

The HonuDB Database is the first AI native distributed database intended for an audience of AI developers who need to manage multi-modal datasets with snapshots that can map to models and model training. A replicated document database, HonuDB provides rapid data ingestion and collection management for different mimetypes including JSON, Parquet, images, video, and more. With privacy in mind from the start, HonuDB has data governance features such as provenance and lineage tracking (including by geographic location), and fine-grain access controls. Data scientists and machine learning engineers can rely on Honu to manage small to extremely large datasets replicated over multiple geographic areas.

{{< button size="large" relref="quickstart/" >}}Get Started Now!{{< /button >}}

## Feature overview

{{< columns >}}

### Collections &amp; Datasets

Collections allow you to manage related data together; Datasets are snapshots of collections that indicate exactly what data was usd to train a model.

<--->

### Full Versioning

All objects in the database are fully versioned to prevent an update from changing the view of a dataset from a model perspective.

<--->

### Provenance Awareness

Regions and unique writers are tracked across all updates so you can monitor how data is changing in your system and implement privacy controls.

{{< /columns >}}

{{< columns >}}

### Smart Replication

Honu uses reinforcement learning anti-entropy to maximize consistency and scale replication to hundreds of nodes without increasing your cloud costs.

<--->

### Fine-Grain Access Control

Collections, objects, and datasets have a hierarchical permission model specifically for AI workloads including training and inferencing permissions.

<--->

### Model Context Protocol

Honu supports the [Model Context Protocol](https://modelcontextprotocol.io/introduction) so that you can directly add data to your LLM contexts using semantic similarity indexes.

{{< /columns >}}