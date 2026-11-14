# Sentry

Checked: 2026-11-14 · Sources: the connection service's Sentry tool
definitions as planque reads them (the generated part beside this file),
and Sentry's API documentation (rate limits, issues, events, releases).

## How it connects

- **A sign-in to Sentry**, held by the connection service; the taf never
  sees a password or token. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`,
  with why.
- **Almost every tool needs the organization's slug** (the short name in
  its Sentry address). Read it first with `SENTRY_LIST_ORGANIZATIONS`, then
  keep it.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.
- **Organizations live in a region** (Sentry names a US and a DE host). Sentry
  says its regional address answers faster; it doesn't say other addresses
  fail.

## Tools a job usually needs

Read means it changes nothing in Sentry. The catalog gives no `returns`
paths for Sentry today; the ids below are those Sentry's documented
answers show. "Not documented" means the reference shows no answer.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SENTRY_LIST_ORGANIZATIONS` | Lists the organizations the sign-in can see | read | not documented |
| `SENTRY_GET_PROJECT_LIST` | Lists the projects the sign-in can see | read | not documented |
| `SENTRY_LIST_AN_ORGANIZATIONS_ISSUES` | Lists issues (grouped errors) across projects, filtered by project, environment, time and search; unresolved by default | read | issue ids, short ids, links, statuses, a paging cursor |
| `SENTRY_GET_ORGANIZATION_ISSUE_DETAILS` | Reads one issue: status, assignee, first and last seen, event count | read | the issue's id, short id, link |
| `SENTRY_RETRIEVE_ISSUE_EVENTS_BY_ID` | Lists the events (single occurrences) behind one issue | read | event ids, titles, dates, tags |
| `SENTRY_FETCH_ISSUE_EVENT_BY_ID` | Reads one event of an issue in full | read | not documented |
| `SENTRY_LIST_ORGANIZATION_RELEASES` | Lists releases, with how many new issues each brought | read | release versions, dates, projects |
| `SENTRY_GET_ORGANIZATION_RELEASE_VERSION` | Reads one release | read | not documented |
| `SENTRY_UPDATE_ISSUE_ATTRIBUTES_IN_ORGANIZATION` | Changes one issue: resolve, ignore, reopen, assign, bookmark | write | not documented |

**What to keep:** the organization's slug; each issue's id, to read or
update that same issue later; its short id and link, to show the owner. A
paging cursor is used once, then dropped. Never look an issue up by its
title: use the id the taf kept.

**Start read-only.** "Tell me each morning what's new and broken" needs only
reads. Resolving, ignoring or assigning an issue changes what the team sees
and is told about, so offer
`SENTRY_UPDATE_ISSUE_ATTRIBUTES_IN_ORGANIZATION` only when the owner asked,
and each change waits for their yes. Sentry changes only the fields sent.

## Don't offer

- Deleting anything: issues (`SENTRY_DELETE_ORGANIZATION_ISSUE`,
  `SENTRY_DELETE_PROJECT_ISSUES`, `SENTRY_BULK_REMOVE_AN_ORGANIZATIONS_ISSUES`),
  projects, releases, alerts, monitors or replays. A deleted issue's history
  is gone.
- Bulk changes (`SENTRY_BULK_MUTATE_AN_ORGANIZATIONS_ISSUES`,
  `SENTRY_UPDATE_PROJECT_ISSUE_STATUS_AND_DETAILS`, which changes many issues
  of a project at once and can merge them). Change one issue at a time.
- Project and organization admin: members, teams, roles, keys, alert rules,
  webhooks, integrations, spike protection, ownership settings.
- Creating releases or deploys (`SENTRY_CREATE_RELEASE_FOR_ORGANIZATION`,
  `SENTRY_CREATE_RELEASE_DEPLOY_FOR_ORG`): that's the team's build pipeline's
  job.
- Sending events in (`SENTRY_INGEST_EVENT_VIA_DSN`): it would put made-up
  errors into the owner's data.
- Starting an automated fix (`SENTRY_START_SEER_ISSUE_FIX`): it acts on
  the team's code outside the owner's view.
- `SENTRY_RETRIEVE_CUSTOM_INTEG_ISSUE_LINKS_GIVEN_SENTRY_ISSUE` is generated
  as a write though it reads; leave it out.

## Access and plans

- **Reading issues and events needs event read access; changing an issue
  needs event write access.** Reading releases needs project or release
  access. A sign-in made for reading doesn't let the taf change issues.
- **Rate:** Sentry limits each caller per tool, per second, and also caps
  how many calls run at once. It publishes no fixed numbers: they differ by
  tool. Each answer says how many calls are left and when the count resets.
  Polling often trips the limit; a taf reads on its schedule and pages with
  the cursor.
- **Pages:** issue lists answer at most 100 a page.
- **Plans** aren't covered here: we checked no plan limit on these tools.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in didn't give this tool's access, or the person's role in the organization doesn't allow it | "Sentry is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Too many calls | Over the per-tool rate or the limit on calls at once | "Sentry asked us to slow down. The taf will try again shortly." |
| The sign-in no longer works (401) | The owner removed access, or the sign-in expired | "Sentry needs you to sign in again. Nothing was changed." |
| Not found (404) | The issue was deleted or merged, or belongs to a project the sign-in can't see | "That issue isn't in Sentry any more, or this connection can't see it. The taf will skip it." |
| Rejected (400) | A value Sentry didn't accept, such as an unknown assignee | "Sentry didn't accept that change. The taf has stopped and will show you what it tried." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't read an issue never resolves or assigns
it.
