---
title: Architecture
weight: 500
---

This section is primarily for contributors and developers of the HonuDB replicated database. In it, we will describe the design principles of the database system, how features are implemented and integrated, and the path towards creating complex systems designed from many simple components.

**Design Goal**<br />
The goal of the database is to provide scalable data retrieval both in terms of number of nodes (e.g. scale to 100s of nodes) and amount of data (hundreds of terabytes). In addition to scale, this database provides data access controls, privacy and provenance, and other security related features. Finally, HonuDB provides artifacts and features to support the traning and inferencing model lifecycle. In short, HonuDB is a distributed data governance database for machine learning and artificial intelligence workloads.

**On This Page**

{{< toc >}}

## System Diagram



## Key Terms

<dl>
    <dt>Engine</dt>
    <dd>A database engine is a component that manages how data is stored and cached on disk.</dd>
</dl>