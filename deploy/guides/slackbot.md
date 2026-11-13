# Slackbot

Checked: 2026-11-13 · Sources: the connection service's Slackbot tool
definitions as planque reads them (105 tools, [generated](generated/slackbot.md)),
and Slack's API reference.

## What it is, beside Slack

**The same workspace, as a bot instead of as the owner.** With
[Slack](slack.md), the taf acts as the owner: its posts show under their
name. With Slackbot, it acts as an app of its own: its posts show under the
app's name, and it sees only the channels it's been added to.

- **Prefer Slackbot** for a taf that posts to a team (a daily summary, an
  alert): people see it's an app, not the owner.
- **Prefer Slack** for a taf that reads the owner's own messages and
  direct messages.

## How it connects

- **A sign-in to one Slack workspace,** which installs the app, held by the
  connection service. A workspace can require an admin to approve the app.
- **Then the owner adds the app to each channel** it should read or post in
  (`/invite @app` in the channel). Until then it can't see the channel.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `SLACKBOT_FIND_CHANNELS` | Finds channels by name | read | **the channel's id** |
| `SLACKBOT_LIST_CONVERSATIONS` | Lists the channels the app is in | read | channel ids |
| `SLACKBOT_FETCH_CONVERSATION_HISTORY` | Reads a channel's messages in a time window | read | message timestamps |
| `SLACKBOT_FETCH_MESSAGE_THREAD_FROM_A_CONVERSATION` | Reads a thread's replies | read | |
| `SLACKBOT_FIND_USER_BY_EMAIL_ADDRESS` | Finds a person by email | read | the person's id |
| `SLACKBOT_SEND_MESSAGE` | Posts a message, or a reply in a thread | write, visible to others | **the message's timestamp** |
| `SLACKBOT_SEND_EPHEMERAL_MESSAGE` | Posts a message only one person sees | write | |
| `SLACKBOT_SCHEDULE_MESSAGE` | Posts at a set time | write | the scheduled message's id |
| `SLACKBOT_UPLOAD_OR_CREATE_A_FILE_IN_SLACK` | Posts a file | write, visible to others | |

**What to keep:** the channel's id and each message's timestamp, as in
[Slack's guide](slack.md).

## Don't offer

- Workspace administration: user groups, channel creation and archiving,
  invites.
- Deleting messages or files other than the app's own; sharing a file
  publicly (`SLACKBOT_ENABLE_PUBLIC_SHARING_OF_A_FILE`).
- Calls, canvases and lists, unless the job is about them.

## Access and plans

- **On Slack's free plan,** only the last 90 days of messages can be read.
- **Rate:** about one message a second per channel; reading history is
  limited more tightly for apps outside Slack's marketplace.
- **A bot can't read direct messages between people,** nor channels it
  wasn't added to.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Not in the channel (`not_in_channel`, `channel_not_found`) | The app wasn't added to it | "Add the app to that channel first: type /invite and its name there." |
| Refused for missing access (`missing_scope`) | The install didn't give this access | "The Slack app is only partly set up: this step needs access it wasn't given. Sign in again and allow it." |
| Too many calls (`ratelimited`) | Over Slack's limit | "Slack asked us to slow down. The taf will continue in a minute." |
| Admin approval needed | The workspace restricts apps | "Your Slack admin needs to approve this app." |
| The install no longer works (`invalid_auth`, `token_revoked`) | Removed from the workspace | "The Slack app was removed. Sign in again to reinstall it." |

After any refusal the taf stops that step and says so.
