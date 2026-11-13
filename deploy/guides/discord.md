# Discord

Checked: 2026-11-13 · Sources: the connection service's Discord tool
definitions as planque reads them (23 tools), Discord's developer docs
(docs.discord.com: OAuth2, Users, Invites, Rate Limits, and Opcodes and
Status Codes).

## Discord or Discord Bot?

- **Discord** (this app) signs in as the owner's own Discord account. It can
  tell who the owner is and which servers they're in. **It can't read or
  send messages:** Discord's sign-in gives outside apps no access to the
  owner's channels or messages.
- **Discord Bot** ([discordbot.md](discordbot.md)) adds a bot to one server
  the owner runs. Use it for anything in a server: reading channels,
  posting, threads, events.
- A job that says "post to our Discord", "summarise our channel" or "answer
  in our server" is a Discord Bot job. Offer Discord only to check the
  owner's identity or list their servers, say to pick the server for the
  bot.

## How it connects

- **A sign-in to Discord**, held by the connection service; the taf never
  sees a password or token. The owner approves on Discord's own page.
  planque's catalog lists the tools and, with a known sign-in, puts the ones
  its access doesn't cover under `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request, so nothing is asked
  twice.

## Tools a job usually needs

Read means it changes nothing on Discord. The catalog gives no returns for
this app; the ids below are the ones Discord's reference documents.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `DISCORD_GET_MY_USER` | Reads the signed-in account | read | the user id, username |
| `DISCORD_LIST_MY_GUILDS` | Lists the servers the owner is in, up to 200 a page (a person can't be in more) | read | server ids and names, whether the owner owns each, their permissions there |
| `DISCORD_GET_MY_GUILD_MEMBER` | Reads the owner's membership in one server | read | their nickname, role ids, join date |
| `DISCORD_GET_OPENID_CONNECT_USERINFO` | Reads the owner's identity, with email if allowed | read | not documented |
| `DISCORD_LIST_MY_CONNECTIONS` | Lists accounts the owner linked to Discord | read | each connection's type, id and name |
| `DISCORD_GET_GUILD_WIDGET` | Reads a server's public widget, if the server turned it on | read | the server id, name, online count |

**What to keep:** the owner's user id and the server id they chose. Use a
kept id next time; don't look the server up by name.

**Start read-only.** Every tool a job needs here is a read. Nothing in this
app posts anything.

## Don't offer

- Writes that change the owner's linked roles or a developer app's
  settings: `DISCORD_UPDATE_USER_APPLICATION_ROLE_CONNECTION`,
  `DISCORD_DELETE_USER_APPLICATION_ROLE_CONNECTION`.
- Command permissions (`DISCORD_EDIT_APPLICATION_COMMAND_PERMISSIONS`):
  it changes who may use an app's commands in a server, an admin task.
- Purchases and subscriptions (`DISCORD_GET_CURRENT_USER_APPLICATION_ENTITLEMENTS`,
  `DISCORD_GET_SKU_SUBSCRIPTION`, `DISCORD_LIST_SKU_SUBSCRIPTIONS`): billing
  data of a developer's own app, no job of a taf.
- Developer plumbing with no owner job: `DISCORD_GET_GATEWAY`,
  `DISCORD_GET_PUBLIC_KEYS`, `DISCORD_LIST_STICKER_PACKS`,
  `DISCORD_GET_APPLICATION_COMMAND_PERMISSIONS`,
  `DISCORD_GET_BATCH_APPLICATION_COMMAND_PERMISSIONS`,
  `DISCORD_GET_USER_APPLICATION_ROLE_CONNECTION`.
- Any promise to read or send messages: no tool here can.

## Access and plans

- **Each tool needs its own part of the sign-in:** the profile, the server
  list, server membership, linked accounts and email are separate. A
  sign-in made for the profile alone doesn't list servers.
- **No paid plan** is needed for any tool above.
- **Rate:** Discord limits calls per route and answers "too many requests"
  with how long to wait. A job here makes a few calls a run and should never
  meet it.
- **Many refused calls in a row** (for missing access, or too fast) can get
  the caller blocked by Discord for a while: on a refusal, stop, don't
  retry in a loop.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, or the tool is under `unusable`) | The sign-in didn't include this part, the server list say | "Discord is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| The sign-in no longer works (401) | The owner removed the app in Discord's settings, or the sign-in expired | "Discord needs you to sign in again. Nothing was changed." |
| Too many calls (429) | Over Discord's limit for that call | "Discord asked us to slow down. The taf will try again shortly." |
| Not found (404), or an unknown server | The server or invite no longer exists, or the owner left it | "That Discord server isn't reachable any more. Is it still yours?" |
| A widget read is refused or empty | The server's widget is turned off | "That server doesn't share its public widget. The taf will skip it." |
| The job needs messages | This app can't read or post them | "To read or post in your server, connect Discord Bot instead. It adds a bot to your server." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
