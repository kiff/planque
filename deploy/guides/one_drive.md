# OneDrive

Checked: 2026-11-13 · Sources: the connection service's OneDrive tool
definitions as planque reads them (101 tools), and OneDrive's API
documentation (items, uploads, error answers, working at scale).

## How it connects

- **A sign-in to OneDrive**, held by the connection service; the taf never
  sees a password or token. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`,
  with why.
- **The sign-in acts as the person.** The taf sees the files that person
  can open, and nothing more.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request, so nothing is asked
  twice.

## Tools a job usually needs

Read means it changes nothing in the drive. Returns are what a taf keeps to
use next time. The catalog's `returns` is empty for OneDrive today; the ids
below are those OneDrive's answers carry. A file or folder's answer holds
its id, its name, a link that opens it in a browser, and its parent folder.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ONE_DRIVE_GET_ROOT` | Reads the top folder of the person's drive | read | the folder's id |
| `ONE_DRIVE_LIST_FOLDER_CHILDREN` | Lists what is in one folder | read | item ids, names, links |
| `ONE_DRIVE_SEARCH_ITEMS` | Finds files and folders by words | read | matching item ids, names, links |
| `ONE_DRIVE_GET_ITEM` | Reads one file or folder's details | read | the item's id, link, parent folder |
| `ONE_DRIVE_GET_ITEM_BY_PATH` | Reads an item by its folder path | read | the item's id, link |
| `ONE_DRIVE_DOWNLOAD_FILE` | Fetches a file's content | read | the file's content |
| `ONE_DRIVE_LIST_ROOT_DRIVE_CHANGES` | Lists what changed in the drive since the last check | read | changed items, a link for the next check |
| `ONE_DRIVE_ONEDRIVE_CREATE_FOLDER` | Makes a new folder | write | the folder's id, link |
| `ONE_DRIVE_ONEDRIVE_UPLOAD_FILE` | Adds a new file to a folder | write | the file's id, link |
| `ONE_DRIVE_UPDATE_FILE_CONTENT` | Replaces a file's content | write | the file's id |
| `ONE_DRIVE_UPDATE_DRIVE_ITEM_METADATA` | Renames a file or changes its details | write | the item |
| `ONE_DRIVE_MOVE_ITEM` | Moves a file or folder to another folder | write | the item, with its new parent |
| `ONE_DRIVE_CREATE_LINK` | Makes a sharing link to a file | write, outside | the sharing link |

**What to keep:** the id of each folder and file the taf made or works on,
to open the same one next time; the change-check link, so the next run asks
only for what is new. Never search for something the taf made: use the id it
kept.

**Start read-only.** "File the receipts into a folder" is a list or search,
a download, and an upload or a move. A sharing link lets other people open
the file, so offer it only when the owner asked for sharing, and each link
waits for the owner's yes.

## Don't offer

- Deleting (`ONE_DRIVE_DELETE_ITEM`, `ONE_DRIVE_DELETE_ITEM_PERMANENTLY`):
  the first sends to the recycle bin, the second can't be undone.
- Sharing and permission tools (`ONE_DRIVE_INVITE_USER_TO_ITEM`,
  `ONE_DRIVE_CREATE_ITEM_PERMISSION`, `ONE_DRIVE_GRANT_SHARE_PERMISSION`,
  `ONE_DRIVE_UPDATE_DRIVE_ITEMS_PERMISSIONS` and the others that add or
  remove access): they change who can open the owner's files.
- The site and list tools: they reach a company's shared sites, not the
  person's own drive.
- Spreadsheet cell tools (the `ONE_DRIVE_*_WORKBOOK_*` ones): a job that
  works inside a spreadsheet needs its own guide and checks.
- Check-out, check-in and change subscriptions: the taf reads on its
  schedule, and a forgotten check-out locks a file for others.

## Access and plans

- **Reading needs file read access; uploading, moving and renaming need
  file write access.** A sign-in made for reading doesn't let the taf write.
- **One upload takes files up to 250 MB.** Bigger files need an upload in
  parts (`ONE_DRIVE_CREATE_UPLOAD_SESSION`), which a usual job doesn't need.
- **A download link** from OneDrive works for about an hour. Don't keep it;
  keep the file's id.
- **A file changed by someone else** since the taf read it can refuse an
  update. Read it again before writing.
- **Rate:** OneDrive slows apps that call too often, more at busy hours, and
  says how long to wait. An app that doesn't wait can be blocked. A job
  that checks a folder often should use the change list rather than reading
  every file each run.
- **Storage:** a full drive refuses uploads until space is freed.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, `accessDenied`) | The sign-in didn't give this tool's access, or the person can't open that file | "OneDrive is only partly connected, or that file isn't yours to change. Want to sign in again and allow it, or do without this step?" |
| The sign-in no longer works (401) | The sign-in expired or was removed | "OneDrive needs you to sign in again. Nothing was changed." |
| Not found (404, `itemNotFound`) | The file or folder was deleted or moved | "That file isn't in OneDrive any more. The taf will skip it." |
| Name taken, or folder missing (409) | A file with that name exists, or the target folder isn't there | "There's already a file with that name, or the folder is gone. The taf will stop and show you what it tried." |
| Changed since read (412, `resourceModified`) | Someone edited the file after the taf read it | "Someone changed that file meanwhile. The taf will read it again before writing." |
| File too large (413) | Over the size one call accepts | "That file is too big to add in one go. The taf has stopped." |
| Too many calls (429 or 503) | OneDrive is slowing this app down | "OneDrive asked us to slow down. The taf will try again later." |
| Drive full (507, `quotaLimitReached`) | The person's storage is used up | "Your OneDrive is full. Free some space, and the taf will try again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find its folder asks, rather than making
another one.
