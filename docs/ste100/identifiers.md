# STE100 Technical Names: identifiers

Version 0.0.0-docs. This document holds the identifier Technical Names.

These terms are the exact names in the Kibtab source and build files.
Use the exact case of each name in every document.
Read [the catalogue](index.md) first.

## Technical Name: RowRepository
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The port that reads and writes rows in one table.
- **Approved Form:** RowRepository
- **Do Not Use:** RowStore, RowGateway, DataAccess
- **Correct Example:** *The core declares `RowRepository` in `ports/`.*
- **Incorrect Example:** *The core declares `RowStore` in `ports/`.*

## Technical Name: TableRegistry
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The port that lists the tables and their metadata.
- **Approved Form:** TableRegistry
- **Do Not Use:** TableCatalogue, SchemaRegistry, TableIndex
- **Correct Example:** *The core declares `TableRegistry` in `ports/`.*
- **Incorrect Example:** *The core declares `TableCatalogue` in `ports/`.*

## Technical Name: AuditWriter
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The port that records a change to a cell.
- **Approved Form:** AuditWriter
- **Do Not Use:** AuditLogger, HistoryWriter, ChangeRecorder
- **Correct Example:** *The core declares `AuditWriter` in `ports/`.*
- **Incorrect Example:** *The core declares `AuditLogger` in `ports/`.*

## Technical Name: Dialect
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The port that quotes a name and pages a query for one engine.
- **Approved Form:** Dialect
- **Do Not Use:** EngineDialect, SqlDialect, DialectPort
- **Correct Example:** *The core declares `Dialect` in `ports/`.*
- **Incorrect Example:** *The core declares `SqlDialect` in `ports/`.*

## Technical Name: TransactionRunner
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The port that runs a group of writes as one transaction.
- **Approved Form:** TransactionRunner
- **Do Not Use:** TransactionManager, UnitOfWork, TxRunner
- **Correct Example:** *The core declares `TransactionRunner` in `ports/`.*
- **Incorrect Example:** *The core declares `UnitOfWork` in `ports/`.*

## Technical Name: Clock
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The port that gives the current time to the core.
- **Approved Form:** Clock
- **Do Not Use:** TimeProvider, SystemClock, Timer
- **Correct Example:** *The core declares `Clock` in `ports/`.*
- **Incorrect Example:** *The core declares `TimeProvider` in `ports/`.*

## Technical Name: SyncService
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The service that accepts a delta and returns a result.
- **Approved Form:** SyncService
- **Do Not Use:** SyncUseCase, SyncHandler, DeltaService
- **Correct Example:** *The core declares `SyncService` in `services/`.*
- **Incorrect Example:** *The core declares `SyncHandler` in `services/`.*

## Technical Name: TableMetadata
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The domain model that holds the name, fields, and version of one table.
- **Approved Form:** TableMetadata
- **Do Not Use:** TableMeta, TableInfo, TableDescription
- **Correct Example:** *The core declares `TableMetadata` in `domain/`.*
- **Incorrect Example:** *The core declares `TableMeta` in `domain/`.*

## Technical Name: CellValue
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The domain model that holds one cell value and its field name.
- **Approved Form:** CellValue
- **Do Not Use:** Cell, FieldValue, CellData
- **Correct Example:** *The core declares `CellValue` in `domain/`.*
- **Incorrect Example:** *The core declares `Cell` in `domain/`.*

## Technical Name: CellDelta
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The domain model that holds one cell change in a sync payload.
- **Approved Form:** CellDelta
- **Do Not Use:** CellChange, CellEdit, CellPatch
- **Correct Example:** *The core declares `CellDelta` in `domain/`.*
- **Incorrect Example:** *The core declares `CellChange` in `domain/`.*

## Technical Name: SyncPayload
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The domain model that holds a table name and a delta from a client.
- **Approved Form:** SyncPayload
- **Do Not Use:** SyncRequest, DeltaPayload, TableDelta
- **Correct Example:** *The core declares `SyncPayload` in `domain/`.*
- **Incorrect Example:** *The core declares `SyncRequest` in `domain/`.*

## Technical Name: SyncResult
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The domain model that holds the result after a sync.
- **Approved Form:** SyncResult
- **Do Not Use:** SyncResponse, SyncReply, SyncStatus
- **Correct Example:** *The core declares `SyncResult` in `domain/`.*
- **Incorrect Example:** *The core declares `SyncResponse` in `domain/`.*

## Technical Name: ValidationError
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The domain error that names a rejected cell value.
- **Approved Form:** ValidationError
- **Do Not Use:** CellError, WriteError, BadValueError
- **Correct Example:** *The core declares `ValidationError` in `domain/`.*
- **Incorrect Example:** *The core declares `CellError` in `domain/`.*

## Technical Name: kibtab
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The name of the project. The module path uses it.
- **Approved Form:** kibtab, Kibtab
- **Do Not Use:** KibTab, Kibtab Engine, Kibtab Server
- **Correct Example:** *The module path is `github.com/kibtab/kibtab`.*
- **Incorrect Example:** *The module path is `github.com/kibtab/KibTab`.*

## Technical Name: COVERAGE_FLOOR
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The variable in `Makefile` that holds the coverage floor.
- **Approved Form:** COVERAGE_FLOOR
- **Do Not Use:** COVERAGE_MIN, COVERAGE_TARGET, COVERAGE_LIMIT
- **Correct Example:** *`Makefile` sets `COVERAGE_FLOOR` to 100.*
- **Incorrect Example:** *`Makefile` sets `COVERAGE_MIN` to 100.*

## Technical Name: CGO_ENABLED
- **Part of Speech:** Noun
- **Category:** Identifier
- **Definition:** The environment variable that turns CGO off.
- **Approved Form:** CGO_ENABLED
- **Do Not Use:** GOCGO, CGO, DISABLE_CGO
- **Correct Example:** *Every build sets `CGO_ENABLED=0`.*
- **Incorrect Example:** *Every build sets `DISABLE_CGO=0`.*

## Next Steps

* Read [the architecture lexicon](architecture-terms.md) for the layer terms.
* Read [the tooling lexicon](tooling-terms.md) for the command terms.