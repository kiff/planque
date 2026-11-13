# Telegram

Checked: 2026-11-13 · Sources: the connection service's Telegram tool
definitions as planque reads them (18 tools, [generated](generated/telegram.md)),
and Telegram's bot API reference.

## How it connects

- **By a bot's key.** The owner makes a bot with Telegram's own BotFather
  chat (`/newbot`), which answers a key; the owner pastes it. The taf then
  acts as **that bot**, never as the owner's own account.
- **A bot only sees chats it's in.** The owner adds the bot to a group, or
  starts a chat with it, before it can read or post there.
- **One approval for the key, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `TELEGRAM_GET_ME` | Reads the bot's own name, to check the key | read | |
| `TELEGRAM_GET_UPDATES` | Reads new messages sent to the bot | read | **each chat's id**, the last update's number |
| `TELEGRAM_GET_CHAT` | Reads a chat's name and kind | read | |
| `TELEGRAM_SEND_MESSAGE` | Posts a message in a chat | write, visible to others | the message's id |
| `TELEGRAM_SEND_DOCUMENT` / `TELEGRAM_SEND_PHOTO` | Posts a file or a photo | write, visible to others | |
| `TELEGRAM_EDIT_MESSAGE` | Edits a message the bot posted | write | |

**What to keep:** the chat's id (a number; groups' are negative), found once
from `TELEGRAM_GET_UPDATES` after the owner writes to the bot; and the last
update's number, so the next read starts after it and never answers twice.

**A bot can't start a chat.** It can only write to someone who wrote to it
first, or in a group it's in.

## Don't offer

- `TELEGRAM_CREATE_CHAT_INVITE_LINK`: it lets new people into a group.
- `TELEGRAM_DELETE_MESSAGE` for others' messages: it needs admin rights and
  removes what people said.
- `TELEGRAM_SET_MY_COMMANDS`: it changes the bot's menu for everyone.
- `TELEGRAM_FORWARD_MESSAGE` out of a group: it copies people's words
  elsewhere.

## Access and plans

- **Free.** Telegram doesn't charge for bots.
- **Rate:** about one message a second to one chat, 20 a minute to one group,
  and 30 a second across all chats.
- **Groups:** with privacy mode on (the default), a bot sees only messages
  that mention it or answer it. The owner turns it off in BotFather if the
  job needs every message.
- **Updates are kept for 24 hours.** A taf that reads less than daily misses
  messages.
- **A bot's key is the bot.** Anyone with it can post as the bot; the owner
  can make a new one in BotFather at any time.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key is refused (401) | A wrong or replaced key | "Telegram didn't accept the bot's key. Copy it again from BotFather." |
| "Chat not found" (400) | The bot isn't in that chat, or no one wrote to it | "The bot can't reach that chat yet. Write to the bot first, or add it to the group." |
| "Bot was blocked by the user" (403) | The person blocked the bot | "That person blocked the bot, so the message wasn't sent." |
| Too many messages (429) | Over Telegram's rate | "Telegram asked us to slow down. The rest goes in a moment." |
| "Conflict" (409) | Another program reads the same bot's messages | "Another app is reading this bot's messages. Use a bot of its own for this taf." |

After any refusal the taf stops that step and says so.
