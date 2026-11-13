# Outlook

Checked: 2026-11-13 · Sources: the connection service's Outlook tool
definitions as planque reads them (317 tools, [generated](generated/outlook.md)),
and Microsoft's mail and calendar API reference.

## How it connects

- **A sign-in to Microsoft**, held by the connection service: a personal
  Outlook account, or a work or school one. The catalog lists the tools and,
  with a known sign-in, puts the ones its access doesn't cover under
  `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked:
  ask for every tool the job needs, mail and calendar, in the first request.
- **One approval for the sign-in, one for the tools and the Card.**
- **A work account's admin** may have to approve the app before anyone in the
  business can sign in. The owner then sees "Need admin approval".

## Tools a job usually needs

Outlook has hundreds of tools; most jobs need these.

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `OUTLOOK_QUERY_EMAILS` | Lists messages matching a filter (sender, date, unread, folder) | read | message ids, conversation ids |
| `OUTLOOK_SEARCH_MESSAGES` | Finds messages by words | read | message ids |
| `OUTLOOK_GET_MESSAGE` | Reads one message | read | its conversation id |
| `OUTLOOK_LIST_MAIL_FOLDERS` | Lists the mailbox's folders | read | folder ids |
| `OUTLOOK_DOWNLOAD_OUTLOOK_ATTACHMENT` | Downloads one attachment | read | |
| `OUTLOOK_CREATE_DRAFT` | Writes a new draft; sends nothing | write | **the draft's id** |
| `OUTLOOK_CREATE_DRAFT_REPLY` | Writes a reply as a draft | write | **the draft's id** |
| `OUTLOOK_SEND_DRAFT` | Sends a draft as it is | write, outside | |
| `OUTLOOK_REPLY_EMAIL` | Sends a reply at once | write, outside | |
| `OUTLOOK_SEND_EMAIL` | Sends a new email at once | write, outside | |
| `OUTLOOK_MOVE_MESSAGE` | Moves a message to a folder | write | |
| `OUTLOOK_LIST_EVENTS` | Lists calendar events | read | event ids |
| `OUTLOOK_FIND_MEETING_TIMES` | Suggests times when people are free | read | |
| `OUTLOOK_CALENDAR_CREATE_EVENT` | Creates an event, with or without guests | write | **the event's id** |
| `OUTLOOK_UPDATE_CALENDAR_EVENT` | Changes an event | write | |

**Start with drafts.** "Draft, and let me check" is a read plus
`OUTLOOK_CREATE_DRAFT_REPLY`. Offer a sending tool only when the owner asked
for sending, and say every send waits for their yes.

**What to keep:** a draft's id, to send that same draft; a message's
conversation id, to reply in the right thread; each event's id.

**Guests get an invitation** when an event with attendees is created or
changed. The Card says so.

## Don't offer

- Deleting for good (`OUTLOOK_PERMANENT_DELETE_MESSAGE`) and deleting
  messages or folders, unless the job is about cleaning up and the owner
  sees each one.
- Mail rules (`OUTLOOK_CREATE_EMAIL_RULE`) and forwarding: they change what
  happens to every future email.
- Calendar sharing (`OUTLOOK_CREATE_ME_CALENDAR_PERMISSION` and the other
  permission tools): it decides who sees the owner's calendar.
- Teams chats, groups and planner tasks: they act beyond the owner's own
  mailbox and need their own job.

## Access and plans

- **Any Microsoft account works,** personal or business. Shared mailboxes
  and other people's calendars need access the owner has been given.
- **Sending caps are Microsoft's:** about 300 messages a day on a personal
  account and 10,000 recipients a day on most business plans, with no more
  than 30 messages a minute.
- **Attachments:** up to about 3 MB in one step; larger ones go by an upload
  session, which the tools do in parts.
- **Rate:** Microsoft limits each mailbox to about 10,000 calls in 10 minutes.
  A job reading thousands of messages should page and pause.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| "Need admin approval" at sign-in | A work account whose admin hasn't approved the app | "Your Microsoft admin needs to approve this connection first. Ask them, or connect a personal account." |
| Refused for missing access (403) | The sign-in didn't give this tool's access: reading only, say, and the tool sends | "Outlook is only partly connected: this step needs access you didn't give. Sign in again and allow it, or do without this step." |
| Too many calls (429) | Over Microsoft's limit for the mailbox | "Outlook asked us to slow down. The taf will try again in a few minutes." |
| Sending refused, limit reached | The day's sending cap | "Outlook's daily sending limit is reached. The rest goes tomorrow." |
| The message isn't found (404) | Deleted, or moved to a folder the taf doesn't read | "That email isn't there any more. The taf will skip it." |
| The sign-in no longer works (401) | Removed, password changed, or expired | "Outlook needs you to sign in again. Nothing was sent." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
