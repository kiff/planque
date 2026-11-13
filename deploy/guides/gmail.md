# Gmail

Checked: 2026-11-13 · Sources: the connection service's Gmail tool
definitions as planque reads them (63 tools), Google's Gmail API reference,
and test runs on a phone on 13 Nov 2026.

## How it connects

- **A sign-in to Google**, held by the connection service; the taf never sees
  a password or token. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **An account that already holds Gmail reuses that sign-in** for a new taf,
  with no browser, when it covers the tools asked. If the new taf needs a
  tool the held sign-in doesn't cover (sending, say, after a read-only
  sign-in), Google asks the owner again.
- **One approval for the sign-in, one for the tools and the Card.** Ask for
  every tool the job needs in the first request, so nothing is asked twice.

## Tools a job usually needs

Read means it changes nothing in the mailbox. Returns are what a taf keeps
to use next time; the catalog's `returns` gives their exact paths.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `GMAIL_FETCH_EMAILS` | Lists messages matching a search (sender, date, label, unread), a page at a time, up to 500 a page | read | message ids, thread ids, a page token |
| `GMAIL_LIST_THREADS` | Lists conversations matching a search | read | thread ids, a page token |
| `GMAIL_FETCH_MESSAGE_BY_THREAD_ID` | Reads every message in one conversation | read | message ids, headers, body |
| `GMAIL_FETCH_MESSAGE_BY_MESSAGE_ID` | Reads one message | read | the message, its thread id, its labels |
| `GMAIL_GET_ATTACHMENT` | Downloads one attachment of a message | read | the file's content |
| `GMAIL_LIST_LABELS` | Lists the mailbox's labels | read | label ids and names |
| `GMAIL_LIST_DRAFTS` / `GMAIL_GET_DRAFT` | Lists or reads drafts | read | draft ids |
| `GMAIL_CREATE_EMAIL_DRAFT` | Writes a draft, in a conversation or new; sends nothing | write | the draft id, its message id |
| `GMAIL_UPDATE_DRAFT` | Replaces a draft's content | write | the draft id |
| `GMAIL_ADD_LABEL_TO_EMAIL` | Adds or removes labels on one message | write | the message's labels |
| `GMAIL_SEND_DRAFT` | Sends a draft as it is | write, outside | the sent message id |
| `GMAIL_REPLY_TO_THREAD` | Sends a reply in a conversation | write, outside | the sent message id |
| `GMAIL_SEND_EMAIL` | Sends a new email at once | write, outside | the sent message id |

**What to keep:** a draft's id, to update or send that same draft later; a
message's thread id, to reply in the right conversation. Never search for
something the taf made: use the id it kept.

**Start read-only.** "Draft, and let me check before anything is sent" is a
read tool plus `GMAIL_CREATE_EMAIL_DRAFT`. Offer a sending tool only when
the owner asked for sending, and say every send waits for their yes.

## Don't offer

- Deleting tools (`GMAIL_DELETE_MESSAGE`, `GMAIL_BATCH_DELETE_MESSAGES`,
  `GMAIL_DELETE_THREAD`): they bypass the bin and can't be undone.
- Settings tools (filters, forwarding, vacation replies, send-as aliases,
  IMAP or POP): they change how the whole mailbox behaves. Forwarding in
  particular can send every future email elsewhere.
- Deprecated tools (`GMAIL_LIST_MESSAGES`, `GMAIL_REMOVE_LABEL`): use their
  replacements above.

## Access and plans

- **Reading and drafting** need Gmail read access plus the right to write
  drafts; sending needs send access. A sign-in made for reading doesn't let
  the taf send.
- **Daily sending caps are Google's:** about 500 recipients a day on a
  personal account, and about 2,000 a day on a business (Workspace) account.
  A job that sends to more each day than that fails partway through.
- **Attachments:** up to about 25 MB in a draft or reply. Larger files go as
  a Drive link.
- **A business account's admin** can block outside apps from Gmail. Then no
  sign-in works until the admin allows it.
- **Rate:** Google limits calls per user per second, and the connection
  service's shared sign-in shares a daily quota. A job reading thousands of
  messages a run should page and pause.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | The sign-in didn't give this tool's access: reading only, say, and the tool sends | "Gmail is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Refused, "rate limit exceeded" (403) or too many calls (429) | Too many calls in a short time | "Gmail asked us to slow down. The taf will try again in a few minutes." |
| The sign-in no longer works (401) | The owner removed access, changed their password, or the sign-in expired | "Gmail needs you to sign in again. Nothing was sent." |
| "App is blocked" at sign-in | A business admin, or Google, blocks this access for outside apps | "Your Google admin blocks this. Ask them to allow it, or connect another account." |
| A message or draft isn't found (404) | It was deleted, or sent (a sent draft is no longer a draft) | "That email isn't there any more. The taf will skip it." |
| Daily sending cap reached | Google's cap for the account | "Gmail's daily sending limit is reached. The rest goes tomorrow." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
