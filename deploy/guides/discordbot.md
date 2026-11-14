# Discord Bot

Checked: 2026-11-14 · Sources: the connection service's Discord Bot tool
definitions as planque reads them (157 tools), Discord's developer docs
(docs.discord.com: OAuth2, Messages, Guilds, Users, Gateway intents, Rate
Limits, and Opcodes and Status Codes).

## Discord Bot or Discord?

- **Discord Bot** (this app) adds a bot to one server the owner runs. The
  bot reads and posts in that server's channels as itself, not as the
  owner. Use it for anything inside a server.
- **Discord** ([discord.md](discord.md)) signs in as the owner's own
  account and only reads their profile and server list. It can't read or
  send messages. Offer it only when the job needs the owner's identity.
- The bot posts under its own name, never as the owner.

## How it connects

- **Adding the bot to a server**, a sign-in held by the connection service;
  the taf never sees a token. On Discord's page the owner picks the server;
  they need the right to add apps there.
- **The bot sees what its role allows.** Channels its role can't view are
  invisible to it, and posting needs the right to send in that channel. A
  server admin can change this in the server's settings.
- **Ask for every tool the job needs in the first request**, so nothing is
  asked twice.

## Tools a job usually needs

Read means it changes nothing in the server. The catalog gives no returns
for this app; the ids below are the ones Discord's reference documents.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `DISCORDBOT_GET_CONNECTED_GUILD` | Reads the server this connection is for | read | not documented |
| `DISCORDBOT_LIST_GUILD_CHANNELS` | Lists the server's channels (not threads) | read | channel ids, names, types |
| `DISCORDBOT_GET_ACTIVE_GUILD_THREADS` | Lists the server's open threads | read | thread ids |
| `DISCORDBOT_LIST_MESSAGES` | Reads a channel's messages, up to 100 a page | read | message ids, authors, times |
| `DISCORDBOT_GET_MESSAGE` | Reads one message | read | the message id, its channel id |
| `DISCORDBOT_SEARCH_GUILD_MESSAGES` | Searches the server's messages, up to 25 a page | read | message ids, a total that may be approximate |
| `DISCORDBOT_CREATE_MESSAGE` | Posts a message, up to 2,000 characters | write, visible to the server | the message id, its channel id |
| `DISCORDBOT_UPDATE_MESSAGE` | Edits a message the bot posted | write, visible | the message id |
| `DISCORDBOT_CREATE_THREAD_FROM_MESSAGE` | Opens a thread on a message | write, visible | not documented |
| `DISCORDBOT_ADD_MY_MESSAGE_REACTION` | Adds the bot's reaction to a message | write, visible | nothing (an empty answer) |

**What to keep:** the server id, the channel id the owner chose, and the id
of each message the bot posted, to edit it or reply in a thread later. Look
a channel up by id, never by name each run.

**Start read-only.** "Summarise the channel and show me" is reads only.
Offer posting only when the owner asked for it, and say each post waits
for their yes: the whole channel sees it.

## Don't offer

- Deleting: `DISCORDBOT_DELETE_MESSAGE`, `DISCORDBOT_BULK_DELETE_MESSAGES`,
  `DISCORDBOT_DELETE_CHANNEL`, `DISCORDBOT_DELETE_GUILD`, and every other
  `DELETE_` tool.
- People: `DISCORDBOT_BAN_USER_FROM_GUILD`, `DISCORDBOT_BULK_BAN_USERS_FROM_GUILD`,
  `DISCORDBOT_DELETE_GUILD_MEMBER` (a kick), `DISCORDBOT_PRUNE_GUILD`,
  `DISCORDBOT_UPDATE_GUILD_MEMBER` (timeouts, nicknames, roles),
  `DISCORDBOT_ADD_GUILD_MEMBER`, `DISCORDBOT_ADD_GUILD_MEMBER_ROLE`,
  `DISCORDBOT_UPDATE_VOICE_STATE`.
- Roles and permissions: `DISCORDBOT_CREATE_GUILD_ROLE`,
  `DISCORDBOT_UPDATE_GUILD_ROLE`, `DISCORDBOT_REORDER_GUILD_ROLES`,
  `DISCORDBOT_SET_CHANNEL_PERMISSION_OVERWRITE`, and the command-permission
  tools.
- Server settings: `DISCORDBOT_UPDATE_GUILD`, `DISCORDBOT_UPDATE_CHANNEL`,
  `DISCORDBOT_CREATE_GUILD_CHANNEL`, and the auto-moderation, onboarding,
  welcome, widget, template, emoji and sticker tools.
- Webhooks and invites: `DISCORDBOT_CREATE_WEBHOOK`, the `EXECUTE_` webhook
  tools, `DISCORDBOT_CREATE_CHANNEL_INVITE`, `DISCORDBOT_INVITE_REVOKE`. A
  webhook is a standing way to post that outlives the taf.
- Direct messages (`DISCORDBOT_CREATE_DM`): they reach members one by one,
  who didn't ask the taf for them.
- `DISCORDBOT_LEAVE_GUILD`, the application-command and interaction tools,
  and the audit log and ban list.

## Access and plans

- **Message text needs a special allowance.** Without Discord's
  message-content access, messages come back with empty text, except
  messages that mention the bot, DMs with it and its own posts. Search
  needs that access too.
- **Reading history** needs the bot's role to see the channel and read its
  history; without that, reads come back empty.
- **From 16 Nov 2026** the channel list leaves out channels the bot can't
  see.
- **Rate:** a bot may make about 50 calls a second in all, and each kind
  of call has its own limits too. A refusal says how long to wait.
- **Many refused calls** (10,000 refused or too-fast calls in 10 minutes)
  get the caller blocked for a while. On a refusal, stop;
  don't retry in a loop.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, missing access (403) | The bot's role can't see that channel | "The bot can't see that channel. In your server's settings, let its role view it, or pick another channel." |
| Refused, missing permissions (403) | The bot can see the channel but not post, read history or do this | "The bot isn't allowed to do that in this channel. A server admin can give its role that right." |
| Messages come back with empty text | No message-content access | "Discord gave the taf those messages without their text. It can only read messages that mention the bot." |
| Search answers "not ready" (202) | Discord is still indexing the server | Wait the time given and try once more; tell the owner only if it fails again. |
| Too many calls (429), or a slow-mode refusal | Over Discord's limit for that channel | "Discord asked us to slow down. The taf will continue in a moment." |
| Not found (404): unknown channel, message or server | It was deleted, or the bot was removed from the server | "That channel or message isn't there any more. Is the bot still in your server?" |
| Refused (400) | A message over 2,000 characters, empty, or a bad field | Fix it and try once; else "Discord didn't accept that post. Nothing was posted." |
| The sign-in no longer works (401) | The bot's sign-in is missing or no longer valid | "Discord Bot needs to be connected again. Nothing was posted." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
