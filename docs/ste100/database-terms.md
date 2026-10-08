# STE100 Technical Names: database terms

Version 0.0.0-docs. This document holds the database Technical Names.

These terms name the tables, rows, and writes of a relational database.
Read [the catalogue](index.md) first.

## Technical Name: Row
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** One record in a table. A row holds one value for each field.
- **Approved Form:** Row (singular), Rows (plural)
- **Do Not Use:** Record, Entry, Item, Tuple
- **Correct Example:** *The repository writes one **row** for each change.*
- **Incorrect Example:** *The repository writes one record for each change.*

## Technical Name: Table
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** A named set of rows that share the same fields.
- **Approved Form:** Table (singular), Tables (plural)
- **Do Not Use:** Dataset, Grid, Sheet, Relation
- **Correct Example:** *The registry lists every **table** in the database.*
- **Incorrect Example:** *The registry lists every dataset in the database.*

## Technical Name: Field
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** One named value in a row. A column holds the fields of a table.
- **Approved Form:** Field (singular), Fields (plural)
- **Do Not Use:** Column, Property, Attribute, Cell
- **Correct Example:** *The schema names each **field** with its type.*
- **Incorrect Example:** *The schema names each column with its type.*

## Technical Name: Primary Key
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The field that gives each row in a table a unique name.
- **Approved Form:** Primary Key (singular), Primary Keys (plural)
- **Do Not Use:** ID, Key, Index, Identifier
- **Correct Example:** *The repository reads the **primary key** of each row.*
- **Incorrect Example:** *The repository reads the key of each row.*

## Technical Name: Transaction
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** A group of writes that the database runs as one unit.
- **Approved Form:** Transaction (singular), Transactions (plural)
- **Do Not Use:** Batch, Group, Sequence, Operation
- **Correct Example:** *The instance runs the writes in one **transaction**.*
- **Incorrect Example:** *The instance runs the writes in one batch.*

## Technical Name: Commit
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The action that saves the writes of a transaction.
- **Approved Form:** Commit (singular), Commits (plural)
- **Do Not Use:** Save, Store, Apply, Flush
- **Correct Example:** *The instance makes a **commit** after the audit row.*
- **Incorrect Example:** *The instance makes a save after the audit row.*

## Technical Name: Rollback
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The action that undoes the writes of a failed transaction.
- **Approved Form:** Rollback (singular), Rollbacks (plural)
- **Do Not Use:** Undo, Revert, Cancel, Reverse
- **Correct Example:** *The instance makes a **rollback** when one write fails.*
- **Incorrect Example:** *The instance makes an undo when one write fails.*

## Technical Name: Engine
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The database product that Kibtab connects to.
  It names only the product. It never names the Kibtab process.
- **Approved Form:** Engine (singular), Engines (plural)
- **Do Not Use:** Database, Server, Backend, Provider, Instance
- **Correct Example:** *Each **engine** runs in its own adapter package.*
- **Incorrect Example:** *Each database runs in its own adapter package.*

## Technical Name: ValidationError
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** A domain error that names a rejected cell value.
- **Approved Form:** ValidationError (singular), ValidationErrors (plural)
- **Do Not Use:** CellError, WriteError, BadValue
- **Correct Example:** *The service returns a **ValidationError** for an invalid cell.*
- **Incorrect Example:** *The service returns a cell error for an invalid cell.*

## Technical Name: Dialect
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The rules of one engine for quoting names and paging a query.
- **Approved Form:** Dialect (singular), Dialects (plural)
- **Do Not Use:** Flavour, Variant, Profile, Mode
- **Correct Example:** *The adapter asks the **dialect** to quote each name.*
- **Incorrect Example:** *The adapter asks the flavour to quote each name.*

## Technical Name: Audit Row
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** One row that records a change to a cell and the time of it.
- **Approved Form:** Audit Row (singular), Audit Rows (plural)
- **Do Not Use:** History Row, Log Entry, Journal Entry, Trace
- **Correct Example:** *The instance writes one **audit row** for each cell change.*
- **Incorrect Example:** *The instance writes one history row for each cell change.*

## Technical Name: Delta
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The set of cell changes that a client sends to the instance.
- **Approved Form:** Delta (singular), Deltas (plural)
- **Do Not Use:** Batch, Change Set, Patch, Mutation
- **Correct Example:** *The instance validates each **delta** before it writes.*
- **Incorrect Example:** *The instance validates each change set before it writes.*

## Technical Name: Cell Value
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** One stored value in a cell. It belongs to a row and a field.
- **Approved Form:** Cell Value (singular), Cell Values (plural)
- **Do Not Use:** Field Value, Column Value, Payload Value
- **Correct Example:** *The repository returns each **cell value** as a string.*
- **Incorrect Example:** *The repository returns each field value as a string.*

## Technical Name: Table Metadata
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The name, field list, and current version of one table.
- **Approved Form:** Table Metadata (singular), Table Metadata (plural)
- **Do Not Use:** Table Schema, Table Description, Table Info
- **Correct Example:** *The registry returns the **table metadata** for one table.*
- **Incorrect Example:** *The registry returns the table schema for one table.*

## Technical Name: Sync Result
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The result that the instance returns after a sync.
- **Approved Form:** Sync Result (singular), Sync Results (plural)
- **Do Not Use:** Sync Reply, Sync Response, Sync Status
- **Correct Example:** *The sync result holds the new **version** for each row.*
- **Incorrect Example:** *The sync reply holds the new revision for each row.*

## Next Steps

* Read [the spreadsheet lexicon](spreadsheet-terms.md) for the cell terms.
* Read [the architecture lexicon](architecture-terms.md) for the layer terms.