# Calendly

Checked: 2026-11-14 · Sources: the connection service's Calendly tool
definitions as planque reads them (53 tools), and Calendly's developer
documentation (rate limits, frequently asked questions).

## How it connects

- **A sign-in to Calendly**, held by the connection service; the taf never
  sees a password or token. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`,
  with why.
- **The sign-in acts as one Calendly user.** It sees that person's event
  types and bookings.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.
- **Calendly ends the sign-in** when the person changes their login email,
  password or way of logging in. Then the owner signs in again.

## Tools a job usually needs

Read means it changes nothing in Calendly. Calendly's catalog entries give
no `returns` paths today; the ids below are those Calendly's answers carry
(each is a link to the object, ending in its id).

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `CALENDLY_GET_USER` | Reads the signed-in person, or another user by id | read | the user's link, the organization's link |
| `CALENDLY_LIST_EVENT_TYPES` | Lists the kinds of meeting people can book | read | event type links, booking page urls |
| `CALENDLY_GET_EVENT_TYPE` | Reads one kind of meeting | read | the event type's link |
| `CALENDLY_LIST_SCHEDULED_EVENTS` | Lists booked meetings in a date range, optionally by guest email | read | event links, a page token |
| `CALENDLY_GET_EVENT` | Reads one booked meeting: time, place, status | read | the event's link |
| `CALENDLY_LIST_EVENT_INVITEES` | Lists who booked a meeting and their answers | read | invitee links, reschedule and cancel urls |
| `CALENDLY_GET_EVENT_INVITEE` | Reads one guest of a meeting | read | the invitee's link |
| `CALENDLY_LIST_EVENT_TYPE_AVAILABLE_TIMES` | Lists open slots for one kind of meeting | read | start times, booking urls |
| `CALENDLY_LIST_USER_BUSY_TIMES` | Lists when the person is busy | read | busy periods |
| `CALENDLY_CREATE_SINGLE_USE_SCHEDULING_LINK` | Makes a booking link that works for one booking | write | the link's url |
| `CALENDLY_POST_INVITEE` | Books a slot for a guest | write, outside | the invitee's link, the event's link |
| `CALENDLY_CANCEL_SCHEDULED_EVENT` | Cancels a booked meeting, with a reason | write, outside | the cancellation |

**What to keep:** the signed-in user's link (most lists need it, or the
organization's); each event's link, to read or cancel that same meeting
later; an invitee's reschedule url, to send to a guest who wants another
time. Never look a meeting up by its name: use the link the taf kept.

**Start read-only.** "Tell me tomorrow's meetings and who's coming" needs
only the read tools. Booking or cancelling for a guest reaches someone
outside, so offer those only when the owner asked for them, and each one
waits for the owner's yes.

## Don't offer

- Rescheduling: Calendly's API has no way to move a meeting. A taf can send
  the guest's reschedule link, or cancel on the owner's yes.
- Deleting data (`CALENDLY_DELETE_INVITEE_DATA`,
  `CALENDLY_DELETE_SCHEDULED_EVENT_DATA`): they remove guests' or
  meetings' data across the whole organization.
- Organization admin (`CALENDLY_ORGANIZATION_INVITATION`,
  `CALENDLY_REMOVE_USER_FROM_ORGANIZATION`,
  `CALENDLY_DELETE_ORGANIZATION_MEMBERSHIP`): they change who is on the
  team.
- Changing meeting types and hours (`CALENDLY_CREATE_EVENT_TYPE`,
  `CALENDLY_UPDATE_EVENT_TYPE`, `CALENDLY_UPDATE_EVENT_TYPE_AVAILABILITY`):
  they change what every future guest can book.
- Webhook tools: the taf reads on its schedule, it doesn't listen.

## Access and plans

- **The API works on every Calendly plan, Free included,** except three
  tools that need its Enterprise plan: the activity log
  (`CALENDLY_LIST_ACTIVITY_LOG_ENTRIES`) and the two data deletions.
- **Rate:** per user, 500 calls a minute on a paid plan and 50 on Free.
- **Booking for a guest has its own, much lower cap:** 5 a day on Free or a
  trial; 10 a minute, 50 an hour and 100 a day on paid plans below
  Enterprise. A job that books for many guests a day won't fit Free.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in doesn't cover this tool, or the account's role doesn't allow it | "Calendly is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Refused, plan required | The tool needs Calendly's Enterprise plan | "This step needs Calendly's Enterprise plan. The taf will do without it." |
| Too many calls (429) | Over the per-minute limit, or the daily booking cap | "Calendly asked us to slow down. The taf will try again shortly; bookings over today's limit wait until tomorrow." |
| The sign-in no longer works (401) | The person changed their login, or removed access | "Calendly needs you to sign in again. Nothing was booked or cancelled." |
| Not found (404) | The meeting or meeting type was deleted, or belongs to someone else | "That meeting isn't in Calendly any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find a meeting never books a new one in
its place.
