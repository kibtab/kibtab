# STE100 Technical Names for Kibtab

Version 0.0.0-docs. This document holds the controlled list of Technical Names.

ASD-STE100 Simplified Technical English defines a controlled vocabulary.
Most words in that vocabulary are standard STE words.
Some Kibtab concepts have no accurate standard STE word.
For each of those concepts, Kibtab approves a Technical Name.

A Technical Name is a non-STE word that Kibtab explicitly approves.
Kibtab approves it because no standard STE word has the same precision.
A Technical Name has one mandatory part of speech and a defined scope.
The scope is the software domain of Kibtab.

Use a Technical Name as defined here in every document, comment, and commit.
Add a new Technical Name to this catalogue before you use it.

Read the vendored skill at `../../skills/simple-english/SKILL.md` first.

## The Lexicons

Each category has its own lexicon page.

| Lexicon | Holds |
| --- | --- |
| [spreadsheet-terms.md](spreadsheet-terms.md) | The cells, ranges, and workbook terms. |
| [database-terms.md](database-terms.md) | The tables, rows, columns, and write terms. |
| [architecture-terms.md](architecture-terms.md) | The layer, port, adapter, and domain terms. |
| [tooling-terms.md](tooling-terms.md) | The commands, targets, files, and gate terms. |
| [identifiers.md](identifiers.md) | The exact names in the Kibtab source. |

## The Fields Of An Entry

Every Technical Name has all seven fields.

| Field | Holds |
| --- | --- |
| Part of Speech | Noun, Verb, or Modifier. |
| Category | The lexicon that holds the entry. |
| Definition | The meaning, restricted to the Kibtab domain. |
| Approved Form | The singular, the plural, and the abbreviation. |
| Do Not Use | The words that an author must not use instead. |
| Correct Example | A sentence that uses the name correctly. |
| Incorrect Example | A sentence that uses a rejected word. |

Keep every field in every entry.
A missing field fails the catalogue check.

## The Rules For A Name

* Keep a name a Noun or a Modifier. Never approve a Verb.
  Use a standard STE verb for the action instead.
  A database word that names an action is a Noun here.
  Use the verb *Make* with that noun.
* Use one name for one concept. Never approve a second name for it.
* Keep each name once in one lexicon. Never repeat it in a second lexicon.
* Approve a concept and its source name in two lexicons.
  The [identifiers.md](identifiers.md) entry names the exact case.
  Each other lexicon names the concept.

## A Known Ambiguity

The word engine names two things in older Kibtab prose.
It names the database product in the port and adapter rules.
It names the running copy of Kibtab in the install and self-hosting guides.

Use **Engine** only for the database product.
Use **Instance** for the running copy of Kibtab.
The newer guides follow this rule.
The older prose still uses engine for both senses.

## What The Check Does

`make docs-check` checks the shape of this catalogue.
The check fails when an entry lacks a field, uses a rejected part of speech,
names an unknown category, repeats a name, or sits outside the index.

The check does not read the Do Not Use field as a banned word list.
Each rejected word names a synonym for one concept.
Most of those words stay correct in another context.
A review reads the Do Not Use field for the concept that it names.
* Use the exact case of the name as it appears in this catalogue.
* Use the plural form as the Approved Form states.
* Never invent a name in a document. Add it here first.
* Keep each scope inside the Kibtab domain.

## The Template

Use this template for a new entry.

```markdown
## Technical Name: Sync
- **Part of Speech:** Noun
- **Category:** Architecture Term
- **Definition:** The flow that sends a client change to the database.
- **Approved Form:** Sync (singular), Syncs (plural)
- **Do Not Use:** Push, Upload, Copy
- **Correct Example:** *The engine runs a **sync** when a cell changes.*
- **Incorrect Example:** *The engine runs an upload when a cell changes.*
```

## The Difference From The Word Lists

The vendored skill holds the STE word rules.
Kibtab does not keep a copy of the ASD-STE100 word lists.
The standard states that no reproduction is allowed without written authority.

This catalogue holds the opposite kind of data.
Each entry names a word that the standard does not define.
Kibtab writes each entry from its own domain.
No entry copies a word from the ASD-STE100 word lists.

## Next Steps

* Read [spreadsheet-terms.md](spreadsheet-terms.md) for the cell terms.
* Read [database-terms.md](database-terms.md) for the table terms.
* Read [architecture-terms.md](architecture-terms.md) for the layer terms.
* Read [tooling-terms.md](tooling-terms.md) for the command terms.
* Read [identifiers.md](identifiers.md) for the source names.