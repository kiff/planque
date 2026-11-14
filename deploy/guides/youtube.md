# YouTube

Checked: 2026-11-14 · Sources: the connection service's YouTube tool
definitions as planque reads them (67 tools, [generated](generated/youtube.md)),
and the YouTube Data API documentation (quota, errors, comments) and the
YouTube Analytics API documentation.

## How it connects

- **A sign-in to the owner's account, choosing the YouTube channel**, held
  by the connection service. The
  catalog lists the tools and puts the ones its access doesn't cover under
  `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `YOUTUBE_LIST_CHANNEL_VIDEOS` | Lists the channel's videos | read | **video ids** |
| `YOUTUBE_GET_VIDEO_DETAILS_BATCH` | Reads videos' titles, views, likes | read | each video's details, by id |
| `YOUTUBE_GET_CHANNEL_STATISTICS` | Reads the channel's subscribers and views | read | the counts |
| `YOUTUBE_QUERY_ANALYTICS` | Reads the channel's analytics by day, video or country | read | rows of figures |
| `YOUTUBE_LIST_COMMENT_THREADS2` | Reads comments on a video or the channel | read | comment ids |
| `YOUTUBE_CREATE_COMMENT_REPLY` | Replies to a comment, as the channel | write, public | the reply's id |
| `YOUTUBE_LOAD_CAPTIONS` | Reads a video's captions | read | the caption text |
| `YOUTUBE_SEARCH_YOU_TUBE` | Searches public videos | read, **rationed** | video ids |

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

- **The API is free but rationed, per app project, each day:** searches
  and uploads each have their own allowance of about 100 calls a day, and
  everything else shares 10,000 units a day. A read usually costs 1 unit,
  and a write (a comment reply, say) about 50. Even a refused call costs at
  least one unit. The allowance resets at midnight Pacific time. A taf that
  searches every hour uses a big share of the searches: prefer listing the
  channel's own videos.
- **Analytics** are usually two to three days behind.
- **Comments** can be held for review before they're published.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, quota (403 `quotaExceeded`) | The day's allowance is used | "YouTube's daily allowance is used up. The taf will continue tomorrow." |
| Refused for missing access (403) | The sign-in didn't give this access, or the wrong channel was chosen | "YouTube is only partly connected, or connected to another channel. Sign in again and pick the right one." |
| Comments disabled (403) | The video has comments off | "Comments are turned off on that video." |
| Not found (404) | The video isn't there, or isn't visible to this channel | "That video isn't available any more." |
| The sign-in no longer works (401) | The sign-in is no longer valid | "YouTube needs you to sign in again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
