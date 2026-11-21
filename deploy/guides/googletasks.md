# Google Tasks

Checked: 2026-11-14 · Sources: the connection service's Google Tasks tool
definitions as planque reads them (16 tools, [generated](generated/googletasks.md)),
and Google's Tasks API reference.
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **A sign-in to Google**, held by the connection service. The grant's catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **One approval for the sign-in, one for the tools and the Card.**
- **Proved on a phone (21 Nov 2026):** create a list, add a task with a due
  date, and list a list's tasks, each write after the owner's yes and done
  once. Add to a list only once its id is kept: a task added before the list
  exists has nowhere to go. Marking a task done and deleting were not tried.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLETASKS_LIST_TASK_LISTS` | Lists the owner's task lists | read | **each list's id** |
| `GOOGLETASKS_LIST_TASKS` | Lists the tasks in one list, by due date or status | read | task ids |
| `GOOGLETASKS_LIST_ALL_TASKS` | Lists tasks across every list | read | task and list ids |
| `GOOGLETASKS_GET_TASK` | Reads one task | read | |
| `GOOGLETASKS_INSERT_TASK` | Adds a task, with notes and a due date | write | **the task's id** |
| `GOOGLETASKS_PATCH_TASK` | Changes some fields of a task, such as marking it done | write | |
| `GOOGLETASKS_CREATE_TASK_LIST` | Creates a list | write | **the list's id** |
| `GOOGLETASKS_MOVE_TASK` | Moves a task within or between lists | write | |

**What to keep:** the id of the list the taf writes to, and of each task it
adds, so it marks that task done later rather than finding it by title.

**Due dates are dates only.** The time part of a due date is ignored; a task
can't be due at 15:00. A job that needs a time uses Google Calendar.

## Don't offer

- `GOOGLETASKS_CLEAR_TASKS`: it hides every completed task in a list at once.
- `GOOGLETASKS_DELETE_TASK_LIST`: it deletes a list with every task in it.
- `GOOGLETASKS_DELETE_TASK`: it removes a task the owner may still need;
  mark it done instead.
- `GOOGLETASKS_UPDATE_TASK_FULL` and `GOOGLETASKS_UPDATE_TASK_LIST`: they
  replace everything, and a field left out is erased. Use the patch tools.
- `GOOGLETASKS_BATCH_EXECUTE`, unless the job writes many tasks at once and
  the owner sees the list first.

## Access and plans

- **What Google's page asks for:** planque has no reviewed order of Google
  Tasks's
  access yet, so the sign-in asks for what the connection service recommends
  for the taf's tools. That can be more than they need, and the approval says
  so. Name the taf's tools when it connects.
- **No paid plan is needed.** A business admin can block outside apps.
- **Limits:** up to 100,000 tasks in all, 20,000 open tasks in one list,
  and 2,000 lists. Google's daily call allowance is shared by everyone on
  the same connection, so a taf keeps its calls few.
- **Assigned tasks** (from a Doc or a Chat space) carry a link back to where
  they were assigned, and some of their fields can't be changed here.

## When it fails

The rows are Google Tasks's own answers: a refusal reaches the taf as
"Google Tasks refused it (it answered …): …", with Google Tasks's reason.
planque's own refusals, and an expired sign-in, are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in gave read access only | "Google Tasks is only partly connected: adding tasks needs access you didn't give. Sign in again and allow it." |
| Refused (404), task or list not found | Deleted in the app | "That task isn't there any more. The taf will skip it." |
| Refused (429), too many calls | Over Google's limit | "Google asked us to slow down. The taf will continue in a minute." |
| Google Tasks's sign-in has expired (401) | Removed, or expired | "Google Tasks needs you to sign in again. Nothing was changed." |
