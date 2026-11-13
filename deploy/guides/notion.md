# Notion

Checked: 2026-11-13 · Sources: the connection service's Notion tool
definitions as planque reads them (56 tools), and Notion's API reference
(request limits, status codes, authorization).

## How it connects

- **A sign-in to Notion**, held by the connection service; the taf never sees
  a password or token. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **The owner chooses which pages the taf can see.** At sign-in Notion shows
  a page picker. Only the pages and databases picked there (and what sits
  under them) are reachable. Anything else answers as if it didn't exist.
  The owner can add pages later from Notion's connection settings, with no
  new sign-in.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in the workspace. Notion's catalog entries
give no `returns` paths today; the ids below are those Notion's answers
carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `NOTION_SEARCH_NOTION_PAGE` | Finds pages and databases shared with the taf, by title | read | page and database ids, a cursor |
| `NOTION_FETCH_DATABASE` | Reads a database's columns and their types | read | the database id, its property ids |
| `NOTION_QUERY_DATABASE_WITH_FILTER` | Lists a database's rows that match a filter, sorted | read | row (page) ids, a cursor |
| `NOTION_FETCH_ROW` | Reads one database row | read | the row's id, its properties |
| `NOTION_RETRIEVE_PAGE` | Reads a page's properties | read | the page id, its url |
| `NOTION_GET_PAGE_MARKDOWN` | Reads a page's content as text | read | the page's text |
| `NOTION_FETCH_COMMENTS` | Lists the comments on a page or block | read | comment ids, discussion ids |
| `NOTION_INSERT_ROW_DATABASE` | Adds a row to a database | write | the new row's id and url |
| `NOTION_UPDATE_ROW_DATABASE` | Changes a row's properties | write | the row's id |
| `NOTION_CREATE_NOTION_PAGE` | Creates a page under a page the taf can reach | write | the new page's id and url |
| `NOTION_ADD_MULTIPLE_PAGE_CONTENT` | Adds paragraphs, headings and lists to a page | write | the new blocks' ids |
| `NOTION_CREATE_COMMENT` | Leaves a comment on a page or in a discussion | write | the comment id |

**What to keep:** the id of a page or row the taf created, to update that
same one later; a database's id and its property ids, so filters and new
rows use the right columns. Never search by title for something the taf
made: titles change and repeat. Use the id it kept.

**Start read-only.** "Summarise what changed in this database" needs only
the read tools. For writes, prefer adding a row or a comment over editing
what a person wrote.

## Don't offer

- Deleting and archiving (`NOTION_DELETE_BLOCK`, `NOTION_ARCHIVE_NOTION_PAGE`,
  and `NOTION_UPDATE_ROW_DATABASE` with its delete option): they take
  content away from the team.
- Replacing a page's content (`NOTION_UPDATE_PAGE_MARKDOWN` in replace mode):
  it overwrites what people wrote.
- Changing a database's columns or views (`NOTION_UPDATE_SCHEMA_DATABASE`,
  `NOTION_CREATE_VIEW`, `NOTION_UPDATE_VIEW`, `NOTION_DELETE_VIEW`): it
  changes the table for everyone who uses it.
- Moving pages (`NOTION_MOVE_PAGE`) and listing the workspace's members
  (`NOTION_LIST_USERS`), unless the job is about exactly that.
- The tools that turn a sentence into a row (`NOTION_INSERT_ROW_FROM_NL`):
  the taf can't show the owner exactly what will be written before it is.

## Access and plans

- **What the taf sees is what the owner shared** at sign-in, not the whole
  workspace. A job that needs a new page or database needs it shared first.
- **The sign-in's abilities** (read content, insert, update, comment) are
  set by Notion's connection. A tool outside them comes under `unusable`.
- **Rate:** Notion gives each connection a budget per minute, set by the
  workspace's plan: 180 calls a minute on most plans, 600 on its Business
  and Enterprise plans. It also limits the whole workspace across all its
  connections, scaled to the plan. Querying a
  database has its own limit of about three a second. A job that reads
  hundreds of rows a run should page and pause.
- **Size per request:** at most 100 blocks in one append, about 1,000
  elements and 500 KB in all, and 2,000 characters in one piece of text.
  Long content goes in several calls.
- **A workspace's plan** can cap how much content it holds. When the cap is
  reached, writes are refused until the owner upgrades or frees space.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Not found (404, `object_not_found`) | The page or database wasn't shared with the taf at sign-in, or it was deleted | "The taf can't see that Notion page. Share it with the taf from the page's connections menu, then try again." |
| Refused, not allowed (403, `restricted_resource`) | The sign-in can't do this kind of change, or the workspace's content cap is reached | "Notion didn't let the taf make this change. Check the taf's access to that page, or your workspace's plan." |
| Too many calls (429, `rate_limited`) | Over the connection's or the workspace's per-minute budget | "Notion asked us to slow down. The taf will try again in a minute." |
| Blocked by Notion (429, `public_api_request_blocked`) | Notion has restricted this connection | "Notion has blocked this connection. Only Notion's support can lift that." |
| The sign-in no longer works (401) | The owner removed the connection, or its access was revoked | "Notion needs you to sign in again. Nothing was changed." |
| Rejected, wrong shape (400, `validation_error`) | A value doesn't fit the column, or the content is too long for one call | "Notion didn't accept that entry: a value doesn't fit its column. The taf will skip it and tell you which." |
| Conflict (409) | Two edits to the same thing at once | "Someone was editing the same page. The taf will try once more." |
| Notion unavailable (503) | Notion is down, or the answer took too long; a write may have saved anyway | "Notion isn't answering right now. The taf will check whether the change was saved before trying again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create or a write: if it can't find its page, it asks, rather
than making another one.
