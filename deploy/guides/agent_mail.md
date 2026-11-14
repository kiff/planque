# Agent Mail

Checked: 2026-11-14 · Sources: the connection service's Agent Mail tool
definitions as planque reads them (93 tools), and Agent Mail's API
documentation (errors, rate limits, drafts, threads, labels).

## How it connects

- **By the owner's API key**, pasted once and held by the connection
  service; the taf never sees it. The key belongs to the owner's Agent Mail
  organization, and can be limited to some inboxes or some actions.
- **Agent Mail gives a taf its own inbox**, separate from the owner's
  personal mailbox. Mail it sends comes from that inbox's address.
- **A restricted key leaves tools out.** With a known sign-in the catalog
  puts the tools the key can't run under `unusable`. Ask for every tool the
  job needs in the first request, so nothing is asked twice.

## Tools a job usually needs

Read means it changes nothing in the inbox. Returns are what a taf keeps to
use next time.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `AGENT_MAIL_LIST_INBOXES` | Lists the organization's inboxes | read | each inbox's id and address |
| `AGENT_MAIL_GET_INBOX` | Reads one inbox's details | read | the inbox id |
| `AGENT_MAIL_LIST_THREADS` | Lists conversations in an inbox, filtered by label or date, a page at a time | read | thread ids, a page token |
| `AGENT_MAIL_GET_THREAD` | Reads the messages in one conversation, a page at a time | read | message ids, the thread id, a page token |
| `AGENT_MAIL_LIST_MESSAGES` | Lists messages in an inbox, filtered by label or date | read | message ids, thread ids, a page token |
| `AGENT_MAIL_GET_MESSAGE` | Reads one message | read | the message, its thread id, its labels |
| `AGENT_MAIL_GET_MESSAGE_ATTACHMENT` | Gets one attachment of a message | read | a download link for the file, which expires |
| `AGENT_MAIL_LIST_DRAFTS` / `AGENT_MAIL_GET_DRAFT` | Lists or reads drafts | read | draft ids |
| `AGENT_MAIL_CREATE_DRAFT` | Writes a draft, new or as a reply; sends nothing | write | the draft id |
| `AGENT_MAIL_UPDATE_DRAFT` | Changes a draft's content or labels | write | the draft id |
| `AGENT_MAIL_UPDATE_MESSAGE_LABELS` | Adds or removes labels on one message | write | the message's labels |
| `AGENT_MAIL_SEND_DRAFT` | Sends a draft as it is | write, outside | the sent message id |
| `AGENT_MAIL_REPLY_TO_MESSAGE` | Sends a reply to a message | write, outside | the sent message id |
| `AGENT_MAIL_SEND_MESSAGE` | Sends a new email at once | write, outside | the sent message id, its thread id |

**What to keep:** the inbox id the taf works in; a draft's id, to update or
send that draft later; a thread id, to reply in the right conversation. A
sent draft stops being a draft, so after sending keep the message id
instead. Never search for something the taf made: use the id it kept.

**Labels track the work.** A taf can mark a message as handled with a label
and list by that label next run, rather than re-reading everything.

**Start with drafts.** "Draft replies for me to check" is the read tools plus
`AGENT_MAIL_CREATE_DRAFT`. Offer a sending tool only when the owner asked for
sending, and say every send waits for their yes.

## Don't offer

- Deleting tools (messages, threads, drafts, files, inboxes): they can't be
  undone.
- Inbox, pod, domain and API key tools: they change the owner's account,
  not one job's mail.
- Webhook tools and allow or block lists: they change what reaches every
  inbox.
- `AGENT_MAIL_RELEASE_HELD_MESSAGE`: it lets through mail the app held
  back for a security review. That is the owner's call, in the app.
- `AGENT_MAIL_FORWARD_MESSAGE`: forwarding sends someone's mail to a third
  person; offer it only when the owner named the recipient.
- `AGENT_MAIL_SEND_EMAIL`: use `AGENT_MAIL_SEND_MESSAGE` or a draft instead.
- Calendar tools in Europe: the app's calendar isn't offered in its European
  region, so they fail there.

## Access and plans

- **A key can be narrowed** to some inboxes and some actions (read, draft,
  send). A key made for reading doesn't let the taf send.
- **Sending caps are per plan:** about 100 emails a day on the free plan,
  1,000 a day (10,000 a month) on the Developer plan, and 15,000 a day on the
  Startup plan. The free plan has 3 inboxes, Developer 10, Startup 150. A job
  that sends past the cap fails partway through.
- **A custom sending domain's** DNS records must be verified before mail can
  go out from it.
- **Attachments** sent inline are capped at about 6 MB a request.
- **Rate:** the app limits calls per organization and doesn't publish the
  numbers. When over, it says how long to wait, usually a second.
- **Sending can be suspended** for an inbox with too many bounces or
  complaints. Reading still works.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, missing permission) | The key doesn't allow this action, sending say | "Agent Mail is only partly connected: this step needs access your key doesn't give. Want to add a key that allows it, or do without this step?" |
| The key no longer works (401, unknown or expired key) | The key was revoked or has expired | "Agent Mail needs a new key. Nothing was sent." |
| Not found (404) | The message, draft or inbox was deleted, or the key can't see that inbox | "That email or inbox isn't there for this connection. The taf will skip it." |
| Send refused (403, message rejected) | A block list, a recipient missing from an allow list, or an attachment link that can't be fetched | "Agent Mail wouldn't send this one, because of a block or allow list, or an attachment it couldn't fetch. Nothing went out." |
| Domain not verified (403) | The sending domain's DNS records aren't checked yet | "Your sending domain isn't verified yet in Agent Mail. Finish its setup there, then the taf can send." |
| Inbox or account suspended (403) | Too many bounces or complaints, or the account is on hold | "Sending is paused on this Agent Mail inbox. Reading still works. Check the inbox's status in Agent Mail." |
| Too many calls, or the daily cap (429) | Over the rate limit or the plan's sending cap | "Agent Mail asked us to slow down, or today's sending limit is reached. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
