# Google Tasks: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

16 tools: 5 read, 11 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLETASKS_BATCH_EXECUTE` | Batch Execute Google Tasks Operations | write | `requests`; *max_requests* |  | `tasks` |
| `GOOGLETASKS_CLEAR_TASKS` | Clear tasks | write | `tasklist` |  | `tasks` |
| `GOOGLETASKS_CREATE_TASK_LIST` | Create a task list | write | `tasklist_title` |  | `tasks` |
| `GOOGLETASKS_DELETE_TASK` | Delete task | write | `task_id`, `tasklist_id` |  | `tasks` |
| `GOOGLETASKS_DELETE_TASK_LIST` | Delete task list | write | `tasklist_id` |  | `tasks` |
| `GOOGLETASKS_GET_TASK` | Get Task | read | `task_id`, `tasklist_id` |  | `tasks`, `tasks.readonly` |
| `GOOGLETASKS_GET_TASK_LIST` | Get task list | read | `tasklist_id` |  | `tasks`, `tasks.readonly` |
| `GOOGLETASKS_INSERT_TASK` | Insert Task | write | `tasklist_id`, `title`; *completed, deleted, due, etag, hidden, id, notes, status and 2 more* |  | `tasks` |
| `GOOGLETASKS_LIST_ALL_TASKS` | List All Tasks Across All Lists | read | *completedMax, completedMin, dueMax, dueMin, max_tasks_total, showAssigned, showCompleted, showDeleted and 2 more* |  | `tasks` |
| `GOOGLETASKS_LIST_TASKS` | List Tasks | read | `tasklist_id`; *completedMax, completedMin, dueMax, dueMin, maxResults, pageToken, showCompleted, showDeleted and 2 more* |  | `tasks` |
| `GOOGLETASKS_LIST_TASK_LISTS` | List task lists | read | *maxResults, pageToken* |  | `tasks`, `tasks.readonly` |
| `GOOGLETASKS_MOVE_TASK` | Move Task | write | `task`, `tasklist`; *destinationTasklist, parent, previous* |  | `tasks` |
| `GOOGLETASKS_PATCH_TASK` | Patch Task | write | `task_id`, `tasklist_id`; *completed, deleted, due, etag, hidden, id, notes, status and 1 more* |  | `tasks` |
| `GOOGLETASKS_PATCH_TASK_LIST` | Patch task list | write | `tasklist_id`; *updated_title* |  | `tasks` |
| `GOOGLETASKS_UPDATE_TASK_FULL` | Update Task (Full Replacement) | write | `id`, `task_id`, `tasklist_id`, `title`; *completed, deleted, due, etag, hidden, links, notes, parent and 2 more* |  | `tasks` |
| `GOOGLETASKS_UPDATE_TASK_LIST` | Update Task List | write | `tasklist_id`, `title` |  | `tasks` |
