# Google Docs

Checked: 2026-11-13 · Sources: the connection service's Google Docs tool
definitions as planque reads them (46 tools, [generated](generated/googledocs.md)),
and Google's Docs API reference.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **One approval for the sign-in, one for the tools and the Card.**
- **Finding a document by name needs Drive access** that a Docs sign-in may
  not give, as with Sheets. Keep the id of each document the taf creates, and
  ask the owner for a link to theirs.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLEDOCS_GET_DOCUMENT_PLAINTEXT` | Reads a document as plain text | read | |
| `GOOGLEDOCS_GET_DOCUMENT_BY_ID` | Reads a document with its structure | read | its revision |
| `GOOGLEDOCS_SEARCH_DOCUMENTS` | Finds documents by name or date | read, **needs Drive access** | document ids |
| `GOOGLEDOCS_EXPORT_DOCUMENT_AS_PDF` | Exports a document as PDF | read | |
| `GOOGLEDOCS_CREATE_DOCUMENT` | Creates a document with a title and text | write | **the document's id** |
| `GOOGLEDOCS_CREATE_DOCUMENT_MARKDOWN` | Creates a document from formatted text | write | **the document's id** |
| `GOOGLEDOCS_COPY_DOCUMENT` | Copies a document, such as a template | write | the copy's id |
| `GOOGLEDOCS_REPLACE_ALL_TEXT` | Replaces every match of a phrase | write | |
| `GOOGLEDOCS_INSERT_TEXT_ACTION` | Inserts text at a place in the document | write | |
| `GOOGLEDOCS_UPDATE_DOCUMENT_MARKDOWN` | Replaces a document's whole content | write | |

**Templates work best as copy then replace:** copy the owner's template,
then replace its placeholders (`{{client}}`) in the copy. The template itself
is never changed.

**Writing by position is fragile.** A place in a document is counted in
characters, and it shifts with every edit. Read the document again before
each write that depends on a position, or replace by phrase instead.

## Don't offer

- `GOOGLEDOCS_UPDATE_DOCUMENT_MARKDOWN` on a document the owner wrote: it
  replaces everything. Offer it only for documents the taf made.
- Deleting content ranges, tabs, headers or tables in the owner's own
  documents, unless the job is clearly about editing them.
- `GOOGLEDOCS_UPLOAD_DRIVE_FILE` and the Drive and Sheets reads in this app:
  use the Drive and Sheets apps, whose guides cover them.

## Access and plans

- **No paid plan is needed.** A business admin can block outside apps.
- **Size:** a document can hold about 1.02 million characters. A taf that
  appends a line every day for years should start a new document each year.
- **Rate:** Google allows about 60 writes a minute per user. A job writing
  many small edits should put them in one write.
- **Suggestions and comments** made by people aren't changed by a replace;
  a taf can't accept or reject them.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | Most often a search, which needs Drive access | "Google Docs is only partly connected: finding documents by name needs access you didn't give. Send me the document's link instead." |
| The document isn't found (404) | Deleted, or not shared with this account | "I can't open that document. It may have been deleted or not shared with this account." |
| A write is refused (400) | A position past the end, after an edit moved it | Read the document again and retry once; tell the owner only if it fails again. |
| Too many calls (429) | Over Google's per-minute limit | "Google Docs asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, or expired | "Google Docs needs you to sign in again. Nothing was changed." |

After any refusal the taf stops that step and says so. A failed search never
leads to creating a new document in its place.
