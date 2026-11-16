# Gmail

Checked: 2026-11-15 · Sources: the connection service's Gmail tool
definitions as planque reads them (61 tools), Google's Gmail API reference,
test runs on a phone on 13 and 14 Nov 2026, and planque's scope choice for
the managed sign-in (`docs/managed-connectors.md`).
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **A sign-in to Google**, held by the connection service; the taf never sees
  a password or token. The grant's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **An account that already holds Gmail reuses that sign-in** for a new taf,
  with no browser, when it covers the tools asked. Through the managed
  sign-in it always does (full mail access, below).
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
the owner asked for sending, and say every send waits for their yes. The
taf's tools and its Card are what keep a job to reading: the sign-in itself
can do more (below).

## Don't offer

- Deleting tools (`GMAIL_DELETE_MESSAGE`, `GMAIL_BATCH_DELETE_MESSAGES`,
  `GMAIL_DELETE_THREAD`): they bypass the bin and can't be undone.
- Settings tools (filters, forwarding, vacation replies, send-as aliases,
  IMAP or POP): they change how the whole mailbox behaves. Forwarding in
  particular can send every future email elsewhere.

## Access and plans

- **Through the managed sign-in, Google's page asks for full mail access**
  ("read, compose, send and permanently delete"), whatever the taf's tools:
  that sign-in offers no narrower Gmail access, not read-only, compose or
  modify. A read-only job's sign-in can, as a sign-in, also send and delete.
  - The approval says so: "Gmail will ask for more access than this taf
    uses". Coach says it plainly before connecting, and never "only what this
    taf needs".
  - What the taf may do is still only its tools and its Card. A taf with no
    sending tool can't send, whatever the sign-in allows.
  - Read-only Gmail needs a sign-in whose own access includes it, which the
    managed one doesn't. Until then, don't promise it.
- **Never contacts, birthday, addresses or phone numbers:** a Gmail job asks
  for none of them.
- The tools are named in the connect: a connect that names none asks for the
  app's widest access, so the taf's planned tools are always named when it
  connects.
- **Daily sending caps are Google's:** about 500 recipients a day on a
  personal account. A business (Workspace) account can send about 2,000
  messages a day, to at most 10,000 recipients in all, fewer on a trial
  account. A job that sends past its cap fails partway through.
- **Attachments:** up to about 25 MB in a draft or reply on a personal
  account; on a business account the admin sets the size. Larger files go
  as a Drive link.
- **A business account's admin** can block outside apps from Gmail. Then no
  sign-in works until the admin allows it.
- **Rate:** Google limits calls per user and per app, a minute at a time. A
  job reading thousands of
  messages a run should page and pause.

## When it fails

The rows are Gmail's own refusals. planque's own, and an expired sign-in,
are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, "insufficient scopes") | The sign-in didn't give this tool's access: an older sign-in, or one made another way than the managed one | "Gmail is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Refused, "rate limit exceeded" (403) or too many calls (429) | Too many calls in a short time | "Gmail asked us to slow down. The taf will try again in a few minutes." |
| Gmail's sign-in has expired (401) | The owner removed access, or the sign-in expired | "Gmail needs you to sign in again. Nothing was sent." |
| At sign-in: "App is blocked" | A business admin, or Google, blocks this access for outside apps | "Your Google admin blocks this. Ask them to allow it, or connect another account." |
| Refused (404), message or draft not found | It was deleted, or sent (a sent draft is no longer a draft) | "That email isn't there any more. The taf will skip it." |
| Refused (400), `FAILED_PRECONDITION`, "Precondition check failed" | The connected Google account has no Gmail mailbox it can read: a work account with Gmail turned off, or one made without Gmail. Signing in again with the same account doesn't fix it | "The Google account connected for Gmail has no mailbox I can read. Connect Gmail again with the account whose inbox you want." |
| Refused, daily sending cap reached | Google's cap for the account | "Gmail's daily sending limit is reached. The rest goes tomorrow." |
