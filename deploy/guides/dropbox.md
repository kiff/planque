# Dropbox

Checked: 2026-11-13 · Sources: the connection service's Dropbox tool
definitions as planque reads them (175 tools), Dropbox's API reference
(the published API spec), its performance, error-handling, OAuth and data
transport limit pages.

## How it connects

- **A sign-in to Dropbox**, held by the connection service; the taf never
  sees a password or token. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`,
  with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  If it doesn't (writing files after a read-only sign-in, say), Dropbox asks
  the owner again.

## Tools a job usually needs

Read means it changes nothing in the Dropbox. Returns are named from
Dropbox's API reference, since the catalog gives none for this app.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `DROPBOX_LIST_FILES_IN_FOLDER` | Lists a folder's files and subfolders, optionally all the way down | read | entries (each with id, path, rev for files), a cursor, has-more |
| `DROPBOX_LIST_FOLDER_CONTINUE` | Reads the next page of that list | read | entries, a new cursor |
| `DROPBOX_SEARCH_FILE_OR_FOLDER` | Finds files and folders by name or content | read | matches with their metadata, a cursor |
| `DROPBOX_SEARCH_CONTINUE` | Reads the next page of a search | read | matches, a cursor |
| `DROPBOX_GET_METADATA` | Reads one item's details | read | id, path, rev, size, dates |
| `DROPBOX_READ_FILE` | Downloads one file | read | the file's content and metadata |
| `DROPBOX_GET_TEMPORARY_LINK` | Makes a private download link that stops working after four hours | read | the link, the file's metadata |
| `DROPBOX_LIST_FILE_REVISIONS` | Lists a file's earlier versions | read | revs, dates |
| `DROPBOX_CREATE_FOLDER` | Creates a folder | write | the folder's id and path |
| `DROPBOX_UPLOAD_FILE` | Saves one file, up to 150 MB | write | id, path, rev |
| `DROPBOX_COPY_FILE_OR_FOLDER` | Copies an item to another path | write | the copy's metadata |
| `DROPBOX_MOVE_FILE_OR_FOLDER` | Moves or renames an item | write | the item's new metadata |
| `DROPBOX_CREATE_SHARED_LINK` | Makes a link others can open | write, outside | the link's url and id |
| `DROPBOX_CREATE_FILE_REQUEST` | Opens a page where others can upload into a folder | write, outside | the request's id and url |

**What to keep:** a file's id (it survives renames and moves) and its rev,
to update that same version later; a list's cursor, to read only what
changed next time. Never search for what the taf made: use the id it kept.

**Uploads don't overwrite by default.** Keep the default "add" mode, or
update a given rev; use overwrite only when the owner asked to replace a
file. A name clash then saves a renamed copy or is refused, never lost.

**Start read-only:** list, search and read. Offer uploads when the job
saves something, and a shared link or file request only when the owner
asked to share; each link is shown to the owner first.

## Don't offer

- Deleting tools (`DROPBOX_DELETE_FILE`, `DROPBOX_DELETE_BATCH`,
  `DROPBOX_DELETE_FILE_REQUESTS`, `DROPBOX_DELETE_ALL_CLOSED_FILE_REQUESTS`,
  `DROPBOX_DELETE_TEAM_FOLDER_PERMANENTLY`).
- People and permissions: members of files and folders, sharing or
  unsharing folders, folder policies, changing or revoking shared links
  (`DROPBOX_SHARE_FOLDER`, `DROPBOX_ADD_FOLDER_MEMBER_ACTION` and the like).
- Team admin: every team member, group, team folder, quota, allowlist,
  device, linked-app and audit-log tool, and admin roles.
- Bulk and outside-content tools: copy, move, folder and upload batches,
  and `DROPBOX_SAVE_URL`, which pulls any web address into the Dropbox.
- Account settings such as `DROPBOX_SET_PROFILE_PHOTO`.

## Access and plans

- **Reading** needs Dropbox's file read access; **saving, copying and
  moving** need file write access; **links and file requests** need
  sharing or file-request access. A sign-in made for reading can't save.
- **File size:** one upload takes up to 150 MB. Bigger files go through an
  upload session (`DROPBOX_START_UPLOAD_SESSION`,
  `DROPBOX_APPEND_UPLOAD_SESSION`, `DROPBOX_FINISH_UPLOAD_SESSION`), up to
  350 GB; offer them only for large files.
- **Pages:** a search answers at most 1,000 matches a page; lists answer a
  page and a cursor for the rest.
- **Rate:** Dropbox limits calls per user without publishing numbers.
  Parallel writes are refused more often: write files one at a time.
- **Some business teams** have a monthly allowance of upload-type calls,
  shared by every app on the team and reset on the first of the month.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | The owner removed access, or a team admin suspended the account | "Dropbox needs you to sign in again. Nothing was changed." |
| Refused for missing access (401, missing scope) | The sign-in didn't give this tool's access | "Dropbox is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Not found (409, path not found) | The file was moved, renamed, deleted or unshared | "That file isn't where it was any more. The taf will skip it." |
| Conflict on upload (409, conflict) | A file with that name exists | Save under a new name once; else "A file with that name is already there. Nothing was replaced." |
| No space (409, insufficient space) | The Dropbox is full | "Your Dropbox is full, so the file wasn't saved. Free some space or upgrade, and the taf will try again." |
| Refused (403), monthly limit | The team used its monthly upload allowance | "Your Dropbox team has used this month's upload allowance. Uploads resume on the 1st; a team admin can check usage in the admin console." |
| Too many calls (429) | Too many calls or parallel writes | "Dropbox asked us to slow down. The taf will try again shortly." |
| Bad input (400) | The tool doesn't fit this account or the input is wrong | "Dropbox didn't accept that step. Nothing was changed." |
| Server error (500) | A brief fault on Dropbox's side | Try again later; if it repeats, "Dropbox is having trouble right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
