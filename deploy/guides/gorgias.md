# Gorgias

Checked: 2026-11-13 · Sources: the connection service's Gorgias tool
definitions as planque reads them (51 tools, [generated](generated/gorgias.md)),
and Gorgias's API reference.

## How it connects

- **A sign-in to the owner's Gorgias account** (`yourshop.gorgias.com`),
  held by the connection service. The catalog lists the tools; offer only what's
  under `tools`.
- **One approval for the connection, one for the tools and the Card.**
- **Not yet proven on a phone:** our gallery's refunds job leaves Gorgias
  out until it is. Offer it only to an owner who uses Gorgias, and say it's
  new.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GORGIAS_LIST_TICKETS` | Lists tickets, newest first, or in a view | read | ticket ids |
| `GORGIAS_LIST_VIEW_ITEMS` | Lists the tickets in one of the owner's views | read | ticket ids |
| `GORGIAS_GET_TICKET` | Reads one ticket | read | the customer's id |
| `GORGIAS_LIST_MESSAGES` | Reads a ticket's messages | read | |
| `GORGIAS_GET_CUSTOMER` | Reads a customer, with their shop data | read | |
| `GORGIAS_LIST_TAGS` | Lists the account's tags | read | tag names |
| `GORGIAS_ADD_TICKET_TAGS` | Adds tags to a ticket | write | |
| `GORGIAS_SEND_TICKET_MESSAGE` | Sends a message on a ticket: to the customer, or an internal note | write, outside if to the customer | the message's id |
| `GORGIAS_UPDATE_TICKET` | Changes a ticket's status or assignee | write | |
| `GORGIAS_CREATE_TICKET` | Opens a ticket | write | **the ticket's id** |

**A message is a note or a reply by its channel.** An internal note stays
with agents; any other channel goes to the customer. Start with notes, and
offer replies only when the owner asked, with each shown first or bounded by
the Card.

**What to keep:** each ticket's id; the tickets the taf already tagged or
answered, so it never acts twice.

## Don't offer

- Deleting tickets or customers (`GORGIAS_DELETE_TICKET`,
  `GORGIAS_DELETE_CUSTOMER`, `GORGIAS_DELETE_CUSTOMERS`).
- `GORGIAS_MERGE_CUSTOMERS`: it can't be undone.
- Rules, macros, views, teams and account settings: they change how every
  agent works.

## Access and plans

- **API access is on every paid plan;** Gorgias bills by ticket volume, and
  tickets a taf creates count.
- **Rate:** about 40 calls every 20 seconds for an account, shared with its
  other apps. Lists answer pages of up to 100.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | The connection's user lacks the right | "Your Gorgias role doesn't allow this. An admin can change it." |
| The sign-in no longer works (401) | Removed, or expired | "Gorgias needs you to sign in again." |
| Not found (404) | A deleted ticket | "That ticket isn't there any more." |
| Too many calls (429) | Over the account's rate | "Gorgias asked us to slow down. The taf will continue in a moment." |
| Refused (400) | A channel or field Gorgias didn't accept | Fix it and try once; else "Gorgias didn't accept that message. Nothing was sent." |

After any refusal the taf stops that step and says so.
