# Telegram

Checked: 2026-11-14 · Sources: the connection service's Telegram tool
definitions as planque reads them (18 tools, [generated](generated/telegram.md)),
and Telegram's bot API reference.

## How it connects

- **By a bot's key.** The owner makes a bot with Telegram's own BotFather
  chat (`/newbot`), which answers a key; the owner pastes it. The taf then
  acts as **that bot**, never as the owner's own account.
- **A bot only sees chats it's in:** private chats with it, and the groups
  and channels it was added to. The owner adds the bot to a group, or writes
  to it, so the taf can find that chat's id.
- **One approval for the key, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `TELEGRAM_GET_ME` | Reads the bot's own name, to check the key | read | the bot's own details |
| `TELEGRAM_GET_UPDATES` | Reads new messages sent to the bot | read | each update's number and message, with **its chat's id** |
| `TELEGRAM_GET_CHAT` | Reads a chat's name and kind | read | the chat's details |
| `TELEGRAM_SEND_MESSAGE` | Posts a message in a chat | write, visible to others | the sent message, with its id |
| `TELEGRAM_SEND_DOCUMENT` / `TELEGRAM_SEND_PHOTO` | Posts a file or a photo | write, visible to others | the sent message, with its id |
| `TELEGRAM_EDIT_MESSAGE` | Edits a message the bot posted | write | the edited message |

**What to keep:** the chat's id, found once from `TELEGRAM_GET_UPDATES` after
the owner writes to the bot; and the last update's number, so the next read
asks for the ones after it and never answers twice.

## Don't offer

- `TELEGRAM_CREATE_CHAT_INVITE_LINK`: it lets new people into a group.
- `TELEGRAM_DELETE_MESSAGE` for others' messages: in a group it needs admin
  rights, and it removes what people said.
- `TELEGRAM_SET_MY_COMMANDS`: it changes the bot's menu for everyone.
- `TELEGRAM_FORWARD_MESSAGE` out of a group: it copies people's words
  elsewhere.

## Access and plans

- **No charge within the limits below.** Only paid broadcasts, sending more
  than about 30 messages a second, cost the bot Telegram Stars.
- **Rate:** about one message a second to one chat, 20 a minute to one group,
  and about 30 a second across all chats. Short bursts may pass; then
  Telegram refuses with a wait.
- **Groups:** with privacy mode on (the default, unless the bot is a group
  admin), a bot sees only commands meant for it and replies to its
  messages. The owner turns it off in BotFather if the job needs every
  message, then adds the bot to the group again.
- **Updates are kept for 24 hours.** A taf that reads less than daily misses
  messages.
- **A bot's key is the bot.** Anyone with it can post as the bot; the owner
  makes a new one in BotFather (`/token`) if it leaks.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key is refused | A wrong or replaced key | "Telegram didn't accept the bot's key. Copy it again from BotFather." |
| The bot can't reach the chat | The bot isn't in that chat, or no one wrote to it | "The bot can't reach that chat yet. Write to the bot first, or add it to the group." |
| The person blocked the bot | The person stopped the bot from writing to them | "That person blocked the bot, so the message wasn't sent." |
| Too many messages (429) | Over Telegram's rate | "Telegram asked us to slow down. The rest goes in a moment." |
| Reading new messages doesn't work | The bot has a webhook set, so its messages go elsewhere | "Another app receives this bot's messages. Use a bot of its own for this taf." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
