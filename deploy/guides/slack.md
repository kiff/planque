# Slack

Checked: 2026-11-14 · Sources: the connection service's Slack tool
definitions as planque reads them (165 tools, [generated](generated/slack.md)),
and Slack's API reference.

## How it connects

- **A sign-in to one Slack workspace**, held by the connection service. The
  catalog lists the tools and, with a known sign-in, puts the ones its access
  doesn't cover under `unusable`, with why. Offer only what's under `tools`.
- **One workspace per sign-in.** An owner in two workspaces signs in to each.
- **A workspace can require an admin's approval** for outside apps. The
  owner then sees a request to send to the admin instead of an "Allow"
  button.
- **The taf acts as the owner**: what it posts shows under their name, and
  it sees only the channels they're in.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `SLACK_FIND_CHANNELS` | Finds channels by name | read | **the channel's id** |
| `SLACK_LIST_ALL_CHANNELS` | Lists the channels the owner can see | read | channel ids |
| `SLACK_FETCH_CONVERSATION_HISTORY` | Reads a channel's messages in a time window | read | each message's timestamp |
| `SLACK_FETCH_MESSAGE_THREAD_FROM_A_CONVERSATION` | Reads the replies in one thread | read | |
| `SLACK_SEARCH_MESSAGES` | Finds messages by words | read | channel ids, timestamps |
| `SLACK_FIND_USERS` / `SLACK_FIND_USER_BY_EMAIL_ADDRESS` | Finds a person | read | **the person's id** |
| `SLACK_RETRIEVE_MESSAGE_PERMALINK_URL` | Gives a message's link | read | the link |
| `SLACK_SEND_MESSAGE` | Posts a message, or a reply in a thread | write, visible to others | **the message's timestamp** |
| `SLACK_UPDATES_A_SLACK_MESSAGE` | Edits a message the taf posted | write | |
| `SLACK_SCHEDULE_MESSAGE` | Posts a message at a set time | write | the scheduled message's id |
| `SLACK_ADD_REACTION_TO_AN_ITEM` | Adds an emoji reaction | write, visible to others | |

**What to keep:** a channel's id, not its name, which people rename; a
message's timestamp (`ts`), which is its id: it's what a thread reply, an
edit or a link needs.

**Every post is seen by the people in the channel.** The Card names the
channel, and the first post is shown to the owner before it goes. A summary
posted to the owner alone goes as a direct message to themselves.

## Don't offer

- Workspace administration: inviting or removing people, user groups,
  default channels, resetting sessions, converting or archiving channels.
- Deleting messages (`SLACK_DELETES_A_MESSAGE_FROM_A_CHAT`) other than the
  taf's own.
- Setting the owner's presence or profile: it speaks for them in a way they
  didn't ask for.
- Reading every channel to search for something: ask the owner which
  channels matter.

## Access and plans

- **On Slack's free plan, only the last 90 days** of messages and files can
  be read; older ones are deleted after a year.
- **Rate:** Slack limits each method, roughly from 1 to 100 calls a minute.
  Posting is about one message a second per channel. Reading history is
  limited more tightly for apps outside Slack's marketplace, so a job
  reading a busy channel should read a small window each run.
- **Search** needs search access, which some workspaces don't allow for
  outside apps.
- **Private channels and direct messages** are visible only if the owner is
  in them, and only with the access for them.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (`missing_scope`) | The sign-in didn't give this tool's access | "Slack is only partly connected: this step needs access you didn't give. Sign in again and allow it, or do without this step." |
| Channel not found, or not in it (`channel_not_found`, `not_in_channel`) | The owner isn't in the channel, it's private, or archived | "I can't post in that channel. Make sure you're in it, or pick another." |
| Too many calls (429, `ratelimited`) | Over Slack's limit for that method | "Slack asked us to slow down. The taf will continue in a minute." |
| Admin approval needed at sign-in | The workspace restricts outside apps | "Your Slack admin needs to approve this connection. Ask them, or use another workspace." |
| The sign-in no longer works (`invalid_auth`, `token_revoked`) | Removed, or the owner left the workspace | "Slack needs you to sign in again. Nothing was posted." |
| Message cut short | Slack cuts a message far over its length guide (about 4,000 characters of text) | Keep messages short; split a long report into several. |

After any refusal the taf stops that step and says so.
