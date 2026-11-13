# Asana: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

247 tools: 109 read, 138 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ASANA_ADD_CUSTOM_FIELD_SETTING_FOR_GOAL` | Add a custom field to a goal | write | `data`, `goal_gid`; *opt_pretty* |  | `goals:write`, `default` |
| `ASANA_ADD_CUSTOM_FIELD_SETTING_FOR_PORTFOLIO` | Add a custom field to a portfolio | write | `data`, `portfolio_gid`; *opt_pretty* |  | `portfolios:write`, `default` |
| `ASANA_ADD_CUSTOM_FIELD_SETTING_FOR_PROJECT` | Add a custom field to a project | write | `data`, `project_gid`; *opt_fields, opt_pretty* |  | `projects:write`, `default` |
| `ASANA_ADD_FOLLOWERS_TO_GOAL` | Add a collaborator to a goal | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_ADD_FOLLOWERS_TO_PROJECT` | Add Followers to Project | write | `followers`, `project_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_ADD_FOLLOWERS_TO_TASK` | Add Followers to Task | write | `followers`, `task_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_ADD_ITEM_TO_PORTFOLIO` | Add item to portfolio | write | `item`, `portfolio_gid`; *insert_after, insert_before* |  | `portfolios:write` |
| `ASANA_ADD_MEMBERS_FOR_PORTFOLIO` | Add users to a portfolio | write | `data`, `portfolio_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_ADD_MEMBERS_TO_PROJECT` | Add Members to Project | write | `members`, `project_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_ADD_PROJECT_FOR_TASK` | Add Project to Task | write | `project`, `task_gid`; *insert_after, insert_before, opt_pretty, section* |  | `tasks:write` |
| `ASANA_ADD_SUPPORTING_RELATIONSHIP` | Add Supporting Relationship to Goal | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_ADD_TAG_TO_TASK` | Add Tag to Task | write | `tag_gid`, `task_gid`; *opt_fields, opt_pretty* |  | `tasks:write` |
| `ASANA_ADD_TASK_DEPENDENCIES` | Add Task Dependencies | write | `dependencies`, `task_gid`; *opt_fields, opt_pretty* |  | `tasks:write` |
| `ASANA_ADD_TASK_DEPENDENTS` | Add Task Dependents | write | `dependents`, `task_gid`; *opt_pretty* |  | `tasks:write`, `default` |
| `ASANA_ADD_TASK_TO_SECTION` | Add task to section | write | `section_gid`, `task_gid`; *insert_after, insert_before* |  |  |
| `ASANA_ADD_USER_TO_TEAM` | Add User for Team | write | `team_gid`, `user`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_ADD_USER_TO_WORKSPACE` | Add User for Workspace | write | `user`, `workspace_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_APPROVE_ACCESS_REQUEST` | Approve Access Request | write | `access_request_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_CREATE_ACCESS_REQUEST` | Create Access Request | write | `target`; *message* |  | `default` |
| `ASANA_CREATE_ALLOCATION` | Create Allocation | write | `assignee_gid`, `effort_type`, `effort_value`, `end_on`, `start_on`, `workspace_gid`; *notes, project_gid* |  |  |
| `ASANA_CREATE_ATTACHMENT_FOR_OBJECT` | Create Attachment for Object | write | `parent`; *connect_to_app, file, name, resource_subtype, url* |  | `attachments:write` |
| `ASANA_CREATE_ATTACHMENT_FOR_TASK` | Create Attachment for Task | write | `file`, `parent_gid`; *connect_to_app, resource_subtype* |  |  |
| `ASANA_CREATE_A_PROJECT` | Create a project | write | `data`; *opt_fields, opt_pretty* |  |  |
| `ASANA_CREATE_A_TAG_IN_A_WORKSPACE` | Create a tag in a workspace | write | `data`, `workspace_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_CREATE_A_TASK` | Create task in asana with specific details | write | `data`; *opt_fields, opt_pretty* |  | `tasks:write`, `default` |
| `ASANA_CREATE_BUDGET` | Create a budget | write | `data`; *opt_pretty* |  | `default` |
| `ASANA_CREATE_CUSTOM_FIELD` | Create Custom Field | write | `data`; *opt_fields, opt_pretty* |  |  |
| `ASANA_CREATE_ENUM_OPTION_FOR_CUSTOM_FIELD` | Create Enum Option for Custom Field | write | `custom_field_gid`, `data`; *opt_fields, opt_pretty* |  |  |
| `ASANA_CREATE_GOAL` | Create a goal | write | `data`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_CREATE_GOAL_METRIC` | Create a goal metric | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_CREATE_MEMBERSHIP` | Create Membership | write | `member`, `parent`; *access_level, opt_pretty, role* |  | `default` |
| `ASANA_CREATE_OOO_ENTRY` | Create an OOO entry | write | `data`; *opt_fields, opt_pretty* |  | `ooo_entries:write`, `default` |
| `ASANA_CREATE_PORTFOLIO` | Create a portfolio | write | `data`; *opt_fields, opt_pretty* |  | `portfolios:write`, `default` |
| `ASANA_CREATE_PROJECT_BRIEF` | Create Project Brief | write | `html_text`, `project_gid`; *opt_fields, opt_pretty* |  | `default`, `projects:write` |
| `ASANA_CREATE_PROJECT_FOR_TEAM` | Create project for team | write | `data`, `team_gid`; *opt_fields, opt_pretty* |  | `projects:write` |
| `ASANA_CREATE_PROJECT_FOR_WORKSPACE` | Create project for workspace | write | `data`, `workspace_gid`; *opt_fields, opt_pretty* |  | `projects:write` |
| `ASANA_CREATE_PROJECT_STATUS_UPDATE` | Create Project Status Update | write | `color`, `project_gid`, `status_type`, `text`, `title` |  |  |
| `ASANA_CREATE_RATE` | Create a rate | write | `data`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_CREATE_ROLE` | Create a role | write | `data`; *opt_fields, opt_pretty* |  | `roles:write`, `default` |
| `ASANA_CREATE_SECTION_IN_PROJECT` | Create a section in a project | write | `name`, `project_gid`; *insert_after, insert_before* |  |  |
| `ASANA_CREATE_STATUS_FOR_OBJECT` | Create Status Update for Object | write | `parent`, `status_type`, `text`; *opt_fields, opt_pretty, title* |  | `default` |
| `ASANA_CREATE_STORY_FOR_GOAL` | Create a story on a goal | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  | `stories:write`, `default` |
| `ASANA_CREATE_SUBTASK` | Create subtask | write | `name`, `task_gid`; *assignee, completed, due_at, due_on, followers, notes* |  |  |
| `ASANA_CREATE_TAG` | Create a tag | write | `data`; *opt_fields, opt_pretty* |  | `tags:write` |
| `ASANA_CREATE_TASK_COMMENT` | Create task comment | write | `task_id`, `text` |  |  |
| `ASANA_CREATE_TEAM` | Create Team | write | `data` |  |  |
| `ASANA_CREATE_TIMESHEET_APPROVAL_STATUS` | Create a timesheet approval status | write | `data`; *opt_fields, opt_pretty* |  | `timesheet_approval_statuses:write`, `default` |
| `ASANA_CREATE_TIME_TRACKING_CATEGORY` | Create a time tracking category | write | `data`; *opt_fields, opt_pretty* |  | `time_tracking_categories:write`, `default` |
| `ASANA_CREATE_TIME_TRACKING_ENTRY` | Create a time tracking entry | write | `data`, `task_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_DELETE_ALLOCATION` | Delete Allocation | write | `allocation_gid`; *opt_pretty* |  |  |
| `ASANA_DELETE_ATTACHMENT` | Delete Attachment | write | `attachment_gid`; *opt_pretty* |  |  |
| `ASANA_DELETE_BUDGET` | Delete a budget | write | `budget_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_CUSTOM_FIELD` | Delete Custom Field | write | `custom_field_gid` |  |  |
| `ASANA_DELETE_GOAL` | Delete a goal | write | `goal_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_MEMBERSHIP` | Delete Membership | write | `membership_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_OOO_ENTRY` | Delete an OOO entry | write | `ooo_entry_gid`; *opt_pretty* |  | `ooo_entries:delete`, `default` |
| `ASANA_DELETE_PORTFOLIO` | Delete a portfolio | write | `portfolio_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_PROJECT` | Delete a project | write | `project_gid`; *opt_pretty* |  |  |
| `ASANA_DELETE_PROJECT_BRIEF` | Delete Project Brief | write | `project_brief_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_PROJECT_STATUS` | Delete Project Status | write | `project_status_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_PROJECT_TEMPLATE` | Delete a project template | write | `project_template_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_RATE` | Delete a rate | write | `rate_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_ROLE` | Delete a role | write | `role_gid`; *opt_pretty* |  | `roles:delete`, `default` |
| `ASANA_DELETE_SECTION` | Delete a section | write | `section_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_STATUS_UPDATE` | Delete Status Update | write | `status_update_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_STORY` | Delete Story | write | `story_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_TAG` | Delete a Tag | write | `tag_gid`; *opt_pretty* |  |  |
| `ASANA_DELETE_TASK` | Delete a task | write | `task_gid`; *opt_pretty* |  |  |
| `ASANA_DELETE_TASK_TEMPLATE` | Delete a task template | write | `task_template_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_TIME_TRACKING_CATEGORY` | Delete a time tracking category | write | `time_tracking_category_gid`; *opt_pretty* |  | `time_tracking_categories:delete`, `default` |
| `ASANA_DELETE_TIME_TRACKING_ENTRY` | Delete a time tracking entry | write | `time_tracking_entry_gid`; *opt_pretty* |  | `default` |
| `ASANA_DELETE_WEBHOOK` | Delete a webhook | write | `webhook_gid`; *opt_pretty* |  | `webhooks:delete`, `default` |
| `ASANA_DUPLICATE_PORTFOLIO` | Duplicate a portfolio | write | `portfolio_gid`; *data, opt_fields, opt_pretty* |  | `portfolios:write`, `default` |
| `ASANA_DUPLICATE_PROJECT` | Duplicate Project | write | `name`, `project_gid`; *include, schedule_dates, team* |  |  |
| `ASANA_DUPLICATE_TASK` | Duplicate Task | write | `task_gid`; *include, name* |  |  |
| `ASANA_GET_ACCESS_REQUESTS` | Get Access Requests | read | `target`; *opt_fields, opt_pretty, user* |  | `default` |
| `ASANA_GET_AI_STUDIO_RUNS` | Get AI Studio credit utilization | read | `workspace_gid`; *division_gid, end_at, limit, offset, start_at* |  | `default` |
| `ASANA_GET_AI_STUDIO_SEATS` | Get AI Studio seats | read | `workspace_gid`; *division_gid, limit, offset, state* |  | `default` |
| `ASANA_GET_ALLOCATION` | Get Allocation | read | `allocation_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_ALLOCATIONS` | Get Allocations | read | *assignee, limit, offset, opt_fields, opt_pretty, parent, workspace* |  |  |
| `ASANA_GET_ATTACHMENT` | Get Attachment | read | `attachment_gid`; *opt_fields* |  |  |
| `ASANA_GET_AUDIT_LOG_EVENTS` | Get Audit Log Events | read | `workspace_gid`; *actor_gid, actor_type, end_at, event_type, limit, offset, resource_gid, start_at* |  |  |
| `ASANA_GET_A_PROJECT` | Get a project | read | `project_gid`; *opt_fields, opt_pretty* |  | `projects:read` |
| `ASANA_GET_A_TASK` | Get a task | read | `task_gid`; *opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_GET_A_USER_TASK_LIST` | Get a user task list | read | `user_task_list_gid`, `workspace`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_BUDGET` | Get a budget | read | `budget_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_BUDGETS` | Get all budgets | read | `parent`; *opt_pretty* |  | `default` |
| `ASANA_GET_CURRENT_USER` | Get current user | read |  |  |  |
| `ASANA_GET_CUSTOM_FIELD` | Get Custom Field | write | `custom_field_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_CUSTOM_FIELDS_FOR_WORKSPACE` | Get Custom Fields for Workspace | write | `workspace_gid`; *limit, offset, opt_fields, opt_pretty* |  |  |
| `ASANA_GET_CUSTOM_FIELD_SETTINGS_FOR_GOAL` | Get a goal's custom fields | write | `goal_gid`; *limit, offset, opt_fields, opt_pretty* |  | `goals:read`, `default` |
| `ASANA_GET_CUSTOM_FIELD_SETTINGS_FOR_PORTFOLIO` | Get a portfolio's custom fields | write | `portfolio_gid`; *limit, offset, opt_fields, opt_pretty* |  | `portfolios:read`, `default` |
| `ASANA_GET_CUSTOM_FIELD_SETTINGS_FOR_TEAM` | Get a team's custom fields | write | `team_gid`; *opt_fields, opt_pretty* |  | `teams:read`, `default` |
| `ASANA_GET_CUSTOM_TYPE` | Get a custom type | write | `custom_type_gid`; *opt_fields, opt_pretty* |  | `custom_types:read`, `default` |
| `ASANA_GET_CUSTOM_TYPES` | Get Custom Types | write | `project`; *limit, offset, opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_EVENTS` | Get Events on a Resource | read | `resource_gid`; *opt_fields, opt_pretty, sync* |  |  |
| `ASANA_GET_FAVORITES_FOR_USER` | Get Favorites for User | read | `user_gid`, `workspace`; *limit, offset, opt_fields, opt_pretty, resource_type* |  | `users:read` |
| `ASANA_GET_GOAL` | Get Goal | read | `goal_gid`; *opt_fields* |  |  |
| `ASANA_GET_GOALS` | Get Goals | read | *archived, is_workspace_level, limit, offset, opt_fields, portfolio, project, team and 2 more* |  |  |
| `ASANA_GET_GOAL_RELATIONSHIP` | Get a goal relationship | read | `goal_relationship_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_GOAL_RELATIONSHIPS` | Get Goal Relationships | read | `goal_gid`; *limit, offset, opt_fields, opt_pretty, resource_subtype* |  | `default` |
| `ASANA_GET_JOB` | Get Job | read | `job_gid`; *opt_fields, opt_pretty* |  | `jobs:read` |
| `ASANA_GET_MEMBERSHIP` | Get Membership | read | `membership_gid`; *opt_pretty* |  | `default` |
| `ASANA_GET_MEMBERSHIPS` | Get Memberships | read | *limit, member, offset, opt_fields, parent* |  |  |
| `ASANA_GET_MULTIPLE_PROJECTS` | Get multiple projects | read | *archived, limit, offset, opt_fields, opt_pretty, team, workspace* |  | `default` |
| `ASANA_GET_MULTIPLE_TASKS` | Get multiple tasks | read | *assignee, completed_since, limit, modified_since, offset, opt_fields, opt_pretty, project and 4 more* |  | `tasks:read` |
| `ASANA_GET_MULTIPLE_USERS` | Get multiple users | read | *limit, offset, opt_fields, opt_pretty, team, workspace* |  | `users:read` |
| `ASANA_GET_MULTIPLE_WORKSPACES` | Get multiple workspaces | read | *limit, offset, opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_OOO_ENTRIES` | Get OOO Entries | read | `user`, `workspace`; *end_date, limit, offset, opt_fields, opt_pretty, start_date* |  | `ooo_entries:read`, `default` |
| `ASANA_GET_OOO_ENTRY` | Get an OOO entry | read | `ooo_entry_gid`; *opt_fields, opt_pretty* |  | `ooo_entries:read`, `default` |
| `ASANA_GET_PARENT_GOALS_FOR_GOAL` | Get parent goals from a goal | read | `goal_gid`; *opt_fields, opt_pretty* |  | `goals:read`, `default` |
| `ASANA_GET_PORTFOLIO` | Get Portfolio | read | `portfolio_gid`; *limit, offset, opt_fields, opt_pretty* |  |  |
| `ASANA_GET_PORTFOLIOS` | Get Portfolios | read | `owner`, `workspace`; *limit, offset, opt_fields* |  |  |
| `ASANA_GET_PORTFOLIO_ITEMS` | Get Portfolio Items | read | `portfolio_gid`; *limit, offset, opt_fields* |  |  |
| `ASANA_GET_PORTFOLIO_MEMBERSHIP` | Get a portfolio membership | read | `portfolio_membership_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_PORTFOLIO_MEMBERSHIPS` | Get Portfolio Memberships | read | *limit, offset, opt_fields, portfolio, user, workspace* |  |  |
| `ASANA_GET_PORTFOLIO_MEMBERSHIPS_FOR_PORTFOLIO` | Get memberships from a portfolio | read | `portfolio_gid`; *limit, offset, opt_fields, opt_pretty, user* |  | `default` |
| `ASANA_GET_PROJECTS_FOR_TASK` | Get Projects for Task | read | `task_gid`; *limit, offset, opt_fields, opt_pretty* |  | `projects:read` |
| `ASANA_GET_PROJECTS_FOR_TEAM` | Get Projects for Team | read | `team_gid`; *archived, limit, offset, opt_fields, opt_pretty* |  |  |
| `ASANA_GET_PROJECT_BRIEF` | Get Project Brief | read | `project_brief_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_PROJECT_CUSTOM_FIELD_SETTINGS` | Get Project Custom Field Settings | write | `project_gid`; *limit, offset, opt_fields, opt_pretty* |  | `projects:read` |
| `ASANA_GET_PROJECT_MEMBERSHIP` | Get Project Membership | read | `project_membership_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_PROJECT_MEMBERSHIPS` | Get Project Memberships For Project | read | `project_gid`; *limit, offset, opt_fields, opt_pretty, user* |  | `default` |
| `ASANA_GET_PROJECT_PORTFOLIO_SETTING` | Get a project portfolio setting | read | `project_portfolio_setting_gid`; *opt_fields, opt_pretty* |  | `project_portfolio_settings:read`, `default` |
| `ASANA_GET_PROJECT_PORTFOLIO_SETTINGS_FOR_PORTFOLIO` | Get project portfolio settings for a portfolio | read | `portfolio_gid`; *limit, offset, opt_fields, opt_pretty* |  | `project_portfolio_settings:read`, `default` |
| `ASANA_GET_PROJECT_PORTFOLIO_SETTINGS_FOR_PROJECT` | Get project portfolio settings for a project | read | `project_gid`; *limit, offset, opt_fields, opt_pretty* |  | `project_portfolio_settings:read`, `default` |
| `ASANA_GET_PROJECT_STATUS` | Get Project Status | read | `project_status_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_PROJECT_STATUS_UPDATES` | Get Project Status Updates | read | `project_gid`; *limit, offset, opt_fields, opt_pretty* |  |  |
| `ASANA_GET_PROJECT_TEMPLATE` | Get Project Template | read | `project_template_gid`; *opt_fields, opt_pretty* |  | `project_templates:read`, `default` |
| `ASANA_GET_PROJECT_TEMPLATES` | Get Multiple Project Templates | read | *limit, offset, opt_fields, team_gid, workspace_gid* |  |  |
| `ASANA_GET_PROJECT_TEMPLATES_FOR_TEAM` | Get Project Templates for Team | read | `team_gid`; *limit, offset, opt_fields, opt_pretty* |  | `project_templates:read` |
| `ASANA_GET_RATE` | Get a rate | read | `rate_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_RATES` | Get multiple rates | read | *limit, offset, opt_fields, opt_pretty, parent, resource* |  | `default` |
| `ASANA_GET_REACTIONS_ON_OBJECT` | Get Reactions on Object | read | `emoji_base`, `target`; *limit, offset, opt_pretty* |  | `default` |
| `ASANA_GET_ROLE` | Get a role | read | `role_gid`; *opt_fields, opt_pretty* |  | `roles:read`, `default` |
| `ASANA_GET_ROLES` | Get multiple roles | read | `workspace`; *archived, limit, offset, opt_fields, opt_pretty* |  | `roles:read`, `default` |
| `ASANA_GET_SECTION` | Get Section | read | `section_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_SECTIONS_IN_PROJECT` | Get sections in a project | read | `project_gid`; *limit, offset* |  |  |
| `ASANA_GET_STATUS` | Get Status Update | read | `status_update_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_STATUS_UPDATES` | Get Status Updates | read | `parent`; *limit, offset, opt_fields* |  |  |
| `ASANA_GET_STORIES_FOR_GOAL` | Get stories from a goal | read | `goal_gid`; *created_after, limit, offset, opt_fields, opt_pretty, resource_subtype, sort_ascending* |  | `stories:read`, `default` |
| `ASANA_GET_STORIES_FOR_TASK` | Get Stories for Task | read | `task_gid`; *limit, offset, opt_fields* |  |  |
| `ASANA_GET_STORY` | Get Story | read | `story_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_TAG` | Get Tag | read | `tag_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_TAGS` | Get Tags | read | *limit, offset, opt_fields, workspace* |  |  |
| `ASANA_GET_TAGS_FOR_TASK` | Get Tags for Task | read | `task_gid`; *limit, offset, opt_fields, opt_pretty* |  | `tags:read` |
| `ASANA_GET_TAGS_FOR_WORKSPACE` | Get Tags For Workspace | read | `workspace_gid`; *limit, offset, opt_fields, opt_pretty* |  | `tags:read` |
| `ASANA_GET_TASKS_FOR_TAG` | Get tasks for tag | read | `tag_gid`; *limit, offset, opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_TASKS_FOR_USER_TASK_LIST` | Get tasks for user task list | read | `user_task_list_gid`; *completed_since, limit, offset, opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_GET_TASKS_FROM_A_PROJECT` | Retrieve tasks for project | read | `project_gid`; *completed_since, limit, offset, opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_GET_TASKS_FROM_SECTION` | Get tasks from a section | read | `section_gid`; *completed_since, limit, offset, opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_GET_TASK_ATTACHMENTS` | Get Task Attachments | read | `parent_gid`; *limit, offset, opt_fields* |  |  |
| `ASANA_GET_TASK_COUNTS_FOR_PROJECT` | Get Task Counts For Project | read | `project_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_TASK_FOR_CUSTOM_ID` | Get a task for a given custom ID | write | `custom_id`, `workspace_gid` |  | `tasks:read`, `default` |
| `ASANA_GET_TASK_SUBTASKS` | Get Task Subtasks | read | `task_gid`; *limit, offset, opt_fields* |  |  |
| `ASANA_GET_TASK_TEMPLATE` | Get a task template | read | `task_template_gid`; *opt_fields, opt_pretty* |  | `task_templates:read`, `default` |
| `ASANA_GET_TASK_TEMPLATES` | Get Task Templates | read | *limit, offset, opt_fields, project_gid, workspace_gid* |  |  |
| `ASANA_GET_TEAM` | Get Team | read | `team_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_TEAMS_FOR_USER` | Get Teams for User | read | `organization`, `user_gid`; *limit, offset, opt_fields, opt_pretty* |  | `teams:read` |
| `ASANA_GET_TEAMS_IN_WORKSPACE` | Get teams in workspace | read | `workspace_gid`; *limit, offset* |  |  |
| `ASANA_GET_TEAM_MEMBERSHIP` | Get Team Membership | read | `team_membership_gid`; *opt_fields, opt_pretty* |  | `team_memberships:read` |
| `ASANA_GET_TEAM_MEMBERSHIPS` | Get Team Memberships | read | *limit, offset, opt_fields, opt_pretty, team, user, workspace* |  |  |
| `ASANA_GET_TEAM_MEMBERSHIPS_FOR_TEAM` | Get Team Memberships for Team | read | `team_gid`; *limit, offset, opt_fields, opt_pretty* |  | `team_memberships:read` |
| `ASANA_GET_TEAM_MEMBERSHIPS_FOR_USER` | Get Team Memberships For User | read | `user_gid`, `workspace`; *limit, offset, opt_fields, opt_pretty* |  | `team_memberships:read` |
| `ASANA_GET_TIMESHEET_APPROVAL_STATUS` | Get a timesheet approval status | read | `timesheet_approval_status_gid`; *opt_fields, opt_pretty* |  | `timesheet_approval_statuses:read`, `default` |
| `ASANA_GET_TIMESHEET_APPROVAL_STATUSES` | Get multiple timesheet approval statuses | read | `workspace`; *approval_statuses, from_date, limit, offset, opt_fields, opt_pretty, to_date, user* |  | `timesheet_approval_statuses:read`, `default` |
| `ASANA_GET_TIME_PERIOD` | Get Time Period | read | `time_period_gid`; *opt_fields, opt_pretty* |  | `time_periods:read` |
| `ASANA_GET_TIME_PERIODS` | Get Time Periods | read | `workspace_gid`; *end_on, limit, offset, opt_fields, opt_pretty, start_on* |  |  |
| `ASANA_GET_TIME_TRACKING_CATEGORIES` | Get time tracking categories for a workspace | read | `workspace`; *is_archived, limit, offset, opt_fields, opt_pretty* |  | `time_tracking_categories:read`, `default` |
| `ASANA_GET_TIME_TRACKING_CATEGORY` | Get a time tracking category | read | `time_tracking_category_gid`; *opt_fields, opt_pretty* |  | `time_tracking_categories:read`, `default` |
| `ASANA_GET_TIME_TRACKING_ENTRIES` | Get Time Tracking Entries | read | *attributable_to, entered_on_end_date, entered_on_start_date, limit, offset, opt_fields, opt_pretty, portfolio and 4 more* |  | `default` |
| `ASANA_GET_TIME_TRACKING_ENTRIES_FOR_TASK` | Get Time Tracking Entries for Task | read | `task_gid`; *limit, offset, opt_fields, opt_pretty* |  | `time_tracking_entries:read`, `default` |
| `ASANA_GET_TIME_TRACKING_ENTRIES_FOR_TIME_TRACKING_CATEGORY` | Get time tracking entries for a time tracking category | read | `time_tracking_category_gid`; *end_date, limit, offset, opt_fields, opt_pretty, start_date* |  | `time_tracking_categories:read`, `default` |
| `ASANA_GET_TIME_TRACKING_ENTRY` | Get a time tracking entry | read | `time_tracking_entry_gid`; *opt_fields, opt_pretty* |  | `time_tracking_entries:read`, `default` |
| `ASANA_GET_TYPEAHEAD_OBJECTS` | Get Objects via Typeahead | read | `workspace_gid`; *count, opt_fields, opt_pretty, query, resource_type* |  |  |
| `ASANA_GET_USER` | Get User | read | `user_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_USERS_FOR_TEAM` | Get Users for Team | read | `team_gid`; *limit, offset, opt_fields, opt_pretty* |  |  |
| `ASANA_GET_USERS_FOR_WORKSPACE` | Get Users in Workspace | read | `workspace_gid`; *limit, offset, opt_fields, opt_pretty* |  |  |
| `ASANA_GET_USER_FOR_WORKSPACE` | Get User in Workspace | read | `user_gid`, `workspace_gid`; *opt_fields, opt_pretty* |  | `users:read`, `default` |
| `ASANA_GET_USER_TASK_LIST` | Get User Task List | read | `user_task_list_gid`; *opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_GET_WEBHOOK` | Get a webhook | read | `webhook_gid`; *opt_fields, opt_pretty* |  | `webhooks:read`, `default` |
| `ASANA_GET_WEBHOOKS` | Get Webhooks | read | `workspace`; *limit, offset, opt_fields, opt_pretty, resource* |  | `webhooks:read` |
| `ASANA_GET_WORKSPACE` | Get Workspace | read | `workspace_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_GET_WORKSPACE_EVENTS` | Get workspace events | read | `workspace_gid`; *opt_pretty, sync* |  | `default` |
| `ASANA_GET_WORKSPACE_MEMBERSHIP` | Get Workspace Membership | read | `workspace_membership_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_WORKSPACE_MEMBERSHIPS` | Get Workspace Memberships | read | `workspace_gid`; *limit, offset, opt_fields, user* |  |  |
| `ASANA_GET_WORKSPACE_MEMBERSHIPS_FOR_USER` | Get Workspace Memberships For User | read | `user_gid`; *limit, offset, opt_fields, opt_pretty* |  | `default` |
| `ASANA_GET_WORKSPACE_PROJECTS` | Get Workspace Projects | read | `workspace_gid`; *limit, offset, opt_expand, opt_fields* |  |  |
| `ASANA_INSERT_ENUM_OPTION_FOR_CUSTOM_FIELD` | Reorder Enum Option for Custom Field | write | `custom_field_gid`, `data`; *opt_fields, opt_pretty* |  |  |
| `ASANA_INSERT_SECTION_FOR_PROJECT` | Move or Insert Section in Project | write | `data`, `project_gid`; *opt_pretty* |  | `default` |
| `ASANA_INSTANTIATE_PROJECT_TEMPLATE` | Instantiate Project Template | write | `name`, `project_template_gid`; *is_strict, opt_fields, opt_pretty, privacy_setting, requested_dates, requested_roles, team* |  | `projects:write` |
| `ASANA_INSTANTIATE_TASK_TEMPLATE` | Instantiate a task from a task template | write | `task_template_gid`; *data, opt_fields, opt_pretty* |  | `default` |
| `ASANA_LIST_GOAL_RELATIONSHIPS` | Get goal relationships | read | `supported_goal`; *limit, offset, opt_fields, opt_pretty, resource_subtype* |  | `default` |
| `ASANA_LIST_TASK_DEPENDENCIES` | List Task Dependencies | read | `task_gid`; *limit, offset, opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_LIST_TASK_DEPENDENTS` | List Task Dependents | read | `task_gid`; *limit, offset, opt_fields, opt_pretty* |  | `tasks:read` |
| `ASANA_REJECT_ACCESS_REQUEST` | Reject access request | write | `access_request_gid` |  | `default` |
| `ASANA_REMOVE_CUSTOM_FIELD_SETTING_FOR_GOAL` | Remove a custom field from a goal | write | `data`, `goal_gid`; *opt_pretty* |  | `goals:write`, `default` |
| `ASANA_REMOVE_CUSTOM_FIELD_SETTING_FOR_PORTFOLIO` | Remove a custom field from a portfolio | write | `data`, `portfolio_gid`; *opt_pretty* |  | `portfolios:write`, `default` |
| `ASANA_REMOVE_CUSTOM_FIELD_SETTING_FOR_PROJECT` | Remove a custom field from a project | write | `data`, `project_gid`; *opt_pretty* |  | `projects:write`, `default` |
| `ASANA_REMOVE_FOLLOWERS_FOR_PROJECT` | Remove Followers For Project | write | `followers`, `project_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_REMOVE_FOLLOWERS_FROM_GOAL` | Remove a collaborator from a goal | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_REMOVE_FOLLOWER_FROM_TASK` | Remove Follower From Task | write | `followers`, `task_gid` |  |  |
| `ASANA_REMOVE_ITEM_FROM_PORTFOLIO` | Remove Item from Portfolio | write | `item`, `portfolio_gid`; *opt_fields, opt_pretty* |  | `portfolios:write` |
| `ASANA_REMOVE_MEMBERS_FOR_PORTFOLIO` | Remove users from a portfolio | write | `data`, `portfolio_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_REMOVE_MEMBERS_FOR_PROJECT` | Remove Members from Project | write | `members`, `project_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_REMOVE_PROJECT_FROM_TASK` | Remove Project from Task | write | `project`, `task_gid`; *opt_pretty* |  | `tasks:write` |
| `ASANA_REMOVE_SUPPORTING_RELATIONSHIP` | Removes a supporting goal relationship | write | `data`, `goal_gid`; *opt_pretty* |  | `default` |
| `ASANA_REMOVE_TAG_FROM_TASK` | Remove Tag from Task | write | `tag`, `task_gid`; *opt_fields, opt_pretty* |  | `tasks:write` |
| `ASANA_REMOVE_TASK_DEPENDENCIES` | Remove Task Dependencies | write | `dependencies`, `task_gid`; *opt_pretty* |  | `tasks:write` |
| `ASANA_REMOVE_TASK_DEPENDENTS` | Remove Task Dependents | write | `dependents`, `task_gid`; *opt_pretty* |  | `tasks:write`, `default` |
| `ASANA_REMOVE_USER_FOR_TEAM` | Remove User For Team | write | `team_gid`, `user`; *opt_pretty* |  | `default` |
| `ASANA_REMOVE_USER_FOR_WORKSPACE` | Remove User from Workspace | write | `user`, `workspace_gid`; *opt_pretty* |  | `default` |
| `ASANA_ROLLUP_PROJECT` | Roll up subtask dates for a project | write | `project_gid`; *opt_fields, opt_pretty* |  | `projects:write`, `default` |
| `ASANA_ROLLUP_TASK` | Roll up subtask dates for a task | write | `task_gid`; *opt_fields, opt_pretty* |  | `tasks:write`, `default` |
| `ASANA_SAVE_PROJECT_AS_TEMPLATE` | Create a project template from a project | write | `data`, `project_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_SEARCH_PROJECTS_FOR_WORKSPACE` | Search projects in a workspace | read | `workspace_gid`; *completed, completed_at.after, completed_at.before, completed_on, completed_on.after, completed_on.before, created_at.after, created_at.before and 22 more* |  | `projects:read`, `default` |
| `ASANA_SEARCH_TASKS_IN_WORKSPACE` | Search Tasks in Workspace | read | `workspace_gid`; *assignee_any, composio_execution_message, limit, opt_fields, opt_pretty, projects_any, resource_subtype, text* |  |  |
| `ASANA_SET_PARENT_FOR_TASK` | Set parent for task | write | `parent`, `task_gid`; *insert_after, insert_before, opt_fields, opt_pretty* |  | `tasks:write`, `default` |
| `ASANA_SUBMIT_PARALLEL_REQUESTS` | Submit Parallel Requests (Batch API) | write | `data` |  |  |
| `ASANA_TRIGGER_RULE` | Trigger a rule | write | `data`, `rule_trigger_gid` |  | `default` |
| `ASANA_UPDATE_ALLOCATION` | Update Allocation | write | `allocation_gid`, `data`; *opt_fields, opt_pretty* |  |  |
| `ASANA_UPDATE_A_TASK` | Update a task | write | `data`, `task_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_UPDATE_BUDGET` | Update a budget | write | `budget_gid`, `data`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_CUSTOM_FIELD` | Update Custom Field | write | `custom_field_gid`, `data`; *opt_fields, opt_pretty* |  |  |
| `ASANA_UPDATE_ENUM_OPTION` | Update Enum Option | write | `data`, `enum_option_gid`; *opt_fields, opt_pretty* |  |  |
| `ASANA_UPDATE_GOAL` | Update a goal | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_GOAL_METRIC` | Update a goal metric | write | `data`, `goal_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_GOAL_RELATIONSHIP` | Update a goal relationship | write | `data`, `goal_relationship_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_MEMBERSHIP` | Update a membership | write | `data`, `membership_gid`; *opt_pretty* |  | `default` |
| `ASANA_UPDATE_OOO_ENTRY` | Update an OOO entry | write | `data`, `ooo_entry_gid`; *opt_fields, opt_pretty* |  | `ooo_entries:write`, `default` |
| `ASANA_UPDATE_PORTFOLIO` | Update a portfolio | write | `data`, `portfolio_gid`; *opt_fields, opt_pretty* |  | `portfolios:write`, `default` |
| `ASANA_UPDATE_PROJECT` | Update a project | write | `project_gid`; *archived, color, custom_fields, default_view, due_date, due_on, html_notes, is_template and 6 more* |  |  |
| `ASANA_UPDATE_PROJECT_BRIEF` | Update Project Brief | write | `project_brief_gid`; *opt_fields, opt_pretty, text* |  | `default` |
| `ASANA_UPDATE_PROJECT_PORTFOLIO_SETTING` | Update a project portfolio setting | write | `data`, `project_portfolio_setting_gid`; *opt_fields, opt_pretty* |  | `project_portfolio_settings:write`, `default` |
| `ASANA_UPDATE_RATE` | Update a rate | write | `data`, `rate_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_ROLE` | Update a role | write | `data`, `role_gid`; *opt_fields, opt_pretty* |  | `roles:write`, `default` |
| `ASANA_UPDATE_SECTION` | Update Section | write | `section_gid`; *insert_after, insert_before, name, opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_STORY` | Update Story | write | `data`, `story_gid`; *opt_fields, opt_pretty* |  | `stories:write` |
| `ASANA_UPDATE_TAG` | Update Tag | write | `data`, `tag_gid` |  |  |
| `ASANA_UPDATE_TEAM` | Update a team | write | `team_gid`; *description, name, opt_fields, opt_pretty, organization* |  |  |
| `ASANA_UPDATE_TIMESHEET_APPROVAL_STATUS` | Update a timesheet approval status | write | `data`, `timesheet_approval_status_gid`; *opt_fields, opt_pretty* |  | `timesheet_approval_statuses:write`, `default` |
| `ASANA_UPDATE_TIME_TRACKING_CATEGORY` | Update a time tracking category | write | `data`, `time_tracking_category_gid`; *opt_fields, opt_pretty* |  | `time_tracking_categories:write`, `default` |
| `ASANA_UPDATE_TIME_TRACKING_ENTRY` | Update a time tracking entry | write | `data`, `time_tracking_entry_gid`; *opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_USER` | Update User | write | `user_gid`; *custom_fields, opt_fields, opt_pretty, workspace* |  | `default` |
| `ASANA_UPDATE_USER_FOR_WORKSPACE` | Update User for Workspace | write | `user_gid`, `workspace_gid`; *custom_fields, opt_fields, opt_pretty* |  | `default` |
| `ASANA_UPDATE_WEBHOOK` | Update a webhook | write | `webhook_gid`; *filters, opt_fields, opt_pretty* |  | `webhooks:write` |
| `ASANA_UPDATE_WORKSPACE` | Update a workspace | write | `data`, `workspace_gid`; *opt_fields, opt_pretty* |  | `default` |
