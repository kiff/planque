# Zendesk

Checked: 2026-11-14 · Sources: the connection service's Zendesk tool
definitions as planque reads them (480 tools, [generated](generated/zendesk.md)),
and Zendesk's Support API reference.

## How it connects

- **A sign-in to the owner's Zendesk** (`yourcompany.zendesk.com`), held by
  the connection service. The taf acts as the owner's account. The catalog lists the tools and puts the ones its access doesn't
  cover under `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**
- **In the gallery:** "Acknowledge every ticket" uses Zendesk.

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ZENDESK_LIST_ZENDESK_TICKETS` | Lists tickets, newest first | read | ticket ids |
| `ZENDESK_SEARCH_ZENDESK` | Finds tickets, users or organisations by fields (`status:new`) | read | ticket ids |
| `ZENDESK_GET_ZENDESK_TICKET_BY_ID` | Reads one ticket | read | the ticket, with its id and the requester's id |
| `ZENDESK_GET_TICKET_COMMENTS` | Reads a ticket's conversation | read | |
| `ZENDESK_GET_USER` | Reads a customer | read | |
| `ZENDESK_REPLY_ZENDESK_TICKET` | Adds a comment: public to the customer, or a private note | write, outside if public | |
| `ZENDESK_UPDATE_ZENDESK_TICKET` | Changes status, priority, assignee or tags | write | |
| `ZENDESK_UPDATE_TICKETS_TAGS` | Adds or removes tags | write | |
| `ZENDESK_CREATE_ZENDESK_TICKET` | Opens a ticket | write | **the ticket's id** |

**Public or private is one setting on the comment:** a public comment, or
an internal note. Start with internal notes, and offer public replies only
when the owner asked, bounded by the Card.

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

- **Rate:** a limit for the whole account, by plan: 200 calls a minute on
  Team, 400 on Growth and Professional, 700 on Enterprise, and up to 2,500
  with Zendesk's higher tiers. Each app is also capped per user. One ticket
  can be updated at most 30 times in 10 minutes by one user. Search answers
  at most 1,000 results, 100 a page.
- **Search lags** a minute or two behind new tickets; listing by time is
  more exact for a frequent job.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | This account isn't allowed to do that in Zendesk | "Zendesk doesn't let your account do this step. A Zendesk admin can check its access." |
| Not found (404) | The ticket isn't there, or the id is wrong | "That ticket isn't in Zendesk any more. The taf will skip it." |
| Refused (422) | A value Zendesk doesn't accept | "Zendesk didn't accept that change. The taf has stopped and will show you what it tried." |
| Too many calls (429) | Over the account's rate | "Zendesk asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | The sign-in isn't accepted any more | "Zendesk needs you to sign in again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
