# OneDrive

Checked: 2026-11-13 · Sources: the connection service's OneDrive tool
definitions as planque reads them (101 tools), Microsoft's OneDrive and
Graph API references (files, uploads, workbooks, error codes) and
Microsoft's app consent documentation.

## How it connects

- **A sign-in to Microsoft**, held by the connection service; the taf never
  sees a password or token. A personal or a work or school account both
  sign in; some tools work only on the second (see below).
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  A sign-in made for reading doesn't let the taf write, so ask for every tool
  the job needs in the first request.
- **A work or school account answers to its organisation.** By default a
  person can allow an app into their own files, but wider access, such as
  every file in the organisation, needs an admin's approval, and an admin can
  turn off people's own approvals altogether.

## Tools a job usually needs

Read means it changes nothing in the drive.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ONE_DRIVE_ONEDRIVE_FIND_FILE` | Finds files by name, in a folder or the whole drive | read | item ids, names, a cursor |
| `ONE_DRIVE_SEARCH_ITEMS` | Searches names and content across the drive | read | item ids, names, web links, a page token |
| `ONE_DRIVE_LIST_FOLDER_CHILDREN` | Lists what's in one folder, by id or path | read | item ids, names, a page token |
| `ONE_DRIVE_GET_ITEM` / `ONE_DRIVE_GET_ITEM_BY_PATH` | Reads one item's details | read | the item id, drive id, eTag, web link |
| `ONE_DRIVE_DOWNLOAD_FILE` | Downloads one file's content | read | the file's content |
| `ONE_DRIVE_LIST_WORKBOOK_WORKSHEETS` | Lists the sheets of a workbook (.xlsx) | read | sheet ids and names |
| `ONE_DRIVE_READ_WORKBOOK_RANGE` | Reads cells of a sheet, such as A1:F50 | read | the address, values, formulas |
| `ONE_DRIVE_ONEDRIVE_CREATE_FOLDER` | Makes a folder | write | the new folder's id |
| `ONE_DRIVE_ONEDRIVE_UPLOAD_FILE` | Uploads a new file into a folder | write | the new file's id, web link |
| `ONE_DRIVE_UPDATE_FILE_CONTENT` | Replaces a file's content with a new version | write | the item id, new eTag |
| `ONE_DRIVE_UPDATE_WORKBOOK_RANGE` | Writes values or formulas into cells | write | the address, values written |

**What to keep:** a file's item id and drive id, to reach the same file
next time instead of searching; its eTag, to replace its content only if
nobody changed it since. Never search for something the taf made: use the id
it kept.

**Start read-only.** "Read the files in this folder and summarise them" is
a list or search tool plus `ONE_DRIVE_DOWNLOAD_FILE`. Offer a writing
tool only when the owner asked for changes, in their own folder or workbook.

## Don't offer

- Deleting tools (`ONE_DRIVE_DELETE_*`), above all
  `ONE_DRIVE_DELETE_ITEM_PERMANENTLY`, which skips the recycle bin.
- Sharing and permission tools: `ONE_DRIVE_CREATE_LINK`,
  `ONE_DRIVE_INVITE_USER_TO_ITEM`, `ONE_DRIVE_CREATE_ITEM_PERMISSION`,
  `ONE_DRIVE_GRANT_SHARE_PERMISSION`, `ONE_DRIVE_UPDATE_DRIVE_ITEMS_PERMISSIONS`,
  `ONE_DRIVE_REVOKE_SHARING_LINK_GRANTS`, and the tools that grant, change or
  remove access on team sites and their lists. They decide who else sees the
  owner's files.
- `ONE_DRIVE_RESTORE_ITEM_VERSION`, `ONE_DRIVE_MOVE_ITEM` and check-in or
  check-out tools, unless the owner asked: they change what others find.
- People and subscription tools (`ONE_DRIVE_GET_USER`,
  `ONE_DRIVE_LIST_SUBSCRIPTIONS`): not a taf's business.

## Access and plans

- **Workbook tools need a work or school account.** Microsoft's sheet and
  cell tools don't yet work on files in a personal OneDrive, and only on
  .xlsx files, never the older .xls. On a personal account, download the
  file instead.
- **Sheet and cell reads ask for write access** in the catalog, so a
  read-only sign-in lists them under `unusable`. Ask for write access when
  the job reads cells, and say the taf still writes nothing.
- **Uploads:** one call carries a file up to 250 MB; larger files go in
  parts through an upload session.
- **Workbook sessions end** after about 5 minutes without a call (about 7 if
  changes aren't kept). A job open longer starts a new one.
- **Storage is the owner's plan.** When the drive is full, uploads and new
  versions are refused until space is freed or bought.
- **Rate:** Microsoft slows down an app that calls too often, and may say how
  long to wait. A job reading thousands of files a run should page and
  pause.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, `accessDenied` (403) | The sign-in lacks this access, or the file isn't shared with the owner | "OneDrive won't let this step through: the connection doesn't cover it or the file isn't yours to change. Want to sign in again and allow it, or skip this step?" |
| Needs admin approval at sign-in | The organisation's admin must allow the app | "Your organisation's Microsoft admin has to approve this connection. Ask them, or connect a personal account." |
| The sign-in no longer works (401) | The owner removed access, or the sign-in expired | "OneDrive needs you to sign in again. Nothing was changed." |
| Not found, `itemNotFound` (404) | The file was moved, renamed away or deleted | "That file isn't where it was any more. The taf will skip it; tell it where the file is now?" |
| `nameAlreadyExists` (409) | A file of that name is already there | Pick a new name or ask; never overwrite silently. |
| `resourceModified` (412) | Someone changed the file since the taf read it | Read it again and retry once; else "The file changed while the taf worked. Nothing was overwritten." |
| `quotaLimitReached` (507) | The drive is full | "Your OneDrive is full. Free some space and the taf will carry on." |
| Too many calls (429), or unavailable (503) | Microsoft is slowing the app down | "OneDrive asked us to slow down. The taf will try again shortly." |
| A workbook call fails on a personal account, or on an .xls file | Workbook tools need a work account and .xlsx | "That spreadsheet can't be edited from here. Save it as .xlsx in a work OneDrive, or the taf can only download it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
