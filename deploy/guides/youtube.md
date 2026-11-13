# YouTube

Checked: 2026-11-13 · Sources: the connection service's YouTube tool
definitions as planque reads them (67 tools, [generated](generated/youtube.md)),
and Google's YouTube Data and Analytics API references.

## How it connects

- **A sign-in to Google, choosing the YouTube channel**, held by the
  connection service. An owner with a brand channel picks it at sign-in. The
  catalog lists the tools and puts the ones its access doesn't cover under
  `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `YOUTUBE_LIST_CHANNEL_VIDEOS` | Lists the channel's videos | read | **video ids** |
| `YOUTUBE_GET_VIDEO_DETAILS_BATCH` | Reads videos' titles, views, likes | read | |
| `YOUTUBE_GET_CHANNEL_STATISTICS` | Reads the channel's subscribers and views | read | |
| `YOUTUBE_QUERY_ANALYTICS` | Reads the channel's analytics by day, video or country | read | |
| `YOUTUBE_LIST_COMMENT_THREADS2` | Reads comments on a video or the channel | read | comment ids |
| `YOUTUBE_CREATE_COMMENT_REPLY` | Replies to a comment, as the channel | write, public | the reply's id |
| `YOUTUBE_LOAD_CAPTIONS` | Reads a video's captions | read | |
| `YOUTUBE_SEARCH_YOU_TUBE` | Searches public videos | read, **costly** | video ids |

**Replies are public and signed by the channel.** Each is shown to the owner
first, or bounded by the Card.

**What to keep:** video ids, and the last comment handled, so a reply job
never answers twice.

## Don't offer

- Deleting videos, playlists, comments or captions.
- Uploading videos (`YOUTUBE_UPLOAD_VIDEO`, `YOUTUBE_MULTIPART_UPLOAD_VIDEO`)
  and changing videos (`YOUTUBE_UPDATE_VIDEO`) unless the job is about
  publishing and the owner approves each one.
- Live broadcasts, subscribing, rating and reporting videos.
- Moderating others' comments
  (`YOUTUBE_SET_COMMENT_MODERATION_STATUS`) unless the job is moderation.

## Access and plans

- **The API is free but rationed:** each project has 10,000 units a day,
  shared by everyone on the connection. A read costs 1 unit, a search 100,
  an upload about 1,600. A taf that searches every hour uses a big share:
  prefer listing the channel's own videos.
- **Analytics** are a day or two behind.
- **Comments** may be held for review by the channel's own settings.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, quota (403 `quotaExceeded`) | The day's units are used | "YouTube's daily allowance is used up. The taf will continue tomorrow." |
| Refused for missing access (403) | The sign-in didn't give this access, or the wrong channel was chosen | "YouTube is only partly connected, or connected to another channel. Sign in again and pick the right one." |
| Comments disabled (403) | The video has comments off | "Comments are turned off on that video." |
| Not found (404) | A deleted or private video | "That video isn't available any more." |
| The sign-in no longer works (401) | Removed, or expired | "YouTube needs you to sign in again." |

After any refusal the taf stops that step and says so.
