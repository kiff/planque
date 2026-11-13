# Sentry

Checked: 2026-11-13 · Sources: the connection service's Sentry tool
definitions as planque reads them ([generated part](generated/sentry.md), 210
tools), and Sentry's API reference: authentication, permissions and scopes,
rate limits, pagination, the issue endpoints, and the data storage location
page.

## How it connects

- **A sign-in to Sentry**, held by the connection service; the taf never sees
  a password or token. With a known sign-in, the catalog puts the tools its
  access doesn't cover under `unusable`, with why. A held sign-in is reused
  for a new taf when it covers the tools asked.
- **One organisation per sign-in.** The owner picks one Sentry organisation
  while signing in, and the sign-in reaches only that one. A second
  organisation needs a second sign-in.
- **One approval for the sign-in, one for the tools and the Card.** Ask for
  every tool the job needs in the first request.

## Tools a job usually needs

Read changes nothing in Sentry. Returns are the ids Sentry's reference
documents for each call.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SENTRY_LIST_ORGANIZATIONS` | Lists the organisation the sign-in reaches | read | organisation ids and slugs |
| `SENTRY_RETRIEVE_ORGANIZATION_PROJECTS` | Lists the organisation's projects | read | project ids and slugs, a page cursor |
| `SENTRY_LIST_AN_ORGANIZATIONS_ISSUES` | Lists issues, filtered by project, environment, period or search; up to 100 a page | read | issue ids, short ids, links, a page cursor |
| `SENTRY_RETRIEVE_PROJECT_ISSUES_LIST` | Lists one project's issues | read | issue ids, short ids, a page cursor |
| `SENTRY_GET_ORGANIZATION_ISSUE_DETAILS` | Reads one issue: counts, first and last seen, status, assignee | read | the issue id, short id, link |
| `SENTRY_RETRIEVE_ISSUE_EVENTS_BY_ID` | Lists the events (occurrences) of one issue | read | event ids |
| `SENTRY_FETCH_ISSUE_EVENT_BY_ID` | Reads one event of an issue, with its stack trace and context | read | the event id |
| `SENTRY_LIST_A_TAGS_VALUES_FOR_AN_ISSUE` | Counts one tag's values for an issue (browser, release, user) | read | tag values and counts |
| `SENTRY_LIST_ORGANIZATION_RELEASES` | Lists releases | read | release versions |
| `SENTRY_UPDATE_ISSUE_ATTRIBUTES_IN_ORGANIZATION` | Changes one issue: resolve, ignore, reopen, assign, bookmark | write, visible to the team | the updated issue |

**What to keep:** the organisation slug and project slug, and each issue's
id (with its short id and link, to show the owner). A list without a search
shows only unresolved issues: send an empty search to see all.

**Start read-only.** A daily or weekly digest of new and rising issues is
read tools only. Offer `SENTRY_UPDATE_ISSUE_ATTRIBUTES_IN_ORGANIZATION` only
when the owner asked the taf to triage, and show each change first. Never
set an issue public: that opens it to anyone with the link.

## Don't offer

- Deleting anything (`SENTRY_DELETE_*`, `SENTRY_BULK_REMOVE_AN_ORGANIZATIONS_ISSUES`):
  Sentry's events can only go by deleting the whole issue, and they don't
  come back.
- Bulk issue changes (`SENTRY_BULK_MUTATE_AN_ORGANIZATIONS_ISSUES`,
  `SENTRY_UPDATE_PROJECT_ISSUE_STATUS_AND_DETAILS`): one call can merge,
  discard or resolve many issues at once.
- People and permissions: inviting or removing members, roles, teams, and
  the SCIM and user email tools.
- Settings: organisation and project settings, client keys, spike
  protection, inbound filters, symbol sources, ownership rules,
  integrations, webhooks, and creating or changing alerts and monitors
  (they decide who gets paged).
- `SENTRY_INGEST_EVENT_VIA_DSN`: it sends raw events into a project and
  would pollute the owner's error data.

## Access and plans

- **Reading issues and events** needs read access to them; changing an issue
  needs write access to issues. Deleting needs admin access, which a taf
  never asks for. Organisation, project, team and member tools each need
  their own access.
- **The sign-in's access lasts 30 days** and is renewed with a refresh
  token. If renewal fails, the owner signs in again.
- **Rate:** Sentry limits each caller per endpoint, both calls per second
  and calls running at once. It doesn't publish one number for all
  endpoints. Polling often hits the limit quickly: read once per run, page
  with the cursor, and don't run calls in parallel.
- **Data location:** an organisation keeps its error data in the US or in
  the EU (Frankfurt), chosen when it was created; it can't be moved.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | Access was removed, or it expired and couldn't be renewed | "Sentry needs you to sign in again. Nothing was changed." |
| Refused for missing access (403) | The sign-in didn't give this tool's access | "Sentry is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Not found (404) | The issue was deleted, or it's in another organisation than the one chosen at sign-in | "I can't find that in the Sentry organisation you connected. If it's in another organisation, connect that one too." |
| Refused for too many calls | Over Sentry's limit for this endpoint | "Sentry asked us to slow down. The taf will try again in a few minutes." |
| Refused as invalid (400) | A bad search, period or status value | Fix it and try once; else "Sentry didn't accept that request. Nothing was changed." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
