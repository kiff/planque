# Google Docs

Checked: 2026-11-14 · Sources: the connection service's Google Docs tool
definitions as planque reads them (46 tools, [generated](generated/googledocs.md)),
and Google's Docs API reference.
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **A sign-in to Google**, held by the connection service. The grant's catalog lists
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
| `GOOGLEDOCS_SEARCH_DOCUMENTS` | Finds documents by name or date; **needs Drive access** | read | document ids |
| `GOOGLEDOCS_EXPORT_DOCUMENT_AS_PDF` | Exports a document as PDF; **needs Drive access** | read | |
| `GOOGLEDOCS_CREATE_DOCUMENT` | Creates a document with a title and text | write | **the document's id** |
| `GOOGLEDOCS_CREATE_DOCUMENT_MARKDOWN` | Creates a document from formatted text | write | **the document's id** |
| `GOOGLEDOCS_COPY_DOCUMENT` | Copies a document, such as a template; **needs Drive access** | write | the copy's id |
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
- `GOOGLEDOCS_UPLOAD_DRIVE_FILE`, `GOOGLEDOCS_CREATE_DRIVE_FILE` and the
  Drive and Sheets reads in this app:
  use the Drive and Sheets apps, whose guides cover them.

## Access and plans

- **What Google's page asks for:** the narrowest access the managed sign-in
  offers for the taf's tools (Google's words, roughly, in quotes), which can
  be wider than the job:
  - **reading or writing a document:** Docs access, up to edit ("see, edit,
    create and delete all your documents");
  - **finding a document by name, exporting it as PDF, reading its file
    details, copying a template, or creating or uploading a file:** Drive
    access as well ("see, edit, create and delete all your Drive files").
    Offer these only when the job needs them.
  - The tools are named in the connect: a connect that names none asks for the
    app's widest access, so the taf's planned tools are always named when it
    connects.
- **The approval shows how far Google's sign-in reaches,** from planque's
  grant (`usual_access`, its `scopes`). Coach doesn't restate it, at any
  step: not in a question, its options, the plan or a message.
- **No paid plan is needed.** A business admin can block outside apps.
- **Size:** a document can hold about 1.02 million characters. A taf that
  appends a line every day for years should start a new document each year.
- **Rate:** Google allows about 60 writes a minute per user. A job writing
  many small edits should spread them out.
- **Suggestions and comments** made by people aren't changed by a replace,
  and this app has no tool to accept or reject them.
- **Copying the owner's template** needs Drive access to that template: on
  access to the taf's own files only, the copy is refused.

## When it fails

The rows are Google Docs's own answers: a refusal reaches the taf as
"Google Docs refused it (it answered …): …", with Google Docs's reason.
planque's own refusals, and an expired sign-in, are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | Most often a search, which needs Drive access | "Google Docs is only partly connected: finding documents by name needs access you didn't give. Send me the document's link instead." |
| Refused (404), document not found | Deleted, or not shared with this account | "I can't open that document. It may have been deleted or not shared with this account." |
| Refused (400), bad write | A position past the end, after an edit moved it | Read the document again and retry once; tell the owner only if it fails again. |
| Refused (429), too many calls | Over Google's per-minute limit | "Google Docs asked us to slow down. The taf will continue in a minute." |
| Google Docs's sign-in has expired (401) | Removed, or expired | "Google Docs needs you to sign in again. Nothing was changed." |

A failed search never leads to creating a new document in its place.
