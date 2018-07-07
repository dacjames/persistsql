Persistent Resource Model
=====

Implements an extensible, persistent resource model with support for Postgres and sqlite.

The resource model tables and views are organized into namespaces, where namespacing is implemented with schemas in Postgres and naming conventions in sqlite. Those namespaces are:

#### ledger

The ledger namespace is the core of the storage model. All tables in ledger are append-only and form something like a relational version of an [OR-Set CRDT](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type). This means that the ledger can be recplicated without the risk of conflicts. The core ledger tables are:

##### ledger.services

Stores a list of all services. The primary purpose of this table is to implement dynamic resource lookup where only the ID is known.

Because the number of services is small, the services table is replicated to all partitions.

##### ledger.resources


