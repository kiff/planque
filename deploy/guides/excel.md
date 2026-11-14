# Excel

Checked: 2026-11-14 · Sources: the connection service's Excel tool
definitions as planque reads them (125 tools, [generated](generated/excel.md)),
and Microsoft's workbook API reference.

## How it connects

- **A sign-in to Microsoft**, held by the connection service: a work or
  school account. Workbooks in a personal OneDrive can't be opened. The catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`.
  Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **A work account's admin** may have to approve the app first.
- **Only workbooks in OneDrive or SharePoint** can be read and written. A file
  on the owner's computer has to be uploaded first.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `EXCEL_SEARCH_FILES` / `EXCEL_LIST_FILES` | Finds workbooks by name, or lists a folder | read | **the workbook's item id** |
| `EXCEL_LIST_WORKSHEETS` | Lists a workbook's sheets | read | sheet names |
| `EXCEL_GET_WORKSHEET_USED_RANGE` | Reads every filled cell of a sheet | read | |
| `EXCEL_GET_RANGE` | Reads one range (`A1:E20`) | read | |
| `EXCEL_LIST_TABLES` / `EXCEL_LIST_TABLE_ROWS` | Lists a sheet's tables, and a table's rows | read | table names |
| `EXCEL_CREATE_WORKBOOK` | Creates a workbook | write | **its item id** |
| `EXCEL_UPDATE_RANGE` | Writes values into one range | write | |
| `EXCEL_ADD_TABLE_ROW` | Adds rows at the end of a table | write | |
| `EXCEL_ADD_WORKSHEET` | Adds a sheet | write | |

**What to keep:** the workbook's item id, from the first run. Searching by
name every time finds the wrong file the day a copy exists.

**Prefer tables.** Adding a row to a table (`EXCEL_ADD_TABLE_ROW`) never
writes over what's there; writing a range can.

## Don't offer

- Sharing and permissions (`EXCEL_CREATE_WORKBOOK_SHARING_LINK`,
  `EXCEL_ADD_WORKBOOK_PERMISSION`, `EXCEL_REVOKE_WORKBOOK_PERMISSION`).
- Deleting files, sheets, tables or ranges (`EXCEL_DELETE_DRIVE_ITEM`,
  `EXCEL_DELETE_WORKSHEET`, `EXCEL_DELETE_TABLE`, `EXCEL_DELETE_RANGE`,
  `EXCEL_CLEAR_RANGE` and the others) unless the job is about removing data.
- Protecting or unprotecting sheets, and the raw batch tools
  (`EXCEL_RUN_GRAPH_MUTATION_BATCH`).

## Access and plans

- **A work or school account is needed.** The workbook must be in the new format
  (`.xlsx`); old `.xls` files can't be opened by the API.
- **Size:** very large ranges can fail; read and write big data in smaller
  ranges.
- **Rate:** Microsoft limits calls per workbook; writing row by row over
  hundreds of rows should be done as one range or a few table rows at a time.
- **A workbook open elsewhere** is still written, and people see the change
  when the file syncs.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in gave read access only, or the file is someone else's | "Excel is only partly connected, or that workbook isn't yours to change. Sign in again and allow editing, or pick another file." |
| The file isn't found (404) | Moved, deleted, or not in OneDrive | "I can't find that workbook in your OneDrive. It may have moved; send me its link." |
| The file is locked (409 or 423) | Someone is changing it right now, or it's checked out | "The workbook is busy. The taf will try again in a minute." |
| The range is refused (400) | A sheet name that doesn't exist, or a range shape that doesn't match the values | Fix it and try once; tell the owner only if it fails again. |
| Too many calls (429) | Over Microsoft's limit | "Microsoft asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, password changed, or expired | "Excel needs you to sign in again. Nothing was changed." |

After any refusal the taf stops that step and says so. A failed search never
leads to creating a new workbook.
