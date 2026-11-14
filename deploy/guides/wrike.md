# Wrike

Checked: 2026-11-14 · Sources: the connection service's Wrike tool
definitions as planque reads them (144 tools, [generated](generated/wrike.md)),
and Wrike's developer documentation (errors, authorisation, ids, tasks and
comments).

## How it connects

- **A sign-in to Wrike**, held by the connection service. The taf acts as the
  owner, with their rights in each space. The catalog lists the tools and,
  with a known sign-in, puts the ones its access doesn't cover under
  `unusable`. Offer only what's under `tools`.
- **Wrike keeps each account's data in one data centre** (the US or the EU),
  and the sign-in tells the connection which one to use.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `WRIKE_GET_FOLDERS` | Lists folders and projects | read | **folder ids** |
| `WRIKE_GET_SPACE` | Reads a space | read | the space |
| `WRIKE_GET_FOLDERS_TASKS` | Lists the tasks in a folder or project | read | task ids |
| `WRIKE_FETCH_ALL_TASKS` | Lists tasks across the account, by status, date or person | read | task ids |
| `WRIKE_GET_TASK_BY_ID` | Reads one task | read | the task, with its id and link |
| `WRIKE_CREATE_TASK` | Creates a task in a folder | write | **the task's id and link** |
| `WRIKE_MODIFY_TASK` | Changes a task: status, dates, assignees | write | the task |
| `WRIKE_CREATE_TASK_COMMENT` | Comments on a task | write, visible to others | the comment's id |

**What to keep:** the folder's id the taf works in, and each task it creates.
Wrike's ids are text of up to 256 characters with no fixed pattern: read
them from the API and keep them as given.

**Assignees are people:** the Card says who the taf may assign.

**A comment can be changed only in the 5 minutes after it's posted.**

## Don't offer

- Deleting tasks, folders, spaces, attachments or comments.
- Account administration: groups, invitations, custom fields, workflows,
  timelog lock periods, webhooks.
- Approvals and time logs, unless the job is about them and the owner sees
  each one.

## Access and plans

- **What a taf sees follows the owner's access** to each space. Wrike can
  refuse an action because of the account's licence or quota.
- **Rate:** 400 calls a minute per address or sign-in. A task list answers
  up to 1,000 tasks a page, with a token for the next page.
- **Custom fields and workflows** differ by account: read them before writing
  a status.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | The owner can't reach that item, or the account's licence or quota doesn't allow the action | "Wrike doesn't let your account do that here. Check your access to that space, or your Wrike plan." |
| Not found (404) | The task was deleted, or the id is wrong | "That task isn't there any more." |
| Refused (400) | A value or field Wrike doesn't accept, such as a status outside this workflow | Read the workflow and try once; else "Wrike didn't accept that change." |
| Too many calls (429) | Over Wrike's rate | "Wrike asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | The sign-in is invalid or expired | "Wrike needs you to sign in again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
