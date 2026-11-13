# Cal

Checked: 2026-11-13 · Sources: the Cal tool definitions as planque reads
them (170 tools, generated/cal.md), and Cal.com's API v2 reference
(introduction, bookings, slots) and its pricing page.

## How it connects

- **A sign-in to Cal** on Cal's own page, held by the connection service;
  the taf never sees a password or token. The catalog lists the tools and,
  with a known sign-in, puts the ones its access doesn't cover under
  `unusable`, with why. Few Cal tools say what access they need, so most
  failures show up only when the tool runs.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request, so nothing is asked
  twice.

## Tools a job usually needs

Read means it changes nothing in Cal. The generated part gives no returns;
the ids below are the ones Cal's API reference shows.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `CAL_RETRIEVE_MY_INFORMATION` | Reads the signed-in owner's profile | read | the user id, username, time zone |
| `CAL_LIST_EVENT_TYPES` | Lists the owner's kinds of meeting (a 30-minute call, say) | read | event type ids, slugs |
| `CAL_RETRIEVE_EVENT_TYPE_BY_ID` | Reads one event type's length, questions and rules | read | the event type id |
| `CAL_GET_AVAILABLE_SLOTS_INFO` | Lists free times for an event type between two dates | read | start times, by date |
| `CAL_FETCH_ALL_BOOKINGS` | Lists bookings, filtered by status or dates | read | booking ids and uids |
| `CAL_RETRIEVE_BOOKING_DETAILS_BY_UID` | Reads one booking with its attendees | read | the booking uid, start, end, status |
| `CAL_RETRIEVE_SCHEDULES_LIST` | Reads the owner's working hours | read | schedule ids |
| `CAL_POST_NEW_BOOKING_REQUEST` | Books a time for an attendee | write, with attendees | the booking id and uid |
| `CAL_RESCHEDULE_BOOKING_BY_UID` | Moves a booking to another time | write, with attendees | the booking uid |
| `CAL_CANCEL_BOOKING_VIA_UID` | Cancels a booking, with an optional reason | write, with attendees | the cancelled booking |
| `CAL_CONFIRM_BOOKING_BY_UID` / `CAL_DECLINE_BOOKING_WITH_REASON` | Accepts or turns down a booking that waits for the owner's approval | write, with attendees | the booking uid, its status |

**What to keep:** a booking's uid (the cancel, reschedule and confirm tools
take the uid, not the number id), and the event type id a job books into.
Never search for a booking the taf made: use the uid it kept.

**Read the free times first.** Book only a time `CAL_GET_AVAILABLE_SLOTS_INFO`
just returned, in the attendee's time zone, so the booking matches what Cal
offers.

**Start read-only.** "Tell me what's booked tomorrow" needs only read
tools. Booking, moving, cancelling or declining touches someone else's day:
offer those tools only when the owner asked, and say each one waits for
their yes. A cancelled recurring booking can take every later one with it.

## Don't offer

- Deleting: `CAL_DELETE_EVENT_TYPE_BY_ID`, `CAL_DELETE_SCHEDULE_BY_ID`,
  `CAL_DELETE_TEAM_BY_ID`, `CAL_DELETE_WEBHOOK_BY_ID` and the other
  `CAL_DELETE_*` tools.
- People and teams: `CAL_ADD_MEMBER_TO_TEAM_USING_ORG_AND_TEAM_ID`,
  `CAL_CREATE_TEAM_MEMBERSHIP_WITH_ROLE`, `CAL_CREATE_TEAM_INVITE_LINK`,
  `CAL_POST_USER_TO_ORGANIZATION`, `CAL_DELETE_USER_FROM_ORGANIZATION`,
  `CAL_MODIFY_ORGANIZATION_MEMBERSHIP_BY_ID`, and the booking reassign tools
  (`CAL_REASSIGN_BOOKING_TO_ANOTHER_USER`, `CAL_REASSIGN_BOOKING_WITH_UID`).
- Connections and credentials: `CAL_FETCH_PROVIDER_ACCESS_TOKEN`,
  `CAL_POST_CALENDAR_CREDENTIALS`, `CAL_CONNECT_TO_CALENDAR`,
  `CAL_DISCONNECT_CALENDAR_USING_CREDENTIAL_ID`,
  `CAL_POST_CONFERENCING_APP_CONNECT`, the tools that link a payment
  account, and every OAuth-client tool.
- Webhooks: `CAL_POST_WEBHOOK_EVENT`, `CAL_CREATE_WEBHOOK_FOR_EVENT_TYPE`,
  `CAL_CREATE_ORGANIZATION_WEBHOOK_BY_ORG_ID` and their update tools. They
  send booking details to another address.
- Settings: `CAL_UPDATE_USER_PROFILE_DETAILS`, `CAL_EDIT_EVENT_TYPE_BY_ID`,
  `CAL_UPDATE_SCHEDULE_BY_ID`, `CAL_SET_DEFAULT_CONFERENCING_APP`, and the
  organisation attribute tools. They change what every visitor to the
  owner's booking page sees.

## Access and plans

- **Team and organisation tools** (team bookings, team event types,
  organisation members): Cal's API reference ties team tools to its team
  plan and organisation tools to its organisation plan. Offer them only to
  an owner on that plan; otherwise use the personal tools.
- **Rate:** Cal limits calls per minute, and each way of connecting has its
  own limit. It doesn't publish the number for a sign-in, so a job that
  reads many bookings pages and pauses.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works | The owner removed access, or the sign-in expired | "Cal needs you to sign in again. Nothing was booked or changed." |
| Refused for missing access | The sign-in doesn't cover this tool, or the owner's plan doesn't include it | "Cal didn't let us do this with your account. It may need a team plan, or a new sign-in. Want to do without this step?" |
| A booking isn't found | It was cancelled, or the uid is wrong | "That booking isn't in Cal any more. The taf will skip it." |
| A booking is refused | The time may no longer be offered, or a required answer is missing | Read the free times again and try once; else "Cal didn't accept that booking. Nothing was booked." |
| Too many calls | Over Cal's per-minute limit | "Cal asked us to slow down. The taf will try again in a minute." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
