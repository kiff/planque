# Freshdesk

Checked: 2026-11-14 · Sources: the connection service's Freshdesk tool
definitions as planque reads them (184 tools, [generated](generated/freshdesk.md)),
and Freshdesk's API reference.
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **Through the connection service, with the owner's Freshdesk key and
  domain** (`yourcompany.freshdesk.com`). The key is in Freshdesk under the
  owner's profile settings. The taf acts as that agent, with that agent's
  rights.
- The catalog lists the tools; offer only what's under `tools`.
- **One approval for the connection, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `FRESHDESK_GET_TICKETS` | Lists tickets, newest first, or by a filter | read | ticket ids |
| `FRESHDESK_GET_SEARCH` | Finds tickets by fields (status, priority, tag) | read | ticket ids |
| `FRESHDESK_VIEW_TICKET` | Reads one ticket | read | the requester's id |
| `FRESHDESK_LIST_ALL_TICKET_CONVERSATIONS` | Reads a ticket's replies and notes | read | |
| `FRESHDESK_GET_CONTACT` | Reads a customer | read | |
| `FRESHDESK_ADD_NOTE_TO_TICKET` | Adds a private note, seen by agents only | write | the note's id |
| `FRESHDESK_REPLY_TO_TICKET` | Replies to the customer | write, outside | |
| `FRESHDESK_UPDATE_TICKET` | Changes a ticket's status, priority, tags or agent | write | |
| `FRESHDESK_CREATE_TICKET` | Opens a ticket | write | **the ticket's id** |

**Start with private notes.** "Draft a reply for an agent to check" is a
note. A reply goes to the customer at once and can't be taken back: offer
`FRESHDESK_REPLY_TO_TICKET` only when the owner asked for it, with each reply
shown first or bounded by the Card.

**What to keep:** each ticket's id; a ticket the taf already answered, so it
never answers twice.

## Don't offer

- Deleting tickets, contacts, conversations or articles
  (`FRESHDESK_DELETE_TICKET`, `FRESHDESK_DELETE_MULTIPLE_TICKETS`,
  `FRESHDESK_HARD_DELETE_CONTACT` and the others).
- Bulk updates (`FRESHDESK_BULK_UPDATE_TICKETS`) unless the owner saw the list.
- Merging tickets (`FRESHDESK_MERGE_TICKETS`): it can't be undone.
- Administration: agents, groups, roles, automation rules, SLA policies,
  ticket fields, mailboxes, forums.

## Access and plans

- **The free plan has no API access;** paid plans do.
- **Rate:** per minute, by plan, from 100 calls a minute on Growth to 700 on
  Enterprise, for the whole account, shared with every other app on it. A list answers 30
  tickets a page, at most 100.
- **Search** answers at most 300 results (10 pages of 30).

## When it fails

The rows are Freshdesk's own answers: a refusal reaches the taf as
"Freshdesk refused it (it answered …): …", with Freshdesk's reason.
planque's own refusals, and an expired sign-in, are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused (403) | The plan has no API access, or the agent lacks the right | "Your Freshdesk plan or role doesn't allow this. Check the plan, or connect with an admin's key." |
| Freshdesk's sign-in has expired (401) | Wrong key or domain, or the key was reset | "Freshdesk didn't accept the key. Check the domain, then tap Connect on this taf and paste the key from your profile again." |
| Refused (404), not found | A deleted ticket, or a wrong domain | "That ticket isn't there any more." |
| Refused (429), too many calls | Over the account's per-minute limit | "Freshdesk asked us to slow down. The taf will continue in a minute." |
| Refused (400) | A field value Freshdesk doesn't accept (a status number, a missing email) | Fix it and try once; else "Freshdesk didn't accept that change. Nothing was updated." |
