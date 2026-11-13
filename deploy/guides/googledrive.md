# Google Drive

Checked: 2026-11-13 · Sources: the connection service's Google Drive tool
definitions as planque reads them (92 tools, [generated](generated/googledrive.md)),
Google's Drive API reference, and a test run on a phone on 13 Nov 2026.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **One approval for the sign-in, one for the tools and the Card.**

## Two kinds of access

Drive access comes in two widths, and a job should ask for the narrowest:

- **Files the taf made** (or the owner opened with it): enough for a taf that
  files its own receipts or reports. It can't see anything else in Drive,
  including files the owner made by hand.
- **The whole Drive**, read-only or full: needed to find or read files the
  owner already has.

A search on narrow access finds only the taf's own files. A Sheets sign-in
on 13 Nov 2026 couldn't search at all for this reason. Keep the id of every
file the taf makes, and ask the owner for a link to anything else.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLEDRIVE_FIND_FILE` | Finds files by name, type, folder or date | read | file ids |
| `GOOGLEDRIVE_FIND_FOLDER` | Finds folders by name | read | folder ids |
| `GOOGLEDRIVE_GET_FILE_METADATA` | Reads a file's name, type, folder, link | read | its link |
| `GOOGLEDRIVE_DOWNLOAD_FILE` | Downloads a file's content | read | |
| `GOOGLEDRIVE_EXPORT_GOOGLE_WORKSPACE_FILE` | Exports a Google Doc, Sheet or Slides file as PDF or text | read | |
| `GOOGLEDRIVE_CREATE_FOLDER` | Creates a folder | write | **the folder's id** |
| `GOOGLEDRIVE_UPLOAD_FILE` | Uploads a file into a folder | write | **the file's id** |
| `GOOGLEDRIVE_CREATE_FILE_FROM_TEXT` | Creates a file from text | write | the file's id |
| `GOOGLEDRIVE_MOVE_FILE` | Moves a file to another folder | write | |
| `GOOGLEDRIVE_TRASH_FILE` | Puts a file in the bin, where it stays 30 days | write | |

**What to keep:** the id of the folder the taf files into, and of each file
it uploads. Finding "Receipts" by name every run breaks the day the owner
makes a second folder with that name.

## Don't offer

- Sharing (`GOOGLEDRIVE_CREATE_PERMISSION`, `GOOGLEDRIVE_UPDATE_PERMISSION`,
  `GOOGLEDRIVE_CREATE_PERMISSIONS_BATCH`): it decides who outside can see
  the owner's files. Offer it only for a job that is about sharing, with each
  share shown to the owner first.
- Deleting for good (`GOOGLEDRIVE_GOOGLE_DRIVE_DELETE_FOLDER_OR_FILE_ACTION`,
  `GOOGLEDRIVE_DELETE_REVISION`): it skips the bin. Use the bin instead.
- Shared drives (create, delete, hide), and the change feeds and watches:
  they act on the whole Drive, beyond one job's files.

## Access and plans

- **Storage is the owner's:** 15 GB free across Gmail, Drive and Photos. An
  upload fails when it's full.
- **A business account's admin** can block outside apps, or block sharing
  outside the business.
- **Size:** a simple upload is fine up to a few MB; larger files go by
  resumable upload, which the taf does in parts.
- **Rate:** Google allows many reads a minute per user, and fewer writes. A
  job moving hundreds of files should pause between batches.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | The sign-in only covers the taf's own files, and this one isn't | "I can only see the files I made in your Drive. Send me the file's link, or sign in again and allow access to your Drive." |
| The file isn't found (404) | It's deleted, or not visible on this access | "I can't open that file. It may have been deleted, or it isn't shared with this account." |
| Drive is full (403, storage quota) | The owner's storage is used up | "Your Google storage is full, so the file wasn't saved. Free some space or add storage." |
| Too many calls (403 rate limit, or 429) | Over Google's limit | "Google Drive asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, or expired | "Google Drive needs you to sign in again. Nothing was changed." |
| Sharing is refused (403) | The business doesn't allow sharing outside | "Your Google admin doesn't allow sharing outside the business." |

After any refusal the taf stops that step and says so. A failed search never
leads to creating a second copy.
