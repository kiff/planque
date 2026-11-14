# Google Meet

Checked: 2026-11-14 · Sources: the connection service's Google Meet tool
definitions as planque reads them (25 tools, [generated](generated/googlemeet.md)),
and Google's Meet API reference.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **One approval for the sign-in, one for the tools and the Card.**
- **To put a meeting in a calendar,** use Google Calendar: an event can carry
  a Meet link. This app is for meeting rooms, and for what happened in them.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLEMEET_CREATE_MEET` | Creates a meeting room with its link | write | **the room's name and link** |
| `GOOGLEMEET_GET_MEET` | Reads a room's settings and link | read | |
| `GOOGLEMEET_LIST_CONFERENCE_RECORDS` | Lists past meetings, by time or room | read | **each meeting's record id** |
| `GOOGLEMEET_LIST_PARTICIPANTS` | Lists who joined a past meeting | read | |
| `GOOGLEMEET_GET_TRANSCRIPTS_BY_CONFERENCE_RECORD_ID` | Lists a meeting's transcripts | read | transcript ids |
| `GOOGLEMEET_LIST_TRANSCRIPT_ENTRIES` | Reads a transcript, line by line | read | |
| `GOOGLEMEET_GET_RECORDINGS_BY_CONFERENCE_RECORD_ID` | Lists a meeting's recordings, as Drive files | read | the Drive file's id |
| `GOOGLEMEET_LIST_SMART_NOTES` | Lists a meeting's notes taken by Google | read | |

**What to keep:** a meeting's record id, to read its transcript once it's
ready. A transcript or recording appears some minutes after the meeting ends;
a taf that reads right away finds nothing, and tries again later rather than
saying there was none.

**A transcript is other people's words.** A job that summarises meetings
names who gets the summary, and the Card says so.

## Don't offer

- `GOOGLEMEET_END_ACTIVE_CONFERENCE`: it ends a meeting for everyone in it.
- Member and space changes (`GOOGLEMEET_CREATE_MEMBER`,
  `GOOGLEMEET_UPDATE_MEMBER`, `GOOGLEMEET_BATCH_UPDATE_MEMBERS`,
  `GOOGLEMEET_DELETE_MEMBER`, `GOOGLEMEET_UPDATE_SPACE`): they decide who can
  join and how.

## Access and plans

- **Transcripts, recordings and notes** exist only if someone turned them on
  in the meeting, and only on plans that include them: Google's business
  plans, and for recordings some paid personal storage plans. On a free
  personal account, a meeting has a list of who joined, and nothing more
  to read.
- **The list shows only meetings the owner organised.** A meeting they only
  joined can be read only by its record id.
- **Records are kept 30 days** by Google; recordings stay in Drive.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| No transcript or recording found | Not turned on, not ready yet, or the plan doesn't include it | "That meeting has no transcript yet. If it was turned on, it can take a few minutes; I'll look again." |
| Refused for missing access (403) | The sign-in didn't give this access, or the meeting isn't the owner's | "I can't read that meeting. It may belong to someone else, or Meet is only partly connected." |
| Not found (404) | Older than 30 days, or a wrong id | "That meeting's record isn't there any more." |
| Too many calls (429) | Over Google's limit | "Google asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, or expired | "Google Meet needs you to sign in again." |

After any refusal the taf stops that step and says so.
