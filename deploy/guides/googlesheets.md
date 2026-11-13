# Google Sheets

Checked: 2026-11-13 · Sources: the connection service's Google Sheets tool
definitions as planque reads them (58 tools), Google's Sheets and Drive API
references, and test runs on a phone on 13 Nov 2026.

## How it connects

- **A sign-in to Google**, held by the connection service. planque's catalog
  lists the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  A tool outside it means Google asks the owner again ("already has some
  access"), so ask for every tool the job needs in the first request.
- **One approval for the sign-in, one for the tools and the Card.**

## Two kinds of access, and why it matters

A Sheets sign-in usually gives **spreadsheet access**: read and write the
cells of sheets the taf is pointed at, by their id. **Finding a sheet by
name needs Drive access**, which a Sheets sign-in often doesn't give.
On 13 Nov 2026 a taf's search was refused for exactly this, and it would
have made a new sheet every week instead of finding its own.

So:

- **A sheet the taf creates:** keep its id from `GOOGLESHEETS_CREATE_GOOGLE_SHEET1`
  and use it every run. Never search for it by name.
- **A sheet the owner already has:** ask the owner for its link once; the id
  is in the link. Search only if the catalog offers the search tool.
- **A failed search never leads to a create.** The taf says it couldn't find
  the sheet and stops.

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `GOOGLESHEETS_VALUES_GET` | Reads the cells of one range, a page of rows at a time | read | the range read, its values |
| `GOOGLESHEETS_BATCH_GET` | Reads several ranges at once | read | each range's values |
| `GOOGLESHEETS_LOOKUP_SPREADSHEET_ROW` | Finds the first row matching a value | read | the row and its number |
| `GOOGLESHEETS_GET_SHEET_NAMES` | Lists the tabs of a spreadsheet | read | tab names |
| `GOOGLESHEETS_GET_SPREADSHEET_INFO` | Reads a spreadsheet's title, tabs and settings | read | the spreadsheet id, tab ids |
| `GOOGLESHEETS_SEARCH_SPREADSHEETS` | Finds spreadsheets by name or date | read, **needs Drive access** | spreadsheet ids and links |
| `GOOGLESHEETS_CREATE_GOOGLE_SHEET1` | Creates a new spreadsheet | write | **the spreadsheet id and its link**: keep them |
| `GOOGLESHEETS_VALUES_UPDATE` | Writes values into one range, replacing what's there | write | the range written and how many cells |
| `GOOGLESHEETS_UPDATE_VALUES_BATCH` | Writes several ranges at once | write | the ranges written |
| `GOOGLESHEETS_SPREADSHEETS_VALUES_APPEND` | Adds rows after the last row of a table | write | the range the rows went to |
| `GOOGLESHEETS_UPSERT_ROWS` | Updates rows by a key column, and appends new ones | write | rows updated and added |
| `GOOGLESHEETS_ADD_SHEET` | Adds a tab | write | the new tab's id |

**What to keep:** a created spreadsheet's id and link; a tab's name or id.
A write answers the range it changed and how many cells; that answer is the
proof the write landed, so the run's report should carry it.

**Prefer** `VALUES_UPDATE` on a known range, or `VALUES_APPEND`, over tools
that clear, delete or sort. A job that only fills a column needs no other
write.

## Don't offer

- Deleting or clearing tools (`GOOGLESHEETS_DELETE_SHEET`,
  `GOOGLESHEETS_DELETE_DIMENSION`, `GOOGLESHEETS_CLEAR_VALUES` and the batch
  clears) unless the job is clearly about removing data.
- `GOOGLESHEETS_FIND_REPLACE`: it can change cells across a whole sheet.
- Moving a spreadsheet between folders, or the Drive change feed: they act
  on the owner's Drive, beyond the sheet.
- Deprecated tools (`GOOGLESHEETS_EXECUTE_SQL`, `GOOGLESHEETS_BATCH_UPDATE`,
  `GOOGLESHEETS_QUERY_TABLE`, `GOOGLESHEETS_SHEET_FROM_JSON` and others marked
  so): use the tools above.

## Access and plans

- **No paid plan is needed**; any Google account works. A business admin can
  block outside apps.
- **Rate:** Google allows about 60 reads and 60 writes a minute for one user,
  and caps each project too. A run that writes row by row over hundreds of
  rows should write in batches.
- **Size:** a new tab has 1,000 rows and 26 columns (A to Z). Writing past
  column Z needs the sheet widened first.
- **Tab names with spaces** must be quoted in a range (`'Leads 2026'!A1:E1`),
  or Google refuses the range.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | The sign-in didn't give this tool's access: most often a search, which needs Drive | "Google Sheets is only partly connected: finding sheets by name needs access you didn't give. Send me the sheet's link instead, or sign in again and allow it." |
| Too many calls (429) | Over Google's per-minute limit | "Google Sheets asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | The owner removed access or the sign-in expired | "Google Sheets needs you to sign in again. Nothing was changed." |
| The sheet isn't found (404) | Deleted, or the taf isn't allowed to open it | "I can't open that sheet. It may have been deleted or not shared with this account." |
| A range is refused (400) | A tab name without quotes, or a range outside the sheet | Fix the range and try once; don't tell the owner unless it fails again. |
| The owner doesn't see a write | Google's phone app shows an open sheet's old cells until it's reopened | "It's written. Close and reopen the sheet to see it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create or a write.
