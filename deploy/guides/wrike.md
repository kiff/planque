# Wrike

Checked: 2026-11-13 · Sources: the connection service's Wrike tool
definitions as planque reads them (144 tools, [generated](generated/wrike.md)),
and Wrike's API reference.

## How it connects

- **A sign-in to Wrike**, held by the connection service. The taf acts as the
  owner, with their rights in each space. The catalog lists the tools and,
  with a known sign-in, puts the ones its access doesn't cover under
  `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `WRIKE_GET_FOLDERS` | Lists folders and projects | read | **the folder's id** |
| `WRIKE_GET_SPACE` | Reads a space | read | |
| `WRIKE_GET_FOLDERS_TASKS` | Lists the tasks in a folder or project | read | task ids |
| `WRIKE_FETCH_ALL_TASKS` | Lists tasks across the account, by status, date or person | read | task ids |
| `WRIKE_GET_TASK_BY_ID` | Reads one task | read | |
| `WRIKE_CREATE_TASK` | Creates a task in a folder | write | **the task's id** |
| `WRIKE_MODIFY_TASK` | Changes a task: status, dates, assignees | write | |
| `WRIKE_CREATE_TASK_COMMENT` | Comments on a task | write, visible to others | the comment's id |

**What to keep:** the folder's id the taf works in, and each task it creates.
Wrike's ids are long codes, not the numbers in its links; read them once
from the API.

**Assigning a task notifies the person.** The Card says who.

## Don't offer

- Deleting tasks, folders, spaces, attachments or comments.
- Account administration: groups, invitations, custom fields, workflows,
  timelog lock periods, webhooks.
- Approvals and time logs, unless the job is about them and the owner sees
  each one.

## Access and plans

- **The API is on every plan,** but what a taf sees follows the owner's
  licence and their access to each space.
- **Rate:** about 400 calls a minute per user; a list answers up to 1,000
  items.
- **Custom fields and workflows** differ by account: read them before writing
  a status.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | The owner can't edit that space or folder | "You don't have edit rights in that Wrike space, so I can't change it." |
| Not found (404) | A deleted task, or an id from a link rather than the API | "That task isn't there any more." |
| Refused (400) | A status or field that doesn't exist in this workflow | Read the workflow and try once; else "Wrike didn't accept that change." |
| Too many calls (429) | Over Wrike's rate | "Wrike asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, or expired | "Wrike needs you to sign in again." |

After any refusal the taf stops that step and says so.
