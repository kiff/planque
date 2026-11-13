# Zendesk

Checked: 2026-11-13 · Sources: the connection service's Zendesk tool
definitions as planque reads them (480 tools, [generated](generated/zendesk.md)),
and Zendesk's Support API reference.

## How it connects

- **A sign-in to the owner's Zendesk** (`yourcompany.zendesk.com`), held by
  the connection service. The taf acts as the owner's agent account, with its
  rights. The catalog lists the tools and puts the ones its access doesn't
  cover under `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**
- **In the gallery:** "Acknowledge every ticket" uses Zendesk.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `ZENDESK_LIST_ZENDESK_TICKETS` | Lists tickets, newest first | read | ticket ids |
| `ZENDESK_SEARCH_ZENDESK` | Finds tickets, users or organisations by fields (`status:new`) | read | ticket ids |
| `ZENDESK_GET_ZENDESK_TICKET_BY_ID` | Reads one ticket | read | the requester's id |
| `ZENDESK_GET_TICKET_COMMENTS` | Reads a ticket's conversation | read | |
| `ZENDESK_GET_USER` | Reads a customer | read | |
| `ZENDESK_REPLY_ZENDESK_TICKET` | Adds a comment: public to the customer, or a private note | write, outside if public | |
| `ZENDESK_UPDATE_ZENDESK_TICKET` | Changes status, priority, assignee or tags | write | |
| `ZENDESK_UPDATE_TICKETS_TAGS` | Adds or removes tags | write | |
| `ZENDESK_CREATE_ZENDESK_TICKET` | Opens a ticket | write | **the ticket's id** |

**Public or private is one setting on the comment.** A private note stays
with agents; a public comment emails the customer. Start with private notes,
and offer public replies only when the owner asked, bounded by the Card.

**What to keep:** each ticket's id the taf handled, or a tag it adds
(`acknowledged`), so it never answers the same ticket twice.

## Don't offer

- Deleting tickets, users or organisations; merging them.
- Redacting comments (`ZENDESK_CREATE_COMMENT_REDACTION`): it erases what a
  customer wrote, for good.
- Administration: triggers, automations, macros, ticket fields and forms,
  brands, support addresses, deletion schedules, help-centre articles.
- Bulk ticket creation or imports.

## Access and plans

- **The API is on every Support plan;** what the taf can do follows the
  agent's role.
- **Rate:** from 200 to 700 calls a minute for the account by plan, shared
  with its other apps. Search answers at most 1,000 results.
- **Search lags** a minute or two behind new tickets; listing by time is
  more exact for a frequent job.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | The agent's role doesn't allow it | "Your Zendesk role doesn't allow this. An admin can change it." |
| Not found (404) | A deleted or merged ticket | "That ticket isn't there any more; it may have been merged." |
| Refused (422) | A field value Zendesk doesn't accept, or a closed ticket | "Zendesk didn't accept that change; a closed ticket can't be updated." |
| Too many calls (429) | Over the account's rate | "Zendesk asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, or expired | "Zendesk needs you to sign in again." |

After any refusal the taf stops that step and says so.
