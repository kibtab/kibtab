# STE100 Technical Names: spreadsheet terms

Version 0.0.0-docs. This document holds the spreadsheet Technical Names.

These terms name the cells, ranges, and workbooks of a spreadsheet client.
Read [the catalogue](index.md) first.

## Technical Name: Cell
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** The smallest box in a sheet that holds one value.
- **Approved Form:** Cell (singular), Cells (plural)
- **Do Not Use:** Box, Element, Field, Slot
- **Correct Example:** *The instance writes the value into one **cell**.*
- **Incorrect Example:** *The instance writes the value into one box.*

## Technical Name: Range
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** A set of cells that forms a rectangle in a sheet.
- **Approved Form:** Range (singular), Ranges (plural)
- **Do Not Use:** Block, Area, Selection, Cluster
- **Correct Example:** *The client sends the **range** that the user changed.*
- **Incorrect Example:** *The client sends the block that the user changed.*

## Technical Name: Sheet
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** One page of cells in a workbook. A sheet holds rows and columns.
- **Approved Form:** Sheet (singular), Sheets (plural)
- **Do Not Use:** Tab, Page, Worksheet, Document
- **Correct Example:** *Each **sheet** maps to one table in the database.*
- **Incorrect Example:** *Each tab maps to one table in the database.*

## Technical Name: Workbook
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** The file that holds every sheet for one user.
- **Approved Form:** Workbook (singular), Workbooks (plural)
- **Do Not Use:** File, Document, Spreadsheet, Book
- **Correct Example:** *The **workbook** opens on the server with the user name.*
- **Incorrect Example:** *The file opens on the server with the user name.*

## Technical Name: Column Header
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** The text at the top of a column. It names the field.
- **Approved Form:** Column Header (singular), Column Headers (plural)
- **Do Not Use:** Header, Caption, Label, Title
- **Correct Example:** *The instance reads the **column header** to find the field.*
- **Incorrect Example:** *The instance reads the caption to find the field.*

## Technical Name: Row Header
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** The text at the left of a row. It names the record key.
- **Approved Form:** Row Header (singular), Row Headers (plural)
- **Do Not Use:** Header, Key Label, Index, Marker
- **Correct Example:** *The instance reads the **row header** to find the key.*
- **Incorrect Example:** *The instance reads the marker to find the key.*

## Technical Name: Delta
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** The difference between two versions of a range.
- **Approved Form:** Delta (singular), Deltas (plural)
- **Do Not Use:** Change, Difference, Patch, Edit
- **Correct Example:** *The instance rejects a **delta** from an old version.*
- **Incorrect Example:** *The instance rejects a change from an old version.*

## Technical Name: Client
- **Part of Speech:** Noun
- **Category:** Spreadsheet Term
- **Definition:** The program that runs in the spreadsheet. It sends a delta.
- **Approved Form:** Client (singular), Clients (plural)
- **Do Not Use:** Frontend, Add-in, Extension, Application
- **Correct Example:** *Each **client** runs in its own spreadsheet.*
- **Incorrect Example:** *Each frontend runs in its own spreadsheet.*

## Next Steps

* Read [the database lexicon](database-terms.md) for the table terms.
* Read [the architecture lexicon](architecture-terms.md) for the layer terms.