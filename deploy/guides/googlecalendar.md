# Google Calendar

Checked: 2026-11-13 · Sources: the connection service's Google Calendar tool
definitions as planque reads them (47 tools, [generated](generated/googlecalendar.md)),
and Google's Calendar API reference.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Reading free time needs less access than writing events, so ask for every
  tool the job needs in the first request.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLECALENDAR_LIST_CALENDARS` | Lists the owner's calendars | read | each calendar's id |
| `GOOGLECALENDAR_EVENTS_LIST` | Lists events in a time window, a page at a time | read | event ids |
| `GOOGLECALENDAR_FIND_EVENT` | Finds events by words and dates | read | event ids |
| `GOOGLECALENDAR_EVENTS_GET` | Reads one event | read | its attendees and times |
| `GOOGLECALENDAR_FIND_FREE_SLOTS` | Finds free time across calendars | read | the free windows |
| `GOOGLECALENDAR_FREE_BUSY_QUERY` | Says when calendars are busy | read | busy windows |
| `GOOGLECALENDAR_GET_CURRENT_DATE_TIME` | Gives today's date and time in a time zone | read | |
| `GOOGLECALENDAR_CREATE_EVENT` | Creates an event, with or without guests | write | **the event's id** |
| `GOOGLECALENDAR_PATCH_EVENT` | Changes some fields of one event | write | the event's id |
| `GOOGLECALENDAR_QUICK_ADD` | Creates an event from one line of text | write | the event's id |
| `GOOGLECALENDAR_DELETE_EVENT` | Deletes one event | write | |

**What to keep:** the id of every event the taf creates, so it moves or
cancels that event later rather than finding it by title. The calendar's id
too, when it isn't the owner's main one (`primary`).

**Guests get an email.** Creating, moving or cancelling an event with guests
notifies them by default. The Card says so, and the owner sees each one
before it's sent, unless the owner chose otherwise.

**Times carry a time zone.** Ask for times with their zone, or read the
calendar's own; never assume the server's.

## Don't offer

- Sharing settings (`GOOGLECALENDAR_ACL_INSERT`, `GOOGLECALENDAR_ACL_UPDATE`,
  `GOOGLECALENDAR_ACL_DELETE`): they decide who sees the owner's calendar.
- Deleting or clearing a whole calendar (`GOOGLECALENDAR_CALENDARS_DELETE`).
- Watching for changes (`GOOGLECALENDAR_EVENTS_WATCH`,
  `GOOGLECALENDAR_CALENDAR_LIST_WATCH`): a taf reads on its schedule instead.
- `GOOGLECALENDAR_DELETE_EVENT` for events the taf didn't create, unless the
  job is clearly about cancelling them.

## Access and plans

- **No paid plan is needed.** Rooms and buildings
  (`GOOGLECALENDAR_LIST_BUILDINGS`, `GOOGLECALENDAR_LIST_CALENDAR_RESOURCES`)
  exist only on business accounts.
- **Free time only** needs the narrowest access: a job that only finds slots
  should ask for nothing more.
- **Rate:** Google allows a generous number of calls per user, but sending
  many invitations in a short time is limited to stop spam. A taf creating
  more than a few dozen events with guests a day should pause between them.
- **Recurring events:** changing one occurrence uses its instance
  (`GOOGLECALENDAR_EVENTS_INSTANCES`); changing the series changes every one.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in gave read-only or free-time access, and the tool writes | "Google Calendar is only partly connected: adding events needs access you didn't give. Sign in again and allow it, or I'll only read." |
| Refused, rate or invitation limit (403 or 429) | Too many calls, or too many invitations sent | "Google asked us to slow down sending invitations. The rest go in a few minutes." |
| The event isn't found (404, or 410 "deleted") | Someone deleted it | "That event isn't on the calendar any more. Nothing was changed." |
| A conflict (409) | An event with the same id exists | Use a new event; don't tell the owner unless it fails again. |
| The sign-in no longer works (401) | Removed, or expired | "Google Calendar needs you to sign in again. Nothing was changed." |
| A time is refused (400) | A time without a zone, or an end before its start | Fix the time and try once; else "Google didn't accept that time. Nothing was booked." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create.
