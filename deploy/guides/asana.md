# Asana

Checked: 2026-11-14 · Sources: the connection service's published Asana tool
list (about 250 tools), Asana's API reference (rate limits, errors,
pagination, access scopes, task search). Not yet run on a phone. The tool
names below are checked against the generated part once the keyed run lands.

## How it connects

- **A sign-in to Asana**, held by the connection service; the taf never sees
  a password or token. With a known sign-in, the catalog puts the tools its
  access doesn't cover under `unusable`, with why.
- **The sign-in acts as the owner.** The taf sees only the workspaces,
  teams and projects the owner can see, and can't do more than the owner
  can.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.
- **Everything lives in a workspace.** Most jobs start by reading the
  owner's workspaces and keeping the right one's id.

## Tools a job usually needs

Read means it changes nothing in Asana. Returns are what a taf keeps to use
next time; the catalog's `returns` gives their exact paths.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ASANA_GET_CURRENT_USER` | Says who the sign-in is, and its workspaces | read | the user's id, workspace ids |
| `ASANA_GET_MULTIPLE_WORKSPACES` | Lists the workspaces the owner belongs to | read | workspace ids and names |
| `ASANA_GET_WORKSPACE_PROJECTS` | Lists a workspace's projects | read | project ids and names |
| `ASANA_GET_SECTIONS_IN_PROJECT` | Lists a project's sections (its columns or groups) | read | section ids and names |
| `ASANA_GET_TASKS_FROM_A_PROJECT` | Lists a project's tasks, a page at a time | read | task ids, a page offset |
| `ASANA_GET_TASKS_FROM_SECTION` | Lists the tasks in one section | read | task ids |
| `ASANA_GET_A_TASK` | Reads one task: name, notes, assignee, due date, fields | read | the task, its link |
| `ASANA_GET_STORIES_FOR_TASK` | Reads a task's comments and history | read | story ids, text |
| `ASANA_SEARCH_TASKS_IN_WORKSPACE` | Finds tasks by text, assignee, date or field (paid plans only) | read | task ids |
| `ASANA_CREATE_A_TASK` | Creates a task in a project or workspace | write | the task's id, its link |
| `ASANA_CREATE_SUBTASK` | Adds a subtask under a task | write | the subtask's id |
| `ASANA_UPDATE_A_TASK` | Changes a task: name, notes, due date, assignee, done | write | the task's id |
| `ASANA_ADD_TASK_TO_SECTION` | Moves a task into a section | write | nothing to keep |
| `ASANA_CREATE_TASK_COMMENT` | Adds a comment to a task | write | the comment's id |

**What to keep:** the workspace id and project id the job works in, and the
id of every task the taf creates. Next run, read or update that task by its
id. Never search by name for something the taf made: names repeat, and the
search is slow to see new tasks.

**Start read-only.** "Tell me what's due this week" is reads only. Offer
creating or updating tasks only when the owner asked for it.

## Don't offer

- **Deleting tools** (tasks, projects, sections, tags, attachments,
  comments): a taf archives or marks done instead.
- **People and access tools** (adding or removing users from a workspace,
  team or project; approving access requests; memberships): they change who
  can see the owner's work.
- **Admin and setup tools** (custom fields, project templates, webhooks,
  rules, roles, rates, budgets, time-tracking settings): they change how the
  workspace works for everyone in it.
- **Parallel requests** (`ASANA_SUBMIT_PARALLEL_REQUESTS`): it bundles many
  calls into one and hides what each does. Name the single tools instead.
- **Event triggers** (task created, moved or updated): they aren't tools a
  taf calls.

## Access and plans

- **Reading needs read access** to tasks, projects and workspaces; creating
  or changing tasks needs write access too. A sign-in made for reading
  doesn't let the taf write.
- **Search needs a paid plan.** Finding tasks by text or field works only
  when the workspace, or one of the owner's teams in it, is on a paid plan.
  On a free plan, list a project's tasks instead and filter them.
- **Other paid features** answer "payment required" on a free workspace.
  Offer a job built on reading and writing tasks first; it works on every
  plan.
- **Rate:** about 150 calls a minute on a free workspace and about 1,500 on
  a paid one, per sign-in. Search is limited to about 60 a minute. At most
  about 50 reads and 15 writes at the same moment.
- **Duplicating** a project or task, and making one from a template, are
  limited to 5 at once per user, and these count the owner's own ones in
  Asana too.
- **Big answers:** lists answer at most 100 items a page. A very large list
  (tens of thousands) can be cut short; read it project by project, or
  section by section.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in didn't give this tool's access: reading only, say, and the tool creates a task | "Asana is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Refused, forbidden on one task or project | The owner can't see or change it (private project, guest) | "I can't reach that task or project with your Asana account. It may be private, or need someone to add you." |
| "Payment required" | Search, or another paid feature, on a free workspace | "That needs a paid Asana plan. I can do it another way, by reading the project's tasks, if that works for you." |
| Too many calls | Over the per-minute limit, or 5 duplications already running | "Asana asked us to slow down. The taf will try again in a minute." |
| Not found | The task or project was deleted, or the id is wrong | "That task isn't in Asana any more. The taf will skip it." |
| Bad input | A required value is missing, a workspace for example | Fix it and try once; else "Asana didn't accept that change. Nothing was written." |
| The list was cut short | The answer was too large | Read it again by section; tell the owner only if it still fails. |
| The sign-in no longer works | Access removed, or the workspace blocked this app | "Asana needs you to sign in again. If your workspace admin blocks outside apps, ask them to allow it." |
| Asana's own error | A problem on Asana's side | "Asana had a problem. The taf will try again later. Nothing was changed." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create.
