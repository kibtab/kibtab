# STE100 Technical Names: architecture terms

Version 0.0.0-docs. This document holds the architecture Technical Names.

These terms name the layers, ports, and adapters of Kibtab.
Read the [STE100 Technical Names](index.md) dictionary first.

## Technical Name: Core
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** The part of Kibtab that holds the domain logic. It has no I/O.
- **Approved Form:** Core (singular), Cores (plural)
- **Do Not Use:** Kernel, Engine, Domain Layer, Business Logic
- **Correct Example:** *The **core** names a need through a port.*
- **Incorrect Example:** *The kernel names a need through a port.*

## Technical Name: Instance
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** One running copy of Kibtab. It serves one database.
- **Approved Form:** Instance (singular), Instances (plural)
- **Do Not Use:** Engine, Server, Process, Node, Deployment
- **Correct Example:** *Each **instance** holds one version for each table.*
- **Incorrect Example:** *Each engine holds one version for each table.*

## Technical Name: Port
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** An interface in the core that names a need of the domain.
- **Approved Form:** Port (singular), Ports (plural)
- **Do Not Use:** Interface, Contract, Hook, Boundary
- **Correct Example:** *The core declares a **port** for reading a row.*
- **Incorrect Example:** *The core declares an interface for reading a row.*

## Technical Name: Adapter
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** The code that supplies a port. It sits outside the core.
- **Approved Form:** Adapter (singular), Adapters (plural)
- **Do Not Use:** Driver, Connector, Plugin, Implementation
- **Correct Example:** *The PostgreSQL **adapter** supplies the RowRepository port.*
- **Incorrect Example:** *The PostgreSQL driver supplies the RowRepository port.*

## Technical Name: Driven Adapter
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** An adapter that calls out to the core. It drives the core.
- **Approved Form:** Driven Adapter (singular), Driven Adapters (plural)
- **Do Not Use:** Inbound Adapter, Controller, Entry Point, Handler
- **Correct Example:** *The HTTP **driven adapter** calls the write service.*
- **Incorrect Example:** *The HTTP inbound adapter calls the write service.*

## Technical Name: Driving Adapter
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** An adapter that the core calls. It drives from the core.
- **Approved Form:** Driving Adapter (singular), Driving Adapters (plural)
- **Do Not Use:** Outbound Adapter, Gateway, Sender, Exporter
- **Correct Example:** *The database **driving adapter** supplies the port.*
- **Incorrect Example:** *The database outbound adapter supplies the port.*

## Technical Name: Service
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** A use case in the core. It holds one operation of the domain.
- **Approved Form:** Service (singular), Services (plural)
- **Do Not Use:** Use Case, Operation, Workflow, Process
- **Correct Example:** *The sync **service** checks the version before it writes.*
- **Incorrect Example:** *The sync use case checks the version before it writes.*

## Technical Name: Repository
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** A port that reads and writes rows in a table.
- **Approved Form:** Repository (singular), Repositories (plural)
- **Do Not Use:** Store, Gateway, DAO, Data Access Layer
- **Correct Example:** *The **repository** runs the query for one table.*
- **Incorrect Example:** *The store runs the query for one table.*

## Technical Name: Registry
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** A port that lists the tables and the metadata for each table.
- **Approved Form:** Registry (singular), Registries (plural)
- **Do Not Use:** Catalogue, Directory, Index, Map
- **Correct Example:** *The **registry** holds the metadata for each table.*
- **Incorrect Example:** *The catalogue holds the metadata for each table.*

## Technical Name: Clock
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** A port that gives the current time to the core.
- **Approved Form:** Clock (singular), Clocks (plural)
- **Do Not Use:** Timer, Time Source, System Time
- **Correct Example:** *A test uses a fake **clock** for a fixed time.*
- **Incorrect Example:** *A test uses a fake timer for a fixed time.*

## Technical Name: Version
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** The number that the engine holds for each table.
- **Approved Form:** Version (singular), Versions (plural)
- **Do Not Use:** Revision, Generation, Sequence, State
- **Correct Example:** *The engine rejects a **version** that is too old.*
- **Incorrect Example:** *The engine rejects a revision that is too old.*

## Technical Name: Idempotent
- **Part of Speech:** Modifier
- **Category:** Architecture Term
- **Definition:** Describes a write that gives the same result when it runs twice.
- **Approved Form:** Idempotent
- **Do Not Use:** Repeatable, Safe, Deduplicated
- **Correct Example:** *The engine keeps the write **idempotent** for one delta.*
- **Incorrect Example:** *The engine keeps the write repeatable for one delta.*

## Technical Name: Hexagonal Architecture
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** The design that puts the core at the centre and the adapters outside it.
- **Approved Form:** Hexagonal Architecture
- **Do Not Use:** Ports and Adapters, Onion Architecture, Clean Architecture
- **Correct Example:** *Kibtab uses **hexagonal architecture**.*
- **Incorrect Example:** *Kibtab uses the onion architecture.*

## Next Steps

* Read [tooling-terms.md](tooling-terms.md) for the command terms.
* Read [identifiers.md](identifiers.md) for the source names.