# Microsoft Teams

Checked: 2026-11-14 · Sources: the connection service's Microsoft Teams tool
definitions as planque reads them (237 tools), and Microsoft Teams' API
documentation (messages, channels, chats, teams, throttling, errors).

## How it connects

- **A sign-in to a work or school account**, held by the connection service;
  the taf never sees a password or token. Personal accounts can't send
  messages through the API at all, so a Teams job needs a work account.
- **The taf acts as the person who signed in.** It sees only the teams,
  channels and chats that person belongs to, and what it posts shows under
  their name.
- **With a known sign-in,** planque's catalog puts the tools its access
  doesn't cover under `unusable`, with why. A held sign-in is reused for a
  new taf when it covers the tools asked: ask for every tool the job needs
  in the first request.

## Tools a job usually needs

Read means it changes nothing in Teams. The catalog's entries give no
`returns` paths today; the ids below are the ones Teams' answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `MICROSOFT_TEAMS_LIST_USER_JOINED_TEAMS` | Lists the teams the person belongs to | read | team ids, names |
| `MICROSOFT_TEAMS_TEAMS_LIST_CHANNELS` | Lists a team's channels | read | channel ids, names, links |
| `MICROSOFT_TEAMS_GET_CHANNEL` | Reads one channel | read | the channel's id, its link |
| `MICROSOFT_TEAMS_TEAMS_LIST_CHANNEL_MESSAGES` | Lists the first messages of a channel's threads | read | message ids, links |
| `MICROSOFT_TEAMS_LIST_MESSAGE_REPLIES` | Lists the replies in one thread | read | reply ids, the thread's id |
| `MICROSOFT_TEAMS_GET_CHANNEL_MESSAGE` | Reads one channel message | read | the message's id, its link |
| `MICROSOFT_TEAMS_CHATS_GET_ALL_CHATS` | Lists the person's one-to-one, group and meeting chats | read | chat ids, links |
| `MICROSOFT_TEAMS_GET_CHAT_MESSAGE` | Reads one chat message | read | the message's id, its chat id |
| `MICROSOFT_TEAMS_TEAMS_POST_CHANNEL_MESSAGE` | Starts a new thread in a channel | write, outside | the new message's id and link |
| `MICROSOFT_TEAMS_POST_MESSAGE_REPLY` | Replies in a channel thread | write, outside | the reply's id and link |
| `MICROSOFT_TEAMS_TEAMS_POST_CHAT_MESSAGE` | Sends a message in an existing chat | write, outside | the new message's id |
| `MICROSOFT_TEAMS_SEND_MESSAGE_TO_SELF` | Sends a note to the person's own chat | write | not documented |

**What to keep:** the team id and channel id the job posts to; the id of a
thread's first message, to reply in that same thread; a chat's id, to write
in the same conversation. A message id is unique only inside its channel or
chat, so keep it with the channel or chat it belongs to. Never search for
something the taf posted: use the id it kept.

**Start read-only.** "Sum up this week's channel threads" needs only read
tools. Every post shows under the owner's name to colleagues, so offer a
posting tool only when the owner asked, and each post waits for their yes.
A note to self is the gentlest way to report.

## Don't offer

- Deleting, archiving or cloning teams and channels, and removing members:
  they change the organization's workspace for everyone. The same goes for
  creating or changing teams, channels and tabs (`MICROSOFT_TEAMS_CREATE_TEAM`,
  `MICROSOFT_TEAMS_CREATE_CHANNEL`, `MICROSOFT_TEAMS_UPDATE_TEAM`,
  `MICROSOFT_TEAMS_UPDATE_CHANNEL`, `MICROSOFT_TEAMS_ADD_TAB`,
  `MICROSOFT_TEAMS_UPDATE_TAB`, `MICROSOFT_TEAMS_DELETE_TAB`) and a channel's
  email address (`MICROSOFT_TEAMS_PROVISION_CHANNEL_EMAIL`,
  `MICROSOFT_TEAMS_REMOVE_CHANNEL_EMAIL`).
- Deleting or editing messages (`MICROSOFT_TEAMS_DELETE_SOFT_MESSAGE`,
  `MICROSOFT_TEAMS_UPDATE_CHANNEL_MESSAGE`,
  `MICROSOFT_TEAMS_UPDATE_CHAT_MESSAGE`): people may have read them already.
- Acting as the person beyond one message: new chats or meetings
  (`MICROSOFT_TEAMS_TEAMS_CREATE_CHAT`, `MICROSOFT_TEAMS_CREATE_MEETING`),
  alerts (`MICROSOFT_TEAMS_SEND_ACTIVITY_NOTIFICATION`), presence
  (`MICROSOFT_TEAMS_SET_PRESENCE*`) and hiding chats
  (`MICROSOFT_TEAMS_HIDE_CHAT_FOR_USER`).
- Membership and admin tools (adding people to teams, chats or groups;
  listing directory roles, licences or app grants).
- Meeting recordings and transcripts: sensitive, and not needed for
  messaging jobs.
- Change subscriptions: the taf reads on its schedule, it doesn't listen.
- Shift, time-off, task-board, file and workbook tools: they belong to
  other jobs and need access a messaging sign-in doesn't ask for.

## Access and plans

- **A work or school account** with Teams is needed. Sending isn't
  supported for personal accounts.
- **Rate:** about one post a second in any one channel
  or chat, and about one a second for one person. Reading a channel's
  messages is also held to about one call a second per channel. A job that
  posts to many places paces itself.
- **Polling:** Microsoft allows checking the same thing for changes only
  once a day. Schedule a job that reads Teams daily; one that must react
  within minutes isn't a fit for this app.
- **When throttled,** wait the time given, then try again more slowly.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, access denied (403) | The sign-in lacks this tool's access, the account has no Teams licence, or an access policy blocks it | "Microsoft Teams is only partly connected: this step needs access you didn't give, or your organization blocks it. Want to sign in again and allow it, or ask your IT admin?" |
| The sign-in no longer works (401) | The sign-in expired, or access was removed | "Microsoft Teams needs you to sign in again. Nothing was posted." |
| Not found (404) | The channel, chat or message was deleted, or the person left it | "That conversation isn't there any more for your account. The taf will skip it." |
| Too many calls (429) | Over the per-channel, per-chat or per-person rate | "Microsoft Teams asked us to slow down. The taf will try again in a moment." |
| Teams unavailable (503) | Maintenance or overload | "Microsoft Teams isn't answering right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find its thread asks, rather than
starting a new one.
