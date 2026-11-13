# Intercom

Checked: 2026-11-13 · Sources: the connection service's Intercom tool
definitions as planque reads them (138 tools), and Intercom's developer
documentation (rate limits, error answers, OAuth scopes, and the reference
pages for replying, contacts and tickets).

## How it connects

- **A sign-in to Intercom**, made on Intercom's own page and held by the
  connection service; the taf never sees a token. planque's catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`, with why.
- **What the sign-in may do is set by the app it signs in to,** not picked
  by the person at sign-in: Intercom has the app's maker choose its
  permissions. So a tool can be offered and still be refused by Intercom.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in Intercom. The catalog gives no `returns`
paths for Intercom today; the ids below are those Intercom's reference pages
show in the answer, and "not documented" means the pages read didn't show
the answer's fields.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `INTERCOM_IDENTIFY_AN_ADMIN` | Reads the teammate the sign-in acts as | read | not documented |
| `INTERCOM_LIST_CONVERSATIONS` | Lists conversations, a page at a time | read | not documented |
| `INTERCOM_SEARCH_CONVERSATIONS` | Finds conversations by state, date, contact and the like | read | not documented |
| `INTERCOM_GET_CONVERSATION` | Reads one conversation with its messages | read | the conversation's id, its state, its message parts |
| `INTERCOM_SEARCH_CONTACTS` | Finds people by email or other details | read | not documented |
| `INTERCOM_GET_A_CONTACT` | Reads one person | read | the contact's id, their own system's id, email, role |
| `INTERCOM_LIST_TAGS` | Lists the workspace's tags | read | not documented |
| `INTERCOM_SEARCH_TICKETS` | Finds tickets | read | not documented |
| `INTERCOM_GET_TICKET` | Reads one ticket | read | the ticket's id, its state |
| `INTERCOM_REPLY_TO_CONVERSATION` | Replies in a conversation, as a teammate | write, outside | the conversation's id and state |
| `INTERCOM_CREATE_A_NOTE` | Adds a note to a person, seen by the team only | write | not documented |
| `INTERCOM_ATTACH_TAG_TO_CONVERSATION` | Tags a conversation | write | not documented |
| `INTERCOM_ASSIGN_CONVERSATION` | Hands a conversation to a teammate or team | write | not documented |
| `INTERCOM_CLOSE_CONVERSATION` | Closes a conversation, with an optional last message | write, outside if a message is sent | not documented |
| `INTERCOM_CREATE_TICKET` | Opens a ticket for one or more people | write | the ticket's id, its state |

**What to keep:** the teammate id the sign-in acts as, since replies, tags
and closing all need it as their author; each conversation's id, to reply in
the same one later; a ticket's `id`, not its `ticket_id`: Intercom says the
second is the number shown in its inbox and isn't for the API. Never search
for something the taf made: use the id it kept.

**Start read-only.** "Tell me which conversations are waiting and draft
answers for me" needs only the read tools. A reply reaches the customer, so
offer replying or closing with a message only when the owner asked for it,
and each reply waits for the owner's yes.

## Don't offer

- Deleting or redacting (`INTERCOM_DELETE_CONVERSATION`,
  `INTERCOM_REDACT_CONVERSATION_CONTENT`, `INTERCOM_DELETE_CONTACT`,
  `INTERCOM_DELETE_TICKET`, `INTERCOM_DELETE_A_TAG_DELETE_TAG`): permanent.
- Blocking a person (`INTERCOM_BLOCK_CONTACT`): it stops them reaching the
  business at all.
- Outbound messages (`INTERCOM_SEND_OUTBOUND_MESSAGE`): they start new
  conversations with customers who didn't write in.
- Team and setup changes: teammates' away status, ticket types and their
  fields, and email subscriptions.

## Access and plans

- **Some parts of the API come only with some Intercom plans.** A tool for a
  feature the workspace's plan lacks is refused, and the taf can't change
  that.
- **Rate:** Intercom counts calls per minute, spread over six ten-second
  windows. An app offered to many workspaces gets 10,000 calls a minute
  of its own, and a workspace allows 25,000 a minute across all its apps.
- **Fast replies to one conversation** can be slowed on their own; Intercom
  then says how many seconds to wait.
- **Two people with the same email** make a lookup by email fail. Use the
  contact's id instead.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, not allowed (403) | The sign-in's app doesn't have this permission | "Intercom is only partly connected: this step needs access the connection doesn't have. The taf will do without it." |
| Refused, not on this plan (403, `api_plan_restricted`) | The workspace's Intercom plan doesn't include this feature | "This needs a feature your Intercom plan doesn't include. The taf will do without it." |
| The sign-in no longer works (401) | Access was removed, or the token is no longer valid | "Intercom needs you to sign in again. Nothing was sent." |
| Not found (404) | The conversation, person or ticket was deleted | "That isn't in Intercom any more. The taf will skip it." |
| Several people match (409) | More than one contact has that email | "Two people in Intercom share that email. Tell the taf which one, or it will skip them." |
| Rejected as invalid (422) | A value is well formed but doesn't fit, such as an unknown teammate | "Intercom didn't accept that change. The taf has stopped and will show you what it tried." |
| Too many calls (429) | Over the per-minute limit, or replies too fast in one conversation | "Intercom asked us to slow down. The taf will try again in a moment." |
| Intercom unavailable (5xx) | Trouble on Intercom's side | "Intercom isn't answering right now. The taf will try once more later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find a conversation never starts a new
one in its place.
