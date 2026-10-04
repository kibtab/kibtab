# STE100 Technical Names: database terms

Version 0.0.0-docs. This document holds the database Technical Names.

These terms name the tables, rows, and writes of a relational database.
Read the [STE100 Technical Names](index.md) dictionary first.

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
- **Correct Example:** *The engine runs the writes in one **transaction**.*
- **Incorrect Example:** *The engine runs the writes in one batch.*

## Technical Name: Commit
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The action that saves the writes of a transaction.
- **Approved Form:** Commit (singular), Commits (plural)
- **Do Not Use:** Save, Store, Apply, Flush
- **Correct Example:** *The engine makes a **commit** after the audit row.*
- **Incorrect Example:** *The engine makes a save after the audit row.*

## Technical Name: Rollback
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The action that undoes the writes of a failed transaction.
- **Approved Form:** Rollback (singular), Rollbacks (plural)
- **Do Not Use:** Undo, Revert, Cancel, Reverse
- **Correct Example:** *The engine makes a **rollback** when one write fails.*
- **Incorrect Example:** *The engine makes an undo when one write fails.*

## Technical Name: Engine
- **Part of Speech:** Noun
- **Category:** Database Term
- **Definition:** The database product that Kibtab connects to.
- **Approved Form:** Engine (singular), Engines (plural)
- **Do Not Use:** Database, Server, Backend, Provider
- **Correct Example:** *Each **engine** runs in its own adapter package.*
- **Incorrect Example:** *Each database runs in its own adapter package.*

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
- **Correct Example:** *The engine writes one **audit row** for each cell change.*
- **Incorrect Example:** *The engine writes one history row for each cell change.*

## Next Steps

* Read [spreadsheet-terms.md](spreadsheet-terms.md) for the cell terms.
* Read [architecture-terms.md](architecture-terms.md) for the layer terms.