# Jira

Checked: 2026-11-13 · Sources: the connection service's Jira tool
definitions as planque reads them (166 tools), and Jira's REST API
documentation (its published API description, rate limiting, and sign-in
for outside apps).

## How it connects

- **A sign-in to Jira**, held by the connection service; the taf never sees
  a password or token. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **The taf acts as the person who signed in.** Jira always applies that
  person's own permissions on top of the sign-in's access, so a tool the
  sign-in allows can still be refused in a project the person can't edit.
- **One or more Jira sites.** The person picks which sites the sign-in
  reaches. A held sign-in is reused for a new taf when it covers the tools
  asked: ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in Jira. The catalog gives no `returns` paths
for Jira today; the ids below are those Jira's answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `JIRA_GET_CURRENT_USER` | Reads who the sign-in acts as | read | the person's account id |
| `JIRA_GET_ALL_PROJECTS` | Lists the projects the person can see, a page at a time | read | project ids and keys |
| `JIRA_GET_PROJECT` | Reads one project | read | the project's id and key |
| `JIRA_GET_ISSUE_TYPES` | Lists the kinds of issue (task, bug and so on) | read | issue type ids and names |
| `JIRA_SEARCH_ISSUES` | Finds issues by project, assignee, label, date or a search query | read | issue ids and keys, a page token |
| `JIRA_GET_ISSUE` | Reads one issue and its fields | read | the issue's id and key |
| `JIRA_LIST_ISSUE_COMMENTS` | Reads an issue's comments | read | comment ids |
| `JIRA_GET_TRANSITIONS` | Lists the status moves open to an issue now | read | transition ids and names |
| `JIRA_FIND_USERS_FOR_PICKER` | Finds a person by name, to assign work | read | account ids |
| `JIRA_CREATE_ISSUE` | Opens a new issue in a project | write | the new issue's id and key |
| `JIRA_EDIT_ISSUE` | Changes an issue's summary, description, labels or due date | write | nothing (or the issue, if asked) |
| `JIRA_ADD_COMMENT` | Adds a comment to an issue | write | the comment's id |
| `JIRA_TRANSITION_ISSUE` | Moves an issue to another status | write | nothing |
| `JIRA_ASSIGN_ISSUE` | Assigns an issue to a person | write | nothing |

**What to keep:** each issue's key (such as `OPS-142`), to read, comment on
or move that same issue later; the project key; and the person's account id.
Never search for an issue the taf opened: use the key it kept.

**Look before moving.** A status move uses a transition id from
`JIRA_GET_TRANSITIONS`, read just before; the moves open to an issue depend
on its workflow and its status now.

**Search with `JIRA_SEARCH_ISSUES`.** The other search tool,
`JIRA_SEARCH_FOR_ISSUES_USING_JQL_GET`, isn't marked read-only in the
catalog, so it counts as a write.

**Start read-only.** "Tell me what's overdue in this project" needs only read
tools. Offer writes when the owner asked for changes; comments and
assignments notify other people in Jira.

## Don't offer

- Deleting tools (`JIRA_DELETE_ISSUE`, `JIRA_DELETE_COMMENT`,
  `JIRA_DELETE_ATTACHMENT`, `JIRA_DELETE_WORKLOG` and the rest): they can't
  be undone from the taf.
- Configuration and admin tools: creating issue types, link types, statuses,
  project roles, workflows or custom field contexts. They change how Jira
  works for everyone.
- Bulk tools (`JIRA_BULK_CREATE_ISSUE`, `JIRA_BULK_MOVE_ISSUES`): one
  mistake repeats across many issues.

## Access and plans

- **Permissions come from Jira.** What the person can't do in Jira, the taf
  can't either, whatever the sign-in allows.
- **Rate:** Jira counts calls in points per hour. An app's default
  allowance is 65,000 points an hour, shared by every site the app reaches
  and reset on the hour. A read
  costs 1 point; reading users and permissions costs 2. There's also a
  per-second limit per site and kind of call (about 100 reads a second).
- **One issue at a time:** at most 20 changes in 2 seconds, and 100 in 30
  seconds, on the same issue.
- **Per-issue caps:** an issue holds a limited number of comments,
  attachments and links. Over the cap, Jira refuses the addition.
- **A sign-in stays valid** while it's used at least every 90 days. After 90
  days unused, the owner signs in again.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing permission (403) | The person can't do this in that project, or the sign-in's access doesn't cover it | "Jira didn't let the taf do this in that project. Check your permissions there, or do without this step." |
| Issue not found (404) | The issue was deleted, or the person can't see it | "The taf can't see that issue in Jira. It may be deleted or in a project you can't open." |
| Change rejected (400) | A required field is missing, a field isn't on the screen, or the status move isn't open from here | "Jira didn't accept that change. The taf has stopped and will show you what it tried." |
| Jira's set-up blocks it (422) | A project or workflow setting prevents the change | "Jira's project settings block this change. Your Jira admin can check them." |
| Conflicting change (409) | Someone else changed the issue at the same moment | "Someone was changing the same issue. The taf will try once more." |
| Issue full (413) | The issue has reached its cap on comments, attachments or links | "That issue has reached Jira's limit for comments or attachments. The taf will skip it." |
| Too many calls (429) | Over the hourly allowance, the per-second limit or the per-issue limit; Jira says how long to wait | "Jira asked us to slow down. The taf will try again shortly." |
| The sign-in no longer works (401) | Access was removed, or unused for 90 days | "Jira needs you to sign in again. Nothing was changed." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find its issue asks, rather than opening
another one.
