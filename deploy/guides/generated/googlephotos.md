# Google Photos: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

13 tools: 6 read, 7 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLEPHOTOS_ADD_ENRICHMENT` | Add Enrichment | write | `albumId`, `albumPosition`, `newEnrichmentItem` |  | `photoslibrary.appendonly`, `photoslibrary.edit.appcreateddata` |
| `GOOGLEPHOTOS_BATCH_ADD_MEDIA_ITEMS` | Batch Add Media Items | write | `albumId`, `mediaItemIds` |  | `photoslibrary.appendonly`, `photoslibrary.edit.appcreateddata` |
| `GOOGLEPHOTOS_BATCH_CREATE_MEDIA_ITEMS` | Batch Create Media Items | write | *albumId, albumPosition, files, media_files, urls* |  | `photoslibrary` |
| `GOOGLEPHOTOS_BATCH_GET_MEDIA_ITEMS` | Batch Get Media Items | read | `mediaItemIds` |  | `photoslibrary.readonly.appcreateddata` |
| `GOOGLEPHOTOS_CREATE_ALBUM` | Create Album | write | `title` |  | `photoslibrary.appendonly` |
| `GOOGLEPHOTOS_GET_ALBUM` | Get Album | read | `albumId` |  | `photoslibrary.readonly` |
| `GOOGLEPHOTOS_GET_MEDIA_ITEM_DOWNLOAD` | Download Photos Media Item | read | `mediaItemId` |  | `photoslibrary.readonly` |
| `GOOGLEPHOTOS_LIST_ALBUMS` | List Albums | read | *excludeNonAppCreatedData, pageSize, pageToken* |  | `photoslibrary.readonly` |
| `GOOGLEPHOTOS_LIST_MEDIA_ITEMS` | List Media Items (App-Created Only) | read | *pageSize, pageToken* |  | `photoslibrary.readonly.appcreateddata` |
| `GOOGLEPHOTOS_SEARCH_MEDIA_ITEMS` | Search Media Items | read | *albumId, filters, orderBy, pageSize, pageToken* |  | `photoslibrary`, `photoslibrary.readonly` |
| `GOOGLEPHOTOS_UPDATE_ALBUM` | Update Album | write | `albumId`; *coverPhotoMediaItemId, title* |  | `photoslibrary.edit.appcreateddata` |
| `GOOGLEPHOTOS_UPDATE_MEDIA_ITEM` | Update Media Item | write | `description`, `mediaItemId` |  | `photoslibrary.readonly` |
| `GOOGLEPHOTOS_UPLOAD_MEDIA` | Upload Media | write | *description, file_name, file_to_upload, url* |  | `photoslibrary.appendonly`, `photoslibrary` |
