# Dynamics 365

Checked: 2026-11-13 · Sources: the connection service's Dynamics 365 tool
definitions as planque reads them (119 tools), Microsoft's Web API
reference for Dynamics 365 data (requests and errors, paging, creating
rows), its service protection limits page and its daily request limits page.

## How it connects

- **A sign-in with the owner's Microsoft work account**, held by the
  connection service; the taf never sees a password or token. planque's
  catalog lists the tools and, with a known sign-in, puts the ones its
  access doesn't cover under `unusable`, with why.
- **One organisation's environment.** Each company's Dynamics 365 lives at
  its own address; the taf works only in the environment signed in to.
- **The taf acts as the owner.** It sees and changes only what the owner's
  own Dynamics 365 role allows.

## Tools a job usually needs

Read means it changes nothing. Returns are named from Microsoft's API
reference, since the catalog gives none for this app. Every record id is
a long unique id (a GUID).

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `DYNAMICS365_WHO_AM_I` | Says who the sign-in is | read | the user's id, business unit id, organisation id |
| `DYNAMICS365_QUERY_ACCOUNTS` / `DYNAMICS365_QUERY_CONTACTS` | Finds companies or people by a filter | read | account ids or contact ids, a next-page link |
| `DYNAMICS365_QUERY_OPPORTUNITIES` | Finds open or closed deals | read | opportunity ids, a next-page link |
| `DYNAMICS365_DYNAMICSCRM_GET_ALL_LEADS` | Lists leads by a filter | read | lead ids |
| `DYNAMICS365_QUERY_CASES` / `DYNAMICS365_QUERY_TASKS` | Finds support cases or tasks | read | case ids or activity ids, a next-page link |
| `DYNAMICS365_GET_ACCOUNT` / `DYNAMICS365_GET_CONTACT` / `DYNAMICS365_GET_OPPORTUNITY` | Reads one record | read | the record and its id |
| `DYNAMICS365_DYNAMICSCRM_GET_A_LEAD` | Reads one lead | read | the lead and its id |
| `DYNAMICS365_DYNAMICSCRM_CREATE_LEAD` | Adds a lead | write | the new lead's id |
| `DYNAMICS365_DYNAMICSCRM_CREATE_CONTACT` | Adds a person | write | the new contact's id |
| `DYNAMICS365_CREATE_TASK` / `DYNAMICS365_CREATE_PHONE_CALL` | Logs a to-do or a call | write | the new activity's id |
| `DYNAMICS365_CREATE_ANNOTATION` | Adds a note, with an optional file, to a record | write | the new note's id |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_LEAD` | Changes a lead's fields | write | not documented (nothing on success) |
| `DYNAMICS365_UPDATE_CONTACT` | Changes a contact's fields | write | not documented (nothing on success) |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_OPPORTUNITY` | Changes a deal's stage, value or close date | write | not documented (nothing on success) |

**What to keep:** the id of every record the taf creates or works on, with
its type (lead, contact, account); and a query's next-page link to carry on
where it stopped. Never search for something the taf made: use the id it
kept.

**Paging:** pass the next-page link (`next_link`) as it came; Microsoft's
API ignores skip-style paging, so the `skip` input of the leads list
doesn't page. Filter and order by a unique field instead.

**Start read-only:** query and get tools. Offer creating leads, contacts, tasks and notes when the
job adds records, and updates only when the owner names what may change.
Each update is shown to the owner first.

## Don't offer

- Every deleting tool (`DYNAMICS365_DELETE_ACCOUNT`,
  `DYNAMICS365_DELETE_LEAD` and the rest), and
  `DYNAMICS365_CLEAR_ACCOUNT_DESCRIPTION`.
- People and permissions: `DYNAMICS365_CREATE_SYSTEM_USER`,
  `DYNAMICS365_UPDATE_SYSTEM_USER`, `DYNAMICS365_ASSIGN_USER_SECURITY_ROLE`.
- Billing and catalogue: invoices, quotes, sales orders and their line
  items, price lists, products and currencies (create and update tools).
- Bulk and raw tools: `DYNAMICS365_EXECUTE_BATCH`, and
  `DYNAMICS365_QUERY_ENTITY_DATA`, which reads any table by name, users
  and settings included.

## Access and plans

- **The owner needs a paid Dynamics 365 licence and a role** that allows
  each action. A role that lets them read accounts may not let them edit
  or create them; the taf can't do more than the owner could.
- **An organisation's Microsoft admin can turn off user consent** to
  outside apps. Then the sign-in waits for an admin to approve it.
- **Daily allowance:** Microsoft counts requests per user per 24 hours:
  40,000 for a Dynamics 365 Enterprise or Professional licence, 6,000 for a
  Team Member licence.
- **Short-term limits, per user:** by default about 6,000 requests in any
  five minutes, 20 minutes of combined processing time in five minutes,
  and 52 requests at once, on each of the environment's servers.
- **Pages:** a query answers up to 5,000 rows a page, with a link to the
  next.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | Expired or revoked | "Dynamics 365 needs you to sign in again. Nothing was changed." |
| Refused (403), privilege | The owner's role doesn't allow this action on this record type | "Your Dynamics 365 role doesn't allow this step. Your administrator can extend it, or the taf does without this step." |
| Sign-in needs approval | The organisation requires an admin to approve outside apps | "Your Microsoft admin needs to approve this connection before it can work." |
| Not found (404) | The record was deleted, or isn't visible to the owner | "That record isn't there any more. The taf will skip it." |
| Precondition failed (412) | The record changed since it was read, or would be a duplicate | Read it again and retry once; else "That record changed or already exists. Nothing was changed." |
| Bad request (400) | A field name or value isn't valid for this organisation | "Dynamics 365 didn't accept that change. Nothing was written." |
| Too many requests (429) | Over the short-term limit; Microsoft says how long to wait | "Dynamics 365 asked us to slow down. The taf will continue in a few minutes." |
| Service unavailable (503) | Microsoft's side is down for a while | "Dynamics 365 isn't answering right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
