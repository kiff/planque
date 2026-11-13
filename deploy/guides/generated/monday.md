# Monday: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

122 tools: 43 read, 79 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `MONDAY_ACCOUNT_TRIGGER_STATISTICS` | Get account trigger statistics | read |  |  | `account:read` |
| `MONDAY_ACTIVATE_USERS` | Activate users | write | `user_ids` |  | `users:write` |
| `MONDAY_ADD_SUBSCRIBERS_TO_OBJECT` | Add subscribers to object | write | `id`, `user_ids`; *kind* |  | `boards:write` |
| `MONDAY_ADD_TEAMS_TO_BOARD` | Add teams to board | write | `board_id`, `team_ids`; *kind* |  | `boards:write` |
| `MONDAY_ADD_USERS_TO_BOARD` | Add users to board | write | `board_id`, `user_ids`; *kind* |  | `boards:write` |
| `MONDAY_ADD_USERS_TO_TEAM` | Add users to team | write | `team_id`, `user_ids` |  | `teams:write` |
| `MONDAY_ADD_USERS_TO_WORKSPACE` | Add users to workspace | write | `user_ids`, `workspace_id`; *kind* |  | `workspaces:write` |
| `MONDAY_AGGREGATE_DATA` | Aggregate board data | read | `query` |  | `boards:read` |
| `MONDAY_APP_SUBSCRIPTION` | Get app subscription | read |  |  | `account:read` |
| `MONDAY_ARCHIVE_BOARD` | Archive board | write | `board_id` |  | `boards:write` |
| `MONDAY_ARCHIVE_ITEM` | Archive item | write | `item_id` |  | `boards:write` |
| `MONDAY_ARCHIVE_OBJECT` | Archive object | write | `id` |  | `boards:write` |
| `MONDAY_ARCHIVE_WORKSPACE` | Archive workspace | write | `workspace_id` |  | `workspaces:write` |
| `MONDAY_AUDIT_LOGS` | Get audit logs | read | *end_time, events, ip_address, limit, page, start_time, user_id* |  | `manage_account_security` |
| `MONDAY_BLOCKS` | Get document blocks | read | `doc_ids`; *limit, page* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_BOARDS` | Get boards | read | *board_kind, hierarchy_type, ids, limit, order_by, page, state, workspace_ids* |  | `boards:read`, `boards:write` |
| `MONDAY_CHANGE_SIMPLE_COLUMN_VALUE` | Change simple column value | write | `board_id`, `column_id`, `item_id`, `value`; *create_labels_if_missing* |  | `boards:write` |
| `MONDAY_COLUMNS` | Get board columns | read | `board_ids`; *column_ids, column_types* |  | `boards:read`, `boards:write` |
| `MONDAY_CONNECTIONS` | Get connections | read | *connection_state, order, page, page_size, pagination, with_automations, with_partial_scopes, with_state_validation* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_CONNECTION_BOARD_IDS` | Get connection board IDs | read | *board_ids, limit* |  | `boards:read` |
| `MONDAY_CONNECT_PROJECT_TO_PORTFOLIO` | Connect project to portfolio | write | `portfolio_board_id`, `project_board_id` |  | `boards:write` |
| `MONDAY_CONVERT_BOARD_TO_PROJECT` | Convert board to project | write | `board_id`, `column_mappings` |  | `boards:write` |
| `MONDAY_CREATE_BOARD` | Create a Monday board | write | `board_kind`, `board_name`; *board_owner_ids, board_subscriber_ids, board_subscriber_teams_ids, description, folder_id, template_id, workspace_id* |  | `boards:write` |
| `MONDAY_CREATE_COLUMN` | Create column | write | `board_id`, `column_type`, `title`; *after_column_id, defaults, description* |  | `boards:write` |
| `MONDAY_CREATE_CUSTOM_ACTIVITY` | Create custom activity | write | `color`, `icon_id`, `name` |  | `boards:write` |
| `MONDAY_CREATE_DASHBOARD` | Create a Monday dashboard | write | `board_ids`, `name`, `workspace_id`; *board_folder_id, kind* |  | `boards:read`, `boards:write`, `workspaces:write` |
| `MONDAY_CREATE_DOC` | Create doc | write | *column_id, item_id, kind, name, workspace_id* |  | `boards:read`, `boards:write`, `items:read`, `items:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `workspaces:read`, `workspaces:write`, `teams:read`, `teams:write`, `tags:read`, `tags:write`, `files:read`, `files:write`, `account:read`, `account:write`, `webhooks:read`, `webhooks:write`, `notifications:read`, `notifications:write`, `settings:read`, `settings:write`, `templates:read`, `templates:write`, `activity_logs:read`, `activity_logs:write`, `audit_logs:read`, `audit_logs:write`, `time_tracking:read`, `time_tracking:write`, `forms:read`, `forms:write`, `docs:read`, `docs:write` |
| `MONDAY_CREATE_FOLDER` | Create a folder | write | `name`, `workspace_id`; *color, parent_folder_id* |  | `workspaces:write` |
| `MONDAY_CREATE_GROUP` | Create group | write | `board_id`, `group_name` |  | `boards:read`, `boards:write`, `users:read`, `users:write`, `teams:read`, `teams:write`, `workspaces:read`, `workspaces:write`, `docs:read`, `docs:write`, `updates:read`, `updates:write`, `notifications:write`, `tags:read`, `webhooks:read`, `webhooks:write`, `account:read`, `assets:read`, `me:read` |
| `MONDAY_CREATE_ITEM` | Create item | write | `board_id`, `item_name`; *column_values, group_id* |  | `boards:write` |
| `MONDAY_CREATE_ITEM_FROM_NL` | Create Item From Natural Language | write | `board_id`, `nl_query`; *group_id* |  | `boards:write` |
| `MONDAY_CREATE_NOTIFICATION` | Create notification | write | `target_id`, `target_type`, `text`, `user_id` |  | `notifications:write` |
| `MONDAY_CREATE_OBJECT` | Create object | write | `query`; *variables* |  | `boards:write`, `docs:write`, `workspaces:write`, `updates:write`, `notifications:write`, `teams:write`, `webhooks:write` |
| `MONDAY_CREATE_TAG` | Create tag | write | `tag_name`; *board_id* |  | `boards:write` |
| `MONDAY_CREATE_TIMELINE_ITEM` | Create timeline item | write | `custom_activity_id`, `item_id`, `timestamp`, `title`; *content, location, phone, summary, time_range, url* |  | `boards:write` |
| `MONDAY_CREATE_UPDATE` | Create update | write | `body`; *item_id, mentions_list, original_creation_date, parent_id, use_app_info* |  | `boards:read`, `boards:write`, `users:read`, `users:write`, `updates:read`, `updates:write`, `workspaces:read`, `workspaces:write`, `docs:read`, `docs:write`, `teams:read`, `teams:write`, `tags:read`, `notifications:write`, `webhooks:read`, `webhooks:write`, `account:read`, `me:read` |
| `MONDAY_CREATE_USER` | Invite and create a Monday user | write | `emails`; *product, user_role* |  | `users:write` |
| `MONDAY_CREATE_WORKSPACE` | Create a Monday workspace | write | `kind`, `name`; *account_product_id, description* |  | `workspaces:write` |
| `MONDAY_CUSTOM_ACTIVITY` | Get custom activity | write |  |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_DEACTIVATE_USERS` | Deactivate users | write | `user_ids` |  | `users:write` |
| `MONDAY_DELETE_ASSET` | Delete asset | write | `asset_id`; *board_id, column_id, item_id, target, update_id* |  | `assets:read`, `assets:write`, `updates:write`, `boards:write` |
| `MONDAY_DELETE_BOARD` | Delete board | write | `board_id` |  | `boards:write` |
| `MONDAY_DELETE_COLUMN` | Delete column | write | `board_id`, `column_id` |  |  |
| `MONDAY_DELETE_CUSTOM_ACTIVITY` | Delete Custom Activity | write | `id` |  | `boards:write` |
| `MONDAY_DELETE_DOC` | Delete doc | write | `doc_id` |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_DELETE_FOLDER` | Delete folder | write | `folder_id` |  | `workspaces:write` |
| `MONDAY_DELETE_GROUP` | Delete group | write | `board_id`, `group_id` |  | `boards:write` |
| `MONDAY_DELETE_ITEM` | Delete item | write | `item_id` |  | `boards:write` |
| `MONDAY_DELETE_OBJECT` | Delete object | write | `id` |  | `boards:write`, `workspaces:write`, `docs:write`, `updates:write`, `webhooks:write`, `teams:write` |
| `MONDAY_DELETE_SUBSCRIBERS_FROM_BOARD` | Delete subscribers from board | write | `board_id`, `user_ids` |  | `boards:write` |
| `MONDAY_DELETE_TAG` | Delete tag | write | `tag_id`; *board_id, item_id* |  | `boards:read`, `boards:write`, `tags:write` |
| `MONDAY_DELETE_TEAM` | Delete team | write | `team_id` |  | `teams:write` |
| `MONDAY_DELETE_TEAMS_FROM_BOARD` | Delete teams from board | write | `board_id`, `team_ids` |  | `boards:write` |
| `MONDAY_DELETE_TEAMS_FROM_WORKSPACE` | Delete teams from workspace | write | `team_ids`, `workspace_id` |  | `workspaces:write` |
| `MONDAY_DELETE_TIMELINE_ITEM` | Delete timeline item | write | `id` |  | `boards:write` |
| `MONDAY_DELETE_UPDATE` | Delete update | write | `update_id` |  | `updates:write` |
| `MONDAY_DELETE_WORKSPACE` | Delete workspace | write | `workspace_id` |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_DOCS` | Retrieve Monday docs | read | *ids, limit, object_ids, order_by, page, workspace_ids* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_DUPLICATE_BOARD` | Duplicate board | write | `board_id`, `duplicate_type`; *board_name, folder_id, keep_subscribers, workspace_id* |  | `boards:write` |
| `MONDAY_DUPLICATE_ITEM` | Duplicate item | write | `board_id`, `item_id`; *with_updates* |  | `boards:write` |
| `MONDAY_EDIT_UPDATE` | Edit update | write | `body`, `id` |  | `boards:write`, `updates:write` |
| `MONDAY_GET_ACCOUNT_INFO` | Get account info | read |  |  | `account:read` |
| `MONDAY_GET_ACTIVITY_LOGS` | Get board activity logs | read | `board_id`; *column_ids, from, group_ids, item_ids, limit, page, to, user_ids* |  | `boards:read` |
| `MONDAY_GET_API_VERSION` | Get API version | read |  |  |  |
| `MONDAY_GET_ASSETS` | Get assets | read | `ids` |  | `assets:read` |
| `MONDAY_GET_BOARD_VIEWS` | Get board views | read | `board_ids`; *view_type* |  | `boards:read` |
| `MONDAY_GET_FAVORITES` | Get favorites | read |  |  | `users:read` |
| `MONDAY_GET_FOLDERS` | Get folders | read | *ids, limit, page, workspace_ids* |  | `workspaces:read` |
| `MONDAY_GET_FORM` | Get form details | read | `form_token` |  | `boards:read` |
| `MONDAY_GET_ITEMS` | Get items by IDs | read | `ids`; *exclude_nonactive, limit, newest_first, page* |  | `boards:read` |
| `MONDAY_GET_ME` | Get current user | read | *api_version* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_GET_MUTATION_COMPLEXITY` | Execute GraphQL and get mutation complexity | write | `query` |  |  |
| `MONDAY_GET_QUERY_COMPLEXITY` | Get query complexity | read |  |  |  |
| `MONDAY_GET_TAGS` | Get tags | read | *ids* |  | `boards:read`, `boards:write`, `users:read`, `users:write`, `workspaces:read`, `workspaces:write`, `updates:read`, `updates:write`, `tags:read`, `tags:write`, `teams:read`, `teams:write`, `docs:read`, `docs:write`, `account:read`, `account:write`, `me:read`, `me:write`, `notifications:read`, `notifications:write`, `webhooks:read`, `webhooks:write`, `assets:read`, `assets:write` |
| `MONDAY_GET_TEAMS` | Get teams | read | *ids* |  | `teams:read`, `teams:write`, `boards:write`, `workspaces:write` |
| `MONDAY_GET_UPDATES` | Get item updates | read | `item_id`; *limit, page* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_GET_UPDATE_REPLIES` | Get update replies | read | `item_id`; *limit, page* |  | `boards:read` |
| `MONDAY_GET_VERSIONS` | Get API versions | read |  |  |  |
| `MONDAY_GET_VIEW_SCHEMA_BY_TYPE` | Get view schema by type | read | `mutation_type`, `type` |  | `boards:read` |
| `MONDAY_GET_WEBHOOKS` | Get webhooks | read | `board_id`; *app_webhooks_only* |  | `boards:read` |
| `MONDAY_GET_WORKSPACES` | Get workspaces | read | *ids, kind, limit, order_by, page, state* |  | `workspaces:read` |
| `MONDAY_IMPORT_DOC_FROM_HTML` | Import doc from HTML | write | `html`, `workspace_id`; *folder_id, kind, title* |  | `docs:write` |
| `MONDAY_ITEMS_PAGE` | Get items page | read | *board_id, cursor, group_id, include_column_values, limit, query_params* |  | `boards:read` |
| `MONDAY_LIKE_UPDATE` | Like update | write | `update_id` |  | `updates:write` |
| `MONDAY_LIST_GROUPS` | List groups | read | `board_id`; *group_ids* |  | `boards:read`, `boards:write`, `me:read`, `users:read`, `users:write`, `teams:read`, `teams:write`, `workspaces:read`, `workspaces:write`, `updates:read`, `updates:write`, `notifications:write`, `tags:read`, `docs:read`, `docs:write`, `account:read`, `assets:read`, `webhooks:read`, `webhooks:write` |
| `MONDAY_LIST_ITEMS` | List items | read | `item_ids` |  | `boards:read` |
| `MONDAY_LIST_ITEMS_BY_COLUMN_VALUES` | List items by column values | read | `board_id`; *columns, cursor, limit* |  | `boards:read` |
| `MONDAY_LIST_SUBITEMS_BY_PARENT` | List subitems by parent | read | `parent_item_ids`; *include_column_values, include_parent_fields* |  | `boards:read` |
| `MONDAY_LIST_TEAM_MEMBERS` | List team members | read | *team_ids, user_filters* |  | `users:read`, `teams:read` |
| `MONDAY_LIST_USERS` | List users | read | *limit, page* |  | `boards:write` |
| `MONDAY_MOVE_ITEM_TO_BOARD` | Move item to board | write | `board_id`, `group_id`, `item_id` |  | `boards:write` |
| `MONDAY_MOVE_ITEM_TO_GROUP` | Move item to group | write | `group_id`, `item_id` |  | `boards:write` |
| `MONDAY_MUTE_BOARD_SETTINGS` | Query board mute settings | read | `board_ids` |  | `boards:write` |
| `MONDAY_PIN_TO_TOP` | Pin update to top | write | `update_id`; *item_id* |  | `updates:write` |
| `MONDAY_PUBLISH_OBJECT` | Publish object | write | `id` |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_QUERY_DASHBOARDS` | Query dashboards | read | *limit* |  | `boards:read` |
| `MONDAY_QUERY_SPRINTS` | Query sprints | read | `board_ids`; *sprint_ids* |  | `boards:read` |
| `MONDAY_REMOVE_USERS_FROM_TEAM` | Remove users from team | write | `team_id`, `user_ids` |  | `teams:write` |
| `MONDAY_REMOVE_USERS_FROM_WORKSPACE` | Remove users from workspace | write | `user_ids`, `workspace_id` |  | `workspaces:write` |
| `MONDAY_SET_BOARD_PERMISSION` | Set board permission | write | `basic_role_name`, `board_id` |  | `boards:write` |
| `MONDAY_TIMELINE` | Get item timeline | read | `id` |  | `boards:read`, `boards:write` |
| `MONDAY_TIMELINE_ITEM` | Get timeline item | read | `id` |  | `boards:read`, `boards:write` |
| `MONDAY_UNLIKE_UPDATE` | Unlike update | write | `update_id` |  | `updates:write` |
| `MONDAY_UNPIN_FROM_TOP` | Unpin update from top | write | `update_id`; *item_id* |  | `updates:write` |
| `MONDAY_UNPUBLISH_OBJECT` | Unpublish object | write | `id` |  | `boards:write` |
| `MONDAY_UPDATE_BOARD` | Update board | write | `board_attribute`, `board_id`, `new_value` |  | `boards:write` |
| `MONDAY_UPDATE_BOARD_HIERARCHY` | Update board hierarchy | write | `attributes`, `board_id` |  | `boards:write` |
| `MONDAY_UPDATE_COLUMN` | Update column | write | `board_id`, `column_id`, `column_property`, `value` |  | `boards:write` |
| `MONDAY_UPDATE_DOC` | Update a doc | write | `doc_id`; *after_block_id, markdown, title* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_UPDATE_EMAIL_DOMAIN` | Update email domain | write | `new_domain`, `user_ids` |  | `users:write` |
| `MONDAY_UPDATE_FOLDER` | Update a Monday folder | write | `folder_id`; *color, name, parent_folder_id, workspace_id* |  | `workspaces:write` |
| `MONDAY_UPDATE_GROUP` | Update group | write | `board_id`, `group_attribute`, `group_id`, `new_value` |  | `boards:read`, `boards:write`, `users:read`, `users:write`, `teams:read`, `teams:write`, `workspaces:read`, `workspaces:write`, `docs:read`, `docs:write`, `updates:read`, `updates:write`, `notifications:write`, `tags:read`, `webhooks:read`, `webhooks:write`, `account:read`, `assets:read`, `me:read` |
| `MONDAY_UPDATE_ITEM` | Update item | write | `board_id`, `column_id`, `item_id`, `value`; *create_labels_if_missing* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_UPDATE_MULTIPLE_USERS` | Update multiple users | write | `user_updates` |  | `users:write` |
| `MONDAY_UPDATE_MUTE_BOARD_SETTINGS` | Update mute board settings | write | `board_id`, `mute_state`; *enabled* |  | `boards:write` |
| `MONDAY_UPDATE_TAG` | Update tag | write | `tag_id`; *color, name* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_UPDATE_TEAM` | Update team | write | `team_id`; *is_guest_team, name, parent_team_id, subscriber_ids* |  | `account:read`, `assets:read`, `boards:read`, `boards:write`, `docs:read`, `docs:write`, `me:read`, `notifications:write`, `tags:read`, `teams:read`, `teams:write`, `updates:read`, `updates:write`, `users:read`, `users:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write` |
| `MONDAY_UPDATE_USERS_ROLE` | Update users role | write | `user_ids`; *new_role, role_id* |  | `users:write` |
| `MONDAY_UPDATE_WORKSPACE` | Update a Monday workspace | write | `workspace_id`; *account_product_id, description, kind, name* |  | `workspaces:read`, `workspaces:write` |
| `MONDAY_UPLOAD_ASSET` | Upload Asset | write | `file`, `target`; *column_id, item_id, update_id* |  | `updates:write`, `boards:write` |
| `MONDAY_USER_CONNECTIONS` | Get user connections | read | *order, page, page_size, with_automations, with_state_validation* |  | `me:read`, `users:read`, `teams:read` |
| `MONDAY_WHO_AM_I` | Who Am I | read |  |  |  |
