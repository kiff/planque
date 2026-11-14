# ClickUp

Checked: 2026-11-14 · Sources: [the generated part](generated/clickup.md)
(163 tools), and ClickUp's developer docs at developer.clickup.com: the API
reference, authentication, rate limits and common errors pages.

## How it connects

- **A sign-in to ClickUp**, held by the connection service; the taf never
  sees a password or token. At sign-in the owner chooses which of their
  workspaces to allow; the taf sees only those, and only what the owner can
  see in them.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.
- **ClickUp's API calls a workspace a "team".** A tool's `team_id` or
  `team_Id` is the workspace id; `CLICKUP_GET_AUTHORIZED_TEAMS_WORKSPACES`
  gives it.

## Tools a job usually needs

Read means it changes nothing in the workspace. Returns are the ids a taf
keeps; they come from ClickUp's API reference, since the catalog lists none.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `CLICKUP_GET_AUTHORIZED_TEAMS_WORKSPACES` | Lists the workspaces the sign-in allows | read | workspace ids, names |
| `CLICKUP_GET_SPACES` / `CLICKUP_GET_FOLDERS` / `CLICKUP_GET_LISTS` / `CLICKUP_GET_FOLDERLESS_LISTS` | Walks the workspace down to its lists | read | space, folder and list ids |
| `CLICKUP_GET_TASKS` | Lists a list's tasks, 100 a page, by status, assignee, dates and more | read | task ids, links |
| `CLICKUP_GET_FILTERED_TEAM_TASKS` | Finds tasks across the whole workspace | read | task ids, links |
| `CLICKUP_GET_TASK` | Reads one task | read | the task, its list, its link |
| `CLICKUP_GET_TASK_COMMENTS` | Reads a task's comments | read | comment ids |
| `CLICKUP_SEARCH_DOCS` / `CLICKUP_GET_DOC_PAGE_CONTENT` | Finds Docs, reads one page | read | doc ids, page ids |
| `CLICKUP_CREATE_TASK` | Adds a task to a list | write, notifies | the task id, its link |
| `CLICKUP_UPDATE_TASK` | Changes a task's name, status, dates, assignees or description | write, notifies | the task id |
| `CLICKUP_SET_CUSTOM_FIELD_VALUE` | Sets one custom field on a task | write | not documented |
| `CLICKUP_CREATE_TASK_COMMENT` | Adds a comment to a task | write, notifies | the comment id |
| `CLICKUP_CREATE_CHECKLIST_ITEM` | Adds an item to a task's checklist | write | the checklist id, its items' ids |

**What to keep:** the workspace id and the list id the job works in, and the
id of every task the taf made, to update that same task next time. Never
search for a task the taf made: use the id it kept.

**Pages:** `CLICKUP_GET_TASKS` leaves out closed tasks and subtasks unless
asked, and pages start at 0. Ask for the next page while a page comes
back with 100 tasks.

**Start read-only.** "Tell me what's overdue" needs only reads. Offer a
writing tool only when the owner asked for one. Assignees and watchers are
notified of new tasks, so every task or comment the taf writes for others
is shown to the owner first.

## Don't offer

- Deleting tools (every `CLICKUP_DELETE_...` tool, and
  `CLICKUP_REMOVE_TASK_FROM_LIST`).
- People and permissions: inviting, editing or removing users and guests
  (`CLICKUP_INVITE_USER_TO_WORKSPACE`, `CLICKUP_INVITE_GUEST_TO_WORKSPACE`,
  `CLICKUP_UPDATE_USER_ON_WORKSPACE`, `CLICKUP_REMOVE_USER_FROM_WORKSPACE`,
  and the `CLICKUP_ADD_GUEST_*` and `CLICKUP_REMOVE_GUEST_*` tools), user groups
  (`CLICKUP_CREATE_TEAM`, `CLICKUP_UPDATE_TEAM`), custom roles, and
  `CLICKUP_UPDATE_WORKSPACE_ACL`.
- Workspace setup: creating or changing spaces, folders, views, tags and
  webhooks (`CLICKUP_CREATE_SPACE`, `CLICKUP_UPDATE_SPACE`,
  `CLICKUP_CREATE_WEBHOOK`, the `CREATE_..._VIEW` tools and
  `CLICKUP_UPDATE_VIEW`).
- Billing and seats (`CLICKUP_GET_WORKSPACE_PLAN`,
  `CLICKUP_GET_WORKSPACE_SEATS`): not a taf's business.
- `CLICKUP_GET_ACCESS_TOKEN`: it's part of signing in, not a job tool.
- Chat tools (creating channels, posting messages): they speak to others
  in the owner's name.

## Access and plans

- **The taf can do only what the owner can,** in the workspaces chosen at
  sign-in. A workspace not chosen answers as not authorised; the owner can
  sign in again to add it.
- **Rate, per sign-in, set by the plan of the workspace:** 100 calls a minute
  on Free Forever, Unlimited and Business; 1,000 on Business Plus; 10,000 on
  Enterprise. A job reading many lists should page and pause.
- **The sign-in doesn't expire on its own** today, but ClickUp says that may
  change. It stops when the owner removes the app.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Workspace not authorised (OAUTH_023, OAUTH_026, OAUTH_027, OAUTH_029 to 045) | That workspace wasn't chosen at sign-in | "This ClickUp workspace isn't shared with the connection. Want to sign in again and include it?" |
| Token not found (OAUTH_019, 021, 025, 077) | The owner removed the app's access in ClickUp | "ClickUp needs you to sign in again. Nothing was changed." |
| Too many calls (429) | Over the plan's calls per minute | "ClickUp asked us to slow down. The taf will continue in a minute." |
| Any other refusal on a write | An input ClickUp didn't accept, or an item deleted or made private | Check the item with a read and try once; else "ClickUp didn't accept that change. Nothing was written." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
