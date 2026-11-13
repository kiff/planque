# Linear

Checked: 2026-11-13 · Sources: the connection service's Linear tool
definitions as planque reads them (47 tools), and Linear's developer
documentation (rate limits, sign-in, the API's answers, paging and
filtering).

## How it connects

- **A sign-in to Linear**, made on Linear's own page and held by the
  connection service; the taf never sees a token. planque's catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`, with why.
- **Linear's access comes in a few sizes:** reading, a narrower right to
  create and update issues, a narrower right to comment, full writing, and
  admin. Reading is always included. A sign-in made to create issues
  doesn't let the taf archive or delete them.
- **By default the taf acts as the person who signed in.** What it creates
  shows under their name in Linear.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in Linear. The catalog gives no `returns`
paths for Linear today; "not documented" means the pages read didn't show
the answer's fields.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `LINEAR_GET_CURRENT_USER` | Reads the person the sign-in acts as | read | not documented |
| `LINEAR_LIST_LINEAR_TEAMS` | Lists the workspace's teams | read | not documented |
| `LINEAR_LIST_LINEAR_STATES` | Lists one team's workflow states (to do, in progress, done…) | read | not documented |
| `LINEAR_LIST_LINEAR_LABELS` | Lists labels, for one team or all | read | not documented |
| `LINEAR_LIST_LINEAR_USERS` | Lists the workspace's people | read | not documented |
| `LINEAR_LIST_ISSUES_BY_TEAM_ID` | Lists one team's issues, a page at a time | read | issues, a page cursor |
| `LINEAR_SEARCH_ISSUES` | Finds issues by words | read | issues, a page cursor |
| `LINEAR_GET_LINEAR_ISSUE` | Reads one issue | read | not documented |
| `LINEAR_LIST_COMMENTS` | Lists comments, a page at a time | read | comments, a page cursor |
| `LINEAR_CREATE_LINEAR_ISSUE` | Opens an issue in a team | write | the new issue's id and title |
| `LINEAR_UPDATE_ISSUE` | Changes an issue's state, assignee, labels, due date or text | write | not documented |
| `LINEAR_CREATE_LINEAR_COMMENT` | Comments on an issue | write | not documented |

**What to keep:** the team's id and its state ids, since creating and
moving issues needs them; each issue's id, to update or comment on that
same issue later. Linear also accepts an issue's short key (like `ENG-123`)
where it takes an issue id. Never search for an issue the taf made: use the
id it kept.

**Paging:** lists come 50 at a time unless asked otherwise; the taf passes
the cursor from one page to get the next, while Linear says there is more.

**Start read-only.** "Tell me what changed this week" needs only the read
tools. Offer creating, updating or commenting only when the owner asked for
it; each shows under the owner's name.

## Don't offer

- Deleting and archiving (`LINEAR_DELETE_LINEAR_ISSUE`,
  `LINEAR_ARCHIVE_ISSUE`, `LINEAR_ARCHIVE_PROJECT`).
- Team setup (`LINEAR_CREATE_TEAM`, `LINEAR_UPDATE_TEAM`): it changes how
  the whole workspace is organised.
- The free-form query tool (`LINEAR_RUN_QUERY_OR_MUTATION`): it can do
  anything the sign-in allows, so a Card can't say what it will do. Use the
  named tools.
- Credential tools (`LINEAR_VALIDATE_CREDENTIAL`): nothing a job needs.

## Access and plans

- **Rate:** per person, 5,000 calls an hour through a sign-in, plus a
  budget of 2,000,000 "complexity" points an hour, where a bigger query
  costs more points. One query over 10,000 points is always refused.
  Business workspaces get twice these limits and Enterprise five times;
  Free and Basic get the standard ones.
- **A sign-in's token is short-lived** (a day) and is renewed for up to 30
  days at a time. If renewal stops (access revoked, or the app removed),
  the owner signs in again.
- **Partial answers:** a query can come back partly answered, with the
  errors listed beside the data. The taf treats an answer with errors as
  failed for the part that's missing.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in was made for reading, or only for issues or comments, and this tool needs more | "Linear is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Too many calls (429, `RATELIMITED`) | Over the hourly calls or complexity budget | "Linear asked us to slow down. The taf will try again later." |
| A query too big | Over Linear's limit for a single query | "That request to Linear was too large. The taf will ask for less at a time." |
| The sign-in no longer works | Access was revoked, or the sign-in couldn't be renewed | "Linear needs you to sign in again. Nothing was changed." |
| Answer with errors | Part of the request failed, such as an issue that no longer exists | "Linear couldn't do part of that. The taf has stopped and will show you what was missing." |
| Linear unavailable (5xx) | Trouble on Linear's side | "Linear isn't answering right now. The taf will try once more later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find its issue never opens a new one in
its place.
