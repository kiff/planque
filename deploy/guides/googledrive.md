# Google Drive

Checked: 2026-11-14 · Sources: the connection service's Google Drive tool
definitions as planque reads them (92 tools, [generated](generated/googledrive.md)),
Google's Drive API reference, and a test run on a phone on 13 Nov 2026.
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **One approval for the sign-in, one for the tools and the Card.**

## Two kinds of access

What Google's page asks for follows the taf's tools, each at the narrowest
access it runs on (Google's words, roughly, in quotes):

- **Reads** (search, list, download, export): **read-only Drive**, "see and
  download all your Drive files".
- **Writes** (upload, create, move, share, delete): **full Drive**, "see,
  edit, create and delete all your Drive files". Access to only the files the
  taf made is narrower, but planque doesn't ask for it: a tool that lists it
  may still need more, and the taf then fails on the owner's own files.
- The tools are named in the connect: a connect that names none asks for the
  app's widest access, so the taf's planned tools are always named when it
  connects.

A sign-in made on files-the-taf-made access (an older one) finds only the
taf's own files: a Sheets sign-in on 13 Nov 2026 couldn't search at all for
this reason. Keep the id of every file the taf makes. A link from the owner doesn't open a file on narrow
access; for files the owner already has, the sign-in needs the whole Drive.

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
  `GOOGLEDRIVE_PATCH_PERMISSION`, `GOOGLEDRIVE_DELETE_PERMISSION`,
  `GOOGLEDRIVE_CREATE_PERMISSIONS_BATCH`), and publishing a revision
  (`GOOGLEDRIVE_UPDATE_FILE_REVISION_METADATA`): it decides who outside can see
  the owner's files. Offer it only for a job that is about sharing, with each
  share shown to the owner first.
- Deleting for good (`GOOGLEDRIVE_GOOGLE_DRIVE_DELETE_FOLDER_OR_FILE_ACTION`,
  `GOOGLEDRIVE_DELETE_REVISION`, `GOOGLEDRIVE_EMPTY_TRASH`): it skips the
  bin. Use the bin instead.
- Shared drives (create, delete, hide), and the change feeds and watches:
  they act on the whole Drive, beyond one job's files.

## Access and plans

- **Storage is the owner's:** 15 GB free across Gmail, Drive and Photos. An
  upload fails when it's full.
- **A business account's admin** can block outside apps, or block sharing
  outside the business.
- **Size:** `GOOGLEDRIVE_UPLOAD_FILE` is for files up to 5 MB; larger ones
  need a resumable upload. An export is limited to 10 MB.
- **Rate:** Google allows many reads a minute per user, and fewer writes. A
  job moving hundreds of files should pause between batches.

## When it fails

The rows are Google Drive's own answers: a refusal reaches the taf as
"Google Drive refused it (it answered …): …", with Google Drive's reason. planque's
own refusals, and an expired sign-in, are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | The sign-in only covers the taf's own files, and this one isn't | "I can only see the files I made in your Drive. To use your other files, sign in again and allow access to your whole Drive." |
| Refused (404), file not found | It's deleted, or not visible on this access | "I can't open that file. It may have been deleted, or it isn't shared with this account." |
| Refused (403), drive full (storage quota) | The owner's storage is used up | "Your Google storage is full, so the file wasn't saved. Free some space or add storage." |
| Refused (403 or 429), rate limit | Over Google's limit | "Google Drive asked us to slow down. The taf will continue in a minute." |
| Google Drive's sign-in has expired (401) | Removed, or expired | "Google Drive needs you to sign in again. Nothing was changed." |
| Refused (403), sharing | The business doesn't allow sharing outside | "Your Google admin doesn't allow sharing outside the business." |

A failed search never leads to creating a second copy.
