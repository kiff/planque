# Google Photos

Checked: 2026-11-14 · Sources: the connection service's Google Photos tool
definitions as planque reads them (13 tools, [generated](generated/googlephotos.md)),
and Google's Photos API reference.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## What a taf can see: only what it made

Google no longer lets outside apps browse a person's whole photo library.
A taf sees **the photos and albums it uploaded or created itself**, and
nothing else. "Find my photos from last summer" can't be done; "put each
receipt photo I send into an album" can.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLEPHOTOS_CREATE_ALBUM` | Creates an album | write | **the album's id** |
| `GOOGLEPHOTOS_UPLOAD_MEDIA` | Uploads a photo or video | write | **the item's id** |
| `GOOGLEPHOTOS_BATCH_CREATE_MEDIA_ITEMS` | Uploads several, into an album | write | the items' ids |
| `GOOGLEPHOTOS_BATCH_ADD_MEDIA_ITEMS` | Adds the taf's items to its album | write | |
| `GOOGLEPHOTOS_LIST_ALBUMS` | Lists the albums the taf made | read | album ids |
| `GOOGLEPHOTOS_LIST_MEDIA_ITEMS` | Lists the items the taf uploaded | read | item ids |
| `GOOGLEPHOTOS_SEARCH_MEDIA_ITEMS` | Lists the taf's items in an album or by date | read | item ids |
| `GOOGLEPHOTOS_GET_MEDIA_ITEM_DOWNLOAD` | Downloads one of the taf's items | read | |

**What to keep:** the album's id, from the first run. The taf can't find an
album the owner made by hand.

## Don't offer

- Any job that reads the owner's existing photos: it can't work, and the
  guide's refusal would only come after the owner signed in.
- Sharing albums: the tools don't offer it, and the owner shares in the app.

## Access and plans

- **Storage is the owner's:** 15 GB free across Gmail, Drive and Photos.
  Uploads count against it.
- **Size:** photos up to 200 MB and videos up to 20 GB each.
- **Rate:** about 10,000 calls a day, and 75,000 for downloading, shared by
  everyone using the same connection. A single taf rarely comes near it.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Nothing found, or refused, on a photo the owner took | Outside apps see only what they uploaded | "I can only see photos I added myself. Send me the photo and I'll add it." |
| Storage full (403) | The owner's storage is used up | "Your Google storage is full, so the photo wasn't saved." |
| Not found (404) | The owner deleted it in the app | "That photo isn't there any more." |
| Too many calls (429) | Over the day's limit | "Google asked us to slow down. The taf will continue later." |
| The sign-in no longer works (401) | Removed, or expired | "Google Photos needs you to sign in again." |

After any refusal the taf stops that step and says so.
