# ClickUp: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

163 tools: 61 read, 102 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `CLICKUP_ADD_DEPENDENCY` | Add dependency | write | `task_id`; *custom_task_ids, dependency_of, depends_on, team_id* |  |  |
| `CLICKUP_ADD_GUEST_TO_FOLDER` | Add guest to folder | write | `folder_id`, `guest_id`, `permission_level`; *include_shared* |  |  |
| `CLICKUP_ADD_GUEST_TO_LIST` | Add guest to list | write | `guest_id`, `list_id`, `permission_level`; *include_shared* |  |  |
| `CLICKUP_ADD_GUEST_TO_TASK` | Add guest to task | write | `guest_id`, `permission_level`, `task_id`; *custom_task_ids, include_shared, team_id* |  |  |
| `CLICKUP_ADD_TAGS_TO_TIME_ENTRIES` | Add tags from time entries | write | `tags`, `team_id`, `time_entry_ids` |  |  |
| `CLICKUP_ADD_TAG_TO_TASK` | Add tag to task | write | `tag_name`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_ADD_TASK_LINK` | Add task link | write | `links_to`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_ADD_TASK_TO_LIST` | Add task to list | write | `list_id`, `task_id` |  |  |
| `CLICKUP_CLICK_UP_UPDATE_DOC_PAGE` | Update Doc Page | write | `doc_id`, `page_id`, `workspace_id`; *content, name* |  |  |
| `CLICKUP_CREATE_A_TIME_ENTRY` | Create a time entry | write | `duration`, `start`, `team_Id`; *assignee, billable, custom_task_ids, description, end, stop, tags, team_id and 1 more* |  |  |
| `CLICKUP_CREATE_CHAT_CHANNEL` | Create chat channel | write | `name`, `workspace_id`; *description, topic, user_ids, visibility* |  |  |
| `CLICKUP_CREATE_CHAT_MESSAGE` | Create chat message | write | `channel_id`, `content`, `type`, `workspace_id`; *assignee, content_format, followers, group_assignee, post_data, reactions, triaged_action, triaged_object_id and 1 more* |  |  |
| `CLICKUP_CREATE_CHAT_REACTION` | Create chat reaction | write | `message_id`, `reaction`, `workspace_id` |  |  |
| `CLICKUP_CREATE_CHAT_VIEW_COMMENT` | Create chat view comment | write | `comment_text`, `notify_all`, `view_id` |  |  |
| `CLICKUP_CREATE_CHECKLIST` | Create checklist | write | `name`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_CREATE_CHECKLIST_ITEM` | Create checklist item | write | `checklist_id`, `name`; *assignee* |  |  |
| `CLICKUP_CREATE_DIRECT_MESSAGE_CHANNEL` | Create direct message channel | write | `workspace_id`; *user_ids* |  |  |
| `CLICKUP_CREATE_DOC` | Create Doc | write | `name`, `workspace_id`; *parent_id, parent_type* |  |  |
| `CLICKUP_CREATE_DOC_PAGE` | Create Doc Page | write | `doc_id`, `workspace_id`; *content, content_format, name, parent_page_id, sub_title* |  |  |
| `CLICKUP_CREATE_FOLDER` | Create folder | write | `name`, `space_id` |  |  |
| `CLICKUP_CREATE_FOLDERLESS_LIST` | Create folderless list | write | `name`, `space_id`; *assignee, content, due_date, due_date_time, priority, status* |  |  |
| `CLICKUP_CREATE_FOLDER_VIEW` | Create folder view | write | `folder_id`, `name`, `type`; *columns__fields, divide__collapsed, divide__dir, divide__field, filters__fields, filters__op, filters__search, filters__show__closed and 18 more* |  |  |
| `CLICKUP_CREATE_GOAL` | Create goal | write | `color`, `description`, `due_date`, `multiple_owners`, `name`, `owners`, `team_id` |  |  |
| `CLICKUP_CREATE_KEY_RESULT` | Create key result | write | `goal_id`, `list_ids`, `name`, `owners`, `steps_end`, `steps_start`, `task_ids`, `type`, `unit` |  |  |
| `CLICKUP_CREATE_LIST` | Create list | write | `folder_id`, `name`; *assignee, content, due_date, due_date_time, priority, status* |  |  |
| `CLICKUP_CREATE_LIST_COMMENT` | Create list comment | write | `comment_text`, `list_id`; *assignee, notify_all* |  |  |
| `CLICKUP_CREATE_LIST_FROM_TEMPLATE` | Create list from template in folder | write | `folder_id`, `name`, `template_id` |  |  |
| `CLICKUP_CREATE_LIST_VIEW` | Create list view | write | `list_id`, `name`, `type`; *columns__fields, divide__collapsed, divide__dir, divide__field, filters__fields, filters__op, filters__search, filters__show__closed and 18 more* |  |  |
| `CLICKUP_CREATE_LOCATION_CHAT_CHANNEL` | Create Location Chat Channel | write | `location`, `workspace_id`; *description, topic, user_ids, visibility* |  |  |
| `CLICKUP_CREATE_REPLY_MESSAGE` | Create reply message | write | `content`, `message_id`, `type`, `workspace_id`; *assignee, content_format, followers, group_assignee, post_data, reactions, triaged_action, triaged_object_id and 1 more* |  |  |
| `CLICKUP_CREATE_SPACE` | Create space | write | `multiple_assignees`, `name`, `team_id`; *features__checklists__enabled, features__custom__fields__enabled, features__dependency__warning__enabled, features__due__dates__enabled, features__due__dates__remap__closed__due__date, features__due__dates__remap__due__dates, features__due__dates__start__date, features__portfolios__enabled and 4 more* |  |  |
| `CLICKUP_CREATE_SPACE_TAG` | Create space tag | write | `space_id`, `tag__name`, `tag__tag_bg`, `tag__tag_fg` |  |  |
| `CLICKUP_CREATE_SPACE_VIEW` | Create space view | write | `name`, `space_id`, `type`; *columns__fields, divide__collapsed, divide__dir, divide__field, filters__fields, filters__op, filters__search, filters__show__closed and 18 more* |  |  |
| `CLICKUP_CREATE_TASK` | Create a task | write | `list_id`, `name`; *assignees, check_required_custom_fields, custom_fields, custom_item_id, custom_task_ids, description, due_date, due_date_time and 10 more* |  |  |
| `CLICKUP_CREATE_TASK_ATTACHMENT` | Create task attachment | write | `attachment`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_CREATE_TASK_COMMENT` | Create task comment | write | `assignee`, `comment_text`, `notify_all`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_CREATE_TASK_FROM_TEMPLATE` | Create task from template | write | `list_id`, `name`, `template_id` |  |  |
| `CLICKUP_CREATE_TEAM` | Create team | write | `members`, `name`, `team_id` |  |  |
| `CLICKUP_CREATE_THREADED_COMMENT` | Create threaded comment | write | `comment_id`, `comment_text`, `notify_all`; *assignee, group_assignee* |  |  |
| `CLICKUP_CREATE_WEBHOOK` | Create webhook | write | `endpoint`, `events`, `team_id`; *folder_id, list_id, space_id, task_id* |  |  |
| `CLICKUP_CREATE_WORKSPACE_EVERYTHING_LEVEL_VIEW` | Create workspace Everything level view | write | `name`, `team_id`, `type`; *columns__fields, divide__collapsed, divide__dir, divide__field, filters__fields, filters__op, filters__search, filters__show__closed and 18 more* |  |  |
| `CLICKUP_DELETE_CHAT_CHANNEL` | Delete chat channel | write | `channel_id`, `workspace_id` |  |  |
| `CLICKUP_DELETE_CHAT_MESSAGE` | Delete chat message | write | `message_id`, `workspace_id` |  |  |
| `CLICKUP_DELETE_CHAT_REACTION` | Delete chat reaction | write | `message_id`, `reaction`, `workspace_id` |  |  |
| `CLICKUP_DELETE_CHECKLIST` | Delete checklist | write | `checklist_id` |  |  |
| `CLICKUP_DELETE_CHECKLIST_ITEM` | Delete checklist item | write | `checklist_id`, `checklist_item_id` |  |  |
| `CLICKUP_DELETE_COMMENT` | Delete comment | write | `comment_id` |  |  |
| `CLICKUP_DELETE_DEPENDENCY` | Delete dependency | write | `task_id`; *custom_task_ids, dependency_of, depends_on, team_id* |  |  |
| `CLICKUP_DELETE_FOLDER` | Delete folder | write | `folder_id` |  |  |
| `CLICKUP_DELETE_GOAL` | Delete goal | write | `goal_id` |  |  |
| `CLICKUP_DELETE_KEY_RESULT` | Delete key result | write | `key_result_id` |  |  |
| `CLICKUP_DELETE_LIST` | Delete list | write | `list_id` |  |  |
| `CLICKUP_DELETE_SPACE` | Delete space | write | `space_id` |  |  |
| `CLICKUP_DELETE_SPACE_TAG` | Delete space tag | write | `space_id`, `tag_name`; *tag__name, tag__tag__bg, tag__tag__fg* |  |  |
| `CLICKUP_DELETE_TASK` | Delete task | write | `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_DELETE_TASK_LINK` | Delete task link | write | `links_to`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_DELETE_TEAM` | Delete team | write | `group_id` |  |  |
| `CLICKUP_DELETE_TIME_ENTRY` | Delete a time entry | write | `team_id`, `timer_id` |  |  |
| `CLICKUP_DELETE_TIME_TRACKED` | Delete time tracked | write | `interval_id`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_DELETE_VIEW` | Delete view | write | `view_id` |  |  |
| `CLICKUP_DELETE_WEBHOOK` | Delete webhook | write | `webhook_id` |  |  |
| `CLICKUP_GET_ACCESSIBLE_CUSTOM_FIELDS` | Get accessible custom fields | write | `list_id` |  |  |
| `CLICKUP_GET_ACCESS_TOKEN` | Get access token | write | `client_id`, `client_secret`, `code` |  |  |
| `CLICKUP_GET_ALL_TAGS_FROM_TIME_ENTRIES` | Get all tags from time entries | read | `team_id` |  |  |
| `CLICKUP_GET_AUTHORIZED_TEAMS_WORKSPACES` | Get authorized teams workspaces | read |  |  |  |
| `CLICKUP_GET_AUTHORIZED_USER` | Get authorized user | read |  |  |  |
| `CLICKUP_GET_BULK_TASKS_TIME_IN_STATUS` | Get bulk tasks time in status | read | `task_ids`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_GET_CHAT_CHANNEL` | Get chat channel | read | `channel_id`, `workspace_id`; *description_format* |  |  |
| `CLICKUP_GET_CHAT_CHANNELS` | Get chat channels | read | `workspace_id`; *channel_types, cursor, description_format, include_closed, is_follower, limit, with_message_since* |  |  |
| `CLICKUP_GET_CHAT_CHANNEL_FOLLOWERS` | Get chat channel followers | read | `channel_id`, `workspace_id`; *cursor, limit* |  |  |
| `CLICKUP_GET_CHAT_CHANNEL_MEMBERS` | Get chat channel members | read | `channel_id`, `workspace_id`; *cursor, limit* |  |  |
| `CLICKUP_GET_CHAT_MESSAGES` | Get chat messages | read | `channel_id`, `workspace_id`; *content_format, cursor, limit* |  |  |
| `CLICKUP_GET_CHAT_MESSAGE_REACTIONS` | Get chat message reactions | read | `message_id`, `workspace_id`; *cursor, limit* |  |  |
| `CLICKUP_GET_CHAT_MESSAGE_REPLIES` | Get chat message replies | read | `message_id`, `workspace_id`; *content_format, cursor, limit* |  |  |
| `CLICKUP_GET_CHAT_MESSAGE_TAGGED_USERS` | Get chat message tagged users | read | `message_id`, `workspace_id`; *cursor, limit* |  |  |
| `CLICKUP_GET_CHAT_VIEW_COMMENTS` | Get chat view comments | read | `view_id`; *start, start_id* |  |  |
| `CLICKUP_GET_CUSTOM_ROLES` | Get custom roles | write | `team_id`; *include_members* |  |  |
| `CLICKUP_GET_CUSTOM_TASK_TYPES` | Get custom task types | write | `team_id` |  |  |
| `CLICKUP_GET_DOC` | Get Doc | read | `doc_id`, `workspace_id` |  |  |
| `CLICKUP_GET_DOC_CONTENT` | Get Doc Content | read | `doc_id`, `workspace_id`; *include_page_listing_only, page_ids* |  |  |
| `CLICKUP_GET_DOC_PAGES_PUBLIC` | Get Doc Pages | read | `doc_id`, `workspace_id`; *content_format, max_page_depth* |  |  |
| `CLICKUP_GET_DOC_PAGE_CONTENT` | Get Doc Page Content | read | `doc_id`, `page_id`, `workspace_id` |  |  |
| `CLICKUP_GET_DOC_PAGE_LISTING` | Get Doc Page Listing | read | `doc_id`, `workspace_id`; *max_page_depth* |  |  |
| `CLICKUP_GET_FILTERED_TEAM_TASKS` | Get filtered team tasks | read | `team_Id`; *assignees, custom_fields, custom_items, custom_task_ids, date_created_gt, date_created_lt, date_done_gt, date_done_lt and 17 more* |  |  |
| `CLICKUP_GET_FOLDER` | Get folder | read | `folder_id` |  |  |
| `CLICKUP_GET_FOLDERLESS_LISTS` | Get folderless lists | read | `space_id`; *archived* |  |  |
| `CLICKUP_GET_FOLDERS` | Get folders | read | `space_id`; *archived* |  |  |
| `CLICKUP_GET_FOLDER_AVAILABLE_FIELDS` | Get Folder Available Fields | read | `folder_id` |  |  |
| `CLICKUP_GET_FOLDER_VIEWS` | Get folder views | read | `folder_id` |  |  |
| `CLICKUP_GET_GOAL` | Get goal | read | `goal_id` |  |  |
| `CLICKUP_GET_GOALS` | Get goals | read | `team_id`; *include_completed* |  |  |
| `CLICKUP_GET_GUEST` | Get guest | read | `guest_id`, `team_id` |  |  |
| `CLICKUP_GET_LIST` | Get list | read | `list_id` |  |  |
| `CLICKUP_GET_LISTS` | Get lists | read | `folder_id`; *archived* |  |  |
| `CLICKUP_GET_LIST_COMMENTS` | Get list comments | read | `list_id`; *start, start_id* |  |  |
| `CLICKUP_GET_LIST_MEMBERS` | Get list members | read | `list_id` |  |  |
| `CLICKUP_GET_LIST_VIEWS` | Get list views | read | `list_id` |  |  |
| `CLICKUP_GET_RUNNING_TIME_ENTRY` | Get running time entry | read | `team_id`; *assignee* |  |  |
| `CLICKUP_GET_SHARED_HIERARCHY` | Shared hierarchy | read | `team_id` |  |  |
| `CLICKUP_GET_SPACE` | Get space | read | `space_id` |  |  |
| `CLICKUP_GET_SPACES` | Get spaces | read | `team_id`; *archived* |  |  |
| `CLICKUP_GET_SPACE_AVAILABLE_FIELDS` | Get space available fields | read | `space_id` |  |  |
| `CLICKUP_GET_SPACE_TAGS` | Get space tags | read | `space_id` |  |  |
| `CLICKUP_GET_SPACE_VIEWS` | Get space views | read | `space_id` |  |  |
| `CLICKUP_GET_SUBTYPES` | Get post subtype IDs | read | `comment_type`, `workspace_id` |  |  |
| `CLICKUP_GET_TASK` | Get task | read | `task_id`; *custom_task_ids, include_markdown_description, include_subtasks, team_id* |  |  |
| `CLICKUP_GET_TASKS` | Get tasks | read | `list_id`; *archived, assignees, custom_fields, custom_items, date_created_gt, date_created_lt, date_done_gt, date_done_lt and 12 more* |  |  |
| `CLICKUP_GET_TASK_COMMENTS` | Get task comments | read | `task_id`; *custom_task_ids, start, start_id, team_id* |  |  |
| `CLICKUP_GET_TASK_MEMBERS` | Get task members | read | `task_id` |  |  |
| `CLICKUP_GET_TASK_TEMPLATES` | Get task templates | read | `page`, `team_id` |  |  |
| `CLICKUP_GET_TASK_TIME_IN_STATUS` | Get task s time in status | read | `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_GET_TEAMS` | Get teams | read | `team_id`; *group_ids* |  |  |
| `CLICKUP_GET_TEAM_AVAILABLE_FIELDS` | Get team available fields | read | `team_id` |  |  |
| `CLICKUP_GET_THREADED_COMMENTS` | Get threaded comments | read | `comment_id` |  |  |
| `CLICKUP_GET_TIME_ENTRIES_IN_DATE_RANGE` | Get time entries within a date range | read | `team_Id`; *assignee, custom_task_ids, end_date, folder_id, include_location_names, include_task_tags, list_id, space_id and 3 more* |  |  |
| `CLICKUP_GET_TIME_ENTRY` | Get singular time entry | read | `team_id`, `timer_id`; *include__task, include_location_names* |  |  |
| `CLICKUP_GET_TIME_ENTRY_HISTORY` | Get time entry history | read | `team_id`, `timer_id` |  |  |
| `CLICKUP_GET_TRACKED_TIME` | Get tracked time | read | `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_GET_USER` | Get user | read | `team_id`, `user_id` |  |  |
| `CLICKUP_GET_VIEW` | Get view | read | `view_id` |  |  |
| `CLICKUP_GET_VIEW_TASKS` | Get view tasks | read | `page`, `view_id` |  |  |
| `CLICKUP_GET_WEBHOOKS` | Get webhooks | read | `team_id` |  |  |
| `CLICKUP_GET_WORKSPACE_EVERYTHING_LEVEL_VIEWS` | Get workspace everything level views | read | `team_id`; *archived* |  |  |
| `CLICKUP_GET_WORKSPACE_PLAN` | Get workspace plan | read | `team_id` |  |  |
| `CLICKUP_GET_WORKSPACE_SEATS` | Get workspace seats | read | `team_id` |  |  |
| `CLICKUP_INVITE_GUEST_TO_WORKSPACE` | Invite guest to workspace | write | `can_create_views`, `can_edit_tags`, `can_see_time_estimated`, `can_see_time_spent`, `email`, `team_id`; *custom_role_id* |  |  |
| `CLICKUP_INVITE_USER_TO_WORKSPACE` | Invite user to workspace | write | `admin`, `email`, `team_id`; *custom_role_id* |  |  |
| `CLICKUP_MOVE_TASK_TO_HOME_LIST` | Move task to new home list | write | `list_id`, `task_id`, `workspace_id`; *custom_fields_to_move, move_custom_fields, status_mappings* |  |  |
| `CLICKUP_REMOVE_CUSTOM_FIELD_VALUE` | Remove custom field value | write | `field_id`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_REMOVE_GUEST_FROM_FOLDER` | Remove guest from folder | write | `folder_id`, `guest_id`; *include_shared* |  |  |
| `CLICKUP_REMOVE_GUEST_FROM_LIST` | Remove guest from list | write | `guest_id`, `list_id`; *include_shared* |  |  |
| `CLICKUP_REMOVE_GUEST_FROM_TASK` | Remove guest from task | write | `guest_id`, `task_id`; *custom_task_ids, include_shared, team_id* |  |  |
| `CLICKUP_REMOVE_GUEST_FROM_WORKSPACE` | Remove guest from workspace | write | `guest_id`, `team_id` |  |  |
| `CLICKUP_REMOVE_TAGS_FROM_TIME_ENTRIES` | Remove tags from time entries | write | `tags`, `team_id`, `time_entry_ids` |  |  |
| `CLICKUP_REMOVE_TAG_FROM_TASK` | Remove tag from task | write | `tag_name`, `task_id`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_REMOVE_TASK_FROM_LIST` | Remove task from list | write | `list_id`, `task_id` |  |  |
| `CLICKUP_REMOVE_USER_FROM_WORKSPACE` | Remove user from workspace | write | `team_id`, `user_id` |  |  |
| `CLICKUP_SEARCH_DOCS` | Search Docs | read | `workspace_id`; *cursor, limit* |  |  |
| `CLICKUP_SET_CUSTOM_FIELD_VALUE` | Set custom field value | write | `field_id`, `task_id`, `value`; *custom_task_ids, team_id, value_options* |  |  |
| `CLICKUP_START_TIME_ENTRY` | Start a time entry | write | `team_Id`; *billable, custom_task_ids, description, tags, team_id, tid* |  |  |
| `CLICKUP_STOP_TIME_ENTRY` | Stop a time entry | write | `team_id` |  |  |
| `CLICKUP_TRACK_TIME` | Track time | write | `end`, `start`, `task_id`, `time`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_UPDATE_CHAT_CHANNEL` | Update chat channel | write | `channel_id`, `workspace_id`; *content_format, description, location, name, topic, visibility* |  |  |
| `CLICKUP_UPDATE_CHAT_MESSAGE` | Update chat message | write | `message_id`, `workspace_id`; *assignee, content, content_format, group_assignee, post_data, resolved* |  |  |
| `CLICKUP_UPDATE_CHECKLIST` | Edit checklist | write | `checklist_id`; *name, position* |  |  |
| `CLICKUP_UPDATE_CHECKLIST_ITEM` | Edit checklist item | write | `checklist_id`, `checklist_item_id`; *assignee, name, parent, resolved* |  |  |
| `CLICKUP_UPDATE_COMMENT` | Update comment | write | `assignee`, `comment_id`, `comment_text`, `resolved` |  |  |
| `CLICKUP_UPDATE_FOLDER` | Update folder | write | `folder_id`, `name` |  |  |
| `CLICKUP_UPDATE_GOAL` | Update goal | write | `add_owners`, `color`, `description`, `due_date`, `goal_id`, `name`, `rem_owners` |  |  |
| `CLICKUP_UPDATE_GUEST_ON_WORKSPACE` | Edit guest on workspace | write | `guest_id`, `team_id`; *can_create_views, can_edit_tags, can_see_time_estimated, can_see_time_spent, custom_role_id, username* |  |  |
| `CLICKUP_UPDATE_KEY_RESULT` | Edit key result | write | `key_result_id`, `note`, `steps_current` |  |  |
| `CLICKUP_UPDATE_LIST` | Update list | write | `list_id`; *assignee, content, due_date, due_date_time, name, priority, status, unset_status* |  |  |
| `CLICKUP_UPDATE_SPACE` | Update space | write | `color`, `multiple_assignees`, `name`, `private`, `space_id`; *admin_can_manage, features__checklists__enabled, features__custom_fields__enabled, features__dependency_warning__enabled, features__due_dates__enabled, features__due_dates__remap_closed_due_date, features__due_dates__remap_due_dates, features__due_dates__start_date and 5 more* |  |  |
| `CLICKUP_UPDATE_SPACE_TAG` | Edit space tag | write | `space_id`, `tag_name`; *tag__bg__color, tag__fg__color, tag__name* |  |  |
| `CLICKUP_UPDATE_TASK` | Update task | write | `task_id`; *archived, assignees__add, assignees__rem, custom_item_id, custom_task_ids, description, due_date, due_date_time and 8 more* |  |  |
| `CLICKUP_UPDATE_TEAM` | Update team | write | `group_id`; *handle, members__add, members__rem, name* |  |  |
| `CLICKUP_UPDATE_TIME_ENTRY` | Update a time entry | write | `team_id`, `timer_id`; *billable, custom_task_ids, description, duration, end, start, tag_action, tags and 1 more* |  |  |
| `CLICKUP_UPDATE_TIME_ENTRY_TAG` | Change tag names from time entries | write | `name`, `new_name`, `tag_bg`, `tag_fg`, `team_id` |  |  |
| `CLICKUP_UPDATE_TIME_TRACKED` | Edit time tracked | write | `end`, `interval_id`, `start`, `task_id`, `time`; *custom_task_ids, team_id* |  |  |
| `CLICKUP_UPDATE_USER_ON_WORKSPACE` | Edit user on workspace | write | `team_id`, `user_id`; *admin, custom_role_id, username* |  |  |
| `CLICKUP_UPDATE_VIEW` | Update view | write | `name`, `type`, `view_id`; *columns__fields, divide__collapsed, divide__dir, divide__field, filters__fields, filters__op, filters__search, filters__show__closed and 20 more* |  |  |
| `CLICKUP_UPDATE_WEBHOOK` | Update webhook | write | `endpoint`, `events`, `status`, `webhook_id` |  |  |
| `CLICKUP_UPDATE_WORKSPACE_ACL` | Update workspace ACL | write | `object_id`, `object_type`, `workspace_id`; *entries, private* |  |  |
