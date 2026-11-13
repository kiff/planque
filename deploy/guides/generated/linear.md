# Linear: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

47 tools: 22 read, 25 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `LINEAR_ARCHIVE_ISSUE` | Archive issue | write | `issue_id` |  | `write` |
| `LINEAR_ARCHIVE_PROJECT` | Archive Project | write | `project_id` |  | `write` |
| `LINEAR_ARCHIVE_PROJECT_UPDATE` | Archive Project Update | write | `project_update_id` |  | `write` |
| `LINEAR_CREATE_ATTACHMENT` | Create attachment | write | `issue_id`, `title`, `url`; *subtitle* |  |  |
| `LINEAR_CREATE_COMMENT_REACTION` | Add reaction to comment | write | `comment_id`, `emoji` |  | `write` |
| `LINEAR_CREATE_LINEAR_COMMENT` | Create a comment | write | `body`, `issueId` |  | `write` |
| `LINEAR_CREATE_LINEAR_ISSUE` | Create linear issue | write | `team_id`, `title`; *assignee_id, cycle_id, description, due_date, estimate, label_ids, parent_id, priority and 2 more* |  | `write` |
| `LINEAR_CREATE_LINEAR_ISSUE_RELATION` | Create issue relation | write | `issue_id`, `related_issue_id`, `relation_type` |  | `write` |
| `LINEAR_CREATE_LINEAR_LABEL` | Create a label | write | `color`, `name`, `team_id`; *description* |  |  |
| `LINEAR_CREATE_LINEAR_PROJECT` | Create Project | write | `name`, `team_ids`; *color, description, icon, lead_id, priority, start_date, target_date* |  | `write` |
| `LINEAR_CREATE_PROJECT_MILESTONE` | Create Project Milestone | write | `name`, `project_id`; *description, sort_order, target_date* |  | `write` |
| `LINEAR_CREATE_PROJECT_UPDATE` | Create Project Update | write | `body`, `project_id`; *health, is_diff_hidden* |  | `write` |
| `LINEAR_CREATE_TEAM` | Create Team | write | `key`, `name`; *color, description, icon, private* |  | `write` |
| `LINEAR_DELETE_LINEAR_ISSUE` | Delete issue | write | `issue_id` |  | `write` |
| `LINEAR_GET_ATTACHMENT` | Download issue files | read | `file_name`; *attachment_id, issue_id, upload_url* |  | `read` |
| `LINEAR_GET_COMMENT` | Get a comment | read | `id` |  | `read` |
| `LINEAR_GET_CURRENT_USER` | Get current user | read |  |  | `read` |
| `LINEAR_GET_CYCLES_BY_TEAM_ID` | Get cycles by team ID | read | `team_id`; *after, first* |  | `read`, `write` |
| `LINEAR_GET_ISSUE_DEFAULTS` | Get create issue default params | read | `team_id` |  | `issues:create` |
| `LINEAR_GET_LINEAR_ISSUE` | Get Linear issue | read | `issue_id` |  | `read`, `write` |
| `LINEAR_GET_LINEAR_PROJECT` | Get Linear project | read | `project_id`; *include_initiatives, include_members, include_teams* |  | `read` |
| `LINEAR_GET_PROJECT_UPDATE` | Get Project Update | read | `id` |  | `read` |
| `LINEAR_LIST_COMMENTS` | List comments | read | *after, before, first, includeArchived* |  | `read` |
| `LINEAR_LIST_ISSUES_BY_TEAM_ID` | List issues by team ID | read | `team_id`; *after, first, include_archived* |  | `read` |
| `LINEAR_LIST_ISSUE_DRAFTS` | List issue drafts | read | *after, first* |  | `read` |
| `LINEAR_LIST_LINEAR_CYCLES` | Get all cycles | read | *after, first* |  | `read` |
| `LINEAR_LIST_LINEAR_ISSUES` | List Linear issues | read | *after, assignee_id, cursor_was_corrupted, first, include_transitions, original_cursor, project_id* |  | `read` |
| `LINEAR_LIST_LINEAR_LABELS` | Get labels | read | *after, first, team_id* |  | `read` |
| `LINEAR_LIST_LINEAR_PROJECTS` | List linear projects | read | *after, first* |  | `read` |
| `LINEAR_LIST_LINEAR_STATES` | List Linear states | read | `team_id`; *after, first* |  | `read` |
| `LINEAR_LIST_LINEAR_TEAMS` | Get teams | read | *after, first, project_id* |  | `read` |
| `LINEAR_LIST_LINEAR_USERS` | List Linear users | read | *after, first* |  | `read` |
| `LINEAR_LIST_PROJECT_UPDATES` | List project updates | read | *after, first* |  | `read` |
| `LINEAR_REMOVE_ISSUE_LABEL` | Remove label from Linear issue | write | `issue_id`, `label_id` |  |  |
| `LINEAR_REMOVE_REACTION` | Remove reaction from comment | write | `reaction_id` |  | `write` |
| `LINEAR_RESOLVE_COMMENT` | Resolve comment | write | `id` |  | `write` |
| `LINEAR_RUN_QUERY_OR_MUTATION` | Run Query or Mutation | write | `query_or_mutation`; *variables* |  | `read` |
| `LINEAR_SEARCH_ISSUES` | Search Linear issues | read | `query`; *after, first, include_archived* |  | `read` |
| `LINEAR_UNARCHIVE_PROJECT_UPDATE` | Unarchive Project Update | write | `id` |  | `write` |
| `LINEAR_UNRESOLVE_COMMENT` | Unresolve a comment | write | `comment_id` |  | `write` |
| `LINEAR_UPDATE_ISSUE` | Update issue | write | `issueId`; *assigneeId, cycleId, description, dueDate, estimate, labelIds, parentId, priority and 4 more* |  |  |
| `LINEAR_UPDATE_LINEAR_COMMENT` | Update a comment | write | `body`, `comment_id` |  | `read` |
| `LINEAR_UPDATE_LINEAR_PROJECT` | Update Project | write | `project_id`; *color, description, icon, lead_id, name, priority, start_date, state and 2 more* |  | `write` |
| `LINEAR_UPDATE_PROJECT_UPDATE` | Update Project Update | write | `id`; *body, health, is_diff_hidden* |  | `write` |
| `LINEAR_UPDATE_TEAM` | Update Team | write | `team_id`; *autoArchivePeriod, autoClosePeriod, autoCloseStateId, color, cycleCooldownTime, cycleDuration, cycleIssueAutoAssignCompleted, cycleIssueAutoAssignStarted and 14 more* |  | `write` |
| `LINEAR_VALIDATE_CREDENTIAL` | Validate credential | read |  |  | `read` |
| `LINEAR_WHO_AM_I` | Who Am I | read |  |  |  |
