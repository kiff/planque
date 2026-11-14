# Slack: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

165 tools: 71 read, 94 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SLACK_ADD_CALL_PARTICIPANTS` | Add call participants | write | `id`, `users` |  | `calls:write` |
| `SLACK_ADD_CHANNEL_BOOKMARK` | Add channel bookmark | write | `channel_id`, `link`, `title`; *emoji, parent_id* |  | `bookmarks:write` |
| `SLACK_ADD_EMOJI` | Add emoji | write | `name`, `url` |  | `admin.teams:write` |
| `SLACK_ADD_EMOJI_ALIAS` | Add an emoji alias | write | `alias_for`, `name` |  | `admin.teams:write` |
| `SLACK_ADD_ENTERPRISE_USER_TO_WORKSPACE` | Add Enterprise user to workspace | write | `team_id`, `user_id`; *channel_ids, is_restricted, is_ultra_restricted* |  | `admin.users:write` |
| `SLACK_ADD_REACTION_TO_AN_ITEM` | Add reaction to message | write | `channel`, `name`, `timestamp` |  | `reactions:write` |
| `SLACK_ADD_REMOTE_FILE` | Add a remote file | write | `external_id`, `external_url`, `title`; *filetype, indexable_file_contents, preview_image* |  |  |
| `SLACK_ADD_STAR` | Add a star to an item | write | *channel, file, file_comment, timestamp* |  | `stars:write` |
| `SLACK_ADMIN_CONVERSATIONS_SEARCH` | Search for channels in Enterprise organization | read | *connected_team_ids, cursor, limit, query, search_channel_types, sort, sort_dir, team_ids and 1 more* |  | `admin.conversations:read` |
| `SLACK_API_TEST` | Test Slack API connection | read | *error, foo* |  |  |
| `SLACK_ARCHIVE_CONVERSATION` | Archive a Slack conversation | write | *channel* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_ASSISTANT_SEARCH_CONTEXT` | Real-time search | read | `query`; *action_token, after, before, channel_types, content_types, context_channel_id, cursor, disable_semantic_search and 11 more* |  | `search:read.public`, `search:read.files`, `search:read.users`, `search:read.private`, `search:read.im`, `search:read.mpim` |
| `SLACK_ASSISTANT_SEARCH_INFO` | Check search capabilities | read |  |  | `search:read.public`, `search:read` |
| `SLACK_CLOSE_DM` | Close conversation channel | write | `channel` |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_COMPLETE_EXTERNAL_FILE_UPLOAD` | Complete external file upload | write | `files`; *blocks, channel_id, channels, initial_comment, thread_ts* |  | `files:write` |
| `SLACK_CONVERT_CHANNEL_TO_PRIVATE` | Convert public channel to private | write | `channel_id`; *name* |  | `admin.conversations:write` |
| `SLACK_CREATE_A_REMINDER` | Create a reminder | write | `text`, `time`; *team_id, user* |  | `reminders:write` |
| `SLACK_CREATE_CANVAS` | Create Slack Canvas | write | *channel_id, document_content, title* |  | `canvases:write` |
| `SLACK_CREATE_CHANNEL` | Create channel | write | `name`; *is_private, team_id* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_CREATE_CHANNEL_BASED_CONVERSATION` | Create a channel-based conversation | write | `is_private`, `name`; *description, org_wide, team_id* |  | `admin.conversations:write` |
| `SLACK_CREATE_ENTERPRISE_TEAM` | Create Enterprise team | write | `team_domain`, `team_name`; *team_description, team_discoverability* |  | `admin.teams:write` |
| `SLACK_CREATE_SLACK_LIST` | Create Slack List | write | `name`; *copy_from_list_id, description_blocks, include_copied_list_records, schema, todo_mode* |  | `lists:write` |
| `SLACK_CREATE_SLACK_LIST_ITEM` | Create Slack List item | write | `list_id`; *duplicated_item_id, initial_fields, parent_item_id* |  | `lists:write` |
| `SLACK_CREATE_USER_GROUP` | Create a Slack user group | write | `name`; *additional_channels, channels, description, enable_section, handle, include_count, team_id* |  | `usergroups:write` |
| `SLACK_CUSTOMIZE_URL_UNFURL` | Customize URL unfurl | write | *channel, metadata, source, ts, unfurl_id, unfurls, user_auth_blocks, user_auth_message and 2 more* |  | `links:write` |
| `SLACK_DELETES_A_MESSAGE_FROM_A_CHAT` | Delete a message from a chat | write | *as_user, channel, ts* |  | `chat:write` |
| `SLACK_DELETE_CANVAS` | Delete Slack Canvas | write | `canvas_id` |  | `canvases:write` |
| `SLACK_DELETE_CHANNEL` | Delete a public or private channel | write | `channel_id` |  | `admin.conversations:write` |
| `SLACK_DELETE_FILE` | Delete a file by ID | write | `file` |  | `files:write` |
| `SLACK_DELETE_FILE_COMMENT` | Delete file comment | write | `file`, `id` |  | `files:write` |
| `SLACK_DELETE_MULTIPLE_SLACK_LIST_ITEMS` | Delete multiple Slack List items | write | `ids`, `list_id` |  | `lists:write` |
| `SLACK_DELETE_REMINDER` | Delete a Slack reminder | write | `reminder`; *team_id* |  | `reminders:write` |
| `SLACK_DELETE_SCHEDULED_MESSAGE` | Delete scheduled chat message | write | `channel`, `scheduled_message_id`; *as_user* |  | `chat:write` |
| `SLACK_DELETE_SLACK_LIST_ACCESS` | Delete Slack List access | write | `list_id`; *channel_ids, user_ids* |  | `lists:write` |
| `SLACK_DELETE_SLACK_LIST_ITEM` | Delete Slack List item | write | `id`, `list_id` |  | `lists:write` |
| `SLACK_DELETE_USER_PROFILE_PHOTO` | Delete user profile photo | write |  |  | `users.profile:write` |
| `SLACK_DISABLE_USER_GROUP` | Disable a Slack user group | write | `usergroup`; *include_count, team_id* |  | `usergroups:write` |
| `SLACK_DOWNLOAD_SLACK_FILE` | Download Slack file | read | `file`; *count, cursor, limit, page* |  | `files:read` |
| `SLACK_EDIT_CANVAS` | Edit Slack Canvas | write | `canvas_id`; *document_content, operation, section_id, title_content* |  | `canvases:write` |
| `SLACK_ENABLE_PUBLIC_SHARING_OF_A_FILE` | Share file public url | write | `file` |  | `files:write` |
| `SLACK_ENABLE_USER_GROUP` | Enable a user group | write | `usergroup`; *include_count, team_id* |  | `usergroups:write` |
| `SLACK_END_CALL` | End a call | write | `id`; *duration* |  | `calls:write` |
| `SLACK_END_DND` | End DND session | write |  |  | `dnd:write` |
| `SLACK_END_SNOOZE` | End snooze | write |  |  | `dnd:write` |
| `SLACK_FETCH_CONVERSATION_HISTORY` | Fetch conversation history | read | `channel`; *cursor, include_all_metadata, inclusive, latest, limit, oldest* |  | `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACK_FETCH_ITEM_REACTIONS` | Fetch item reactions | read | *channel, file, file_comment, full, timestamp* |  | `reactions:read` |
| `SLACK_FETCH_MESSAGE_THREAD_FROM_A_CONVERSATION` | Retrieve conversation replies | read | `channel`, `ts`; *cursor, include_all_metadata, inclusive, latest, limit, oldest* |  | `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACK_FETCH_TEAM_INFO` | Fetch team info | read | *domain, team* |  | `team:read` |
| `SLACK_FIND_CHANNELS` | Find channels | read | `query`; *cursor, exact_match, exclude_archived, limit, member_only, team_id, types* |  | `channels:read`, `groups:read`, `mpim:read`, `im:read` |
| `SLACK_FIND_USERS` | Find users | read | *cursor, email, exact_match, include_bots, include_deleted, include_locale, include_restricted, limit and 2 more* |  | `users:read`, `users:read.email` |
| `SLACK_FIND_USER_BY_EMAIL_ADDRESS` | Lookup users by email | read | `email` |  | `users:read`, `users:read.email` |
| `SLACK_GET_AUDIT_ACTION_TYPES` | Get Audit Action Types | read |  |  |  |
| `SLACK_GET_AUDIT_SCHEMAS` | Get Audit Schemas | read |  |  |  |
| `SLACK_GET_BOT_USER` | Fetch bot user information | read | *bot, team_id* |  | `users:read` |
| `SLACK_GET_CALL_INFO` | Retrieve call information | read | `id` |  | `calls:read` |
| `SLACK_GET_CHANNEL_CONVERSATION_PREFERENCES` | Get channel conversation preferences | read | `channel_id` |  | `admin.conversations:read` |
| `SLACK_GET_REMINDER` | Get reminder information | read | `reminder`; *team_id* |  | `reminders:read` |
| `SLACK_GET_REMOTE_FILE` | Get remote file | read | *external_id, file* |  | `remote_files:read` |
| `SLACK_GET_SLACK_LIST_ITEM` | Get Slack List item | read | `id`, `list_id`; *include_is_subscribed* |  | `lists:read` |
| `SLACK_GET_TEAM_PROFILE` | Retrieve team profile details | read | *team_id, visibility* |  | `users.profile:read` |
| `SLACK_GET_UNREAD_MESSAGES_FROM_USER` | Get unread messages from user | read | `user_id`; *limit* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write`, `channels:read`, `groups:read`, `im:read`, `mpim:read`, `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACK_GET_USER_DND_STATUS` | Get team DND status | read | `users`; *team_id* |  | `dnd:read` |
| `SLACK_GET_USER_PRESENCE` | Retrieve user presence | read | *user* |  | `users:read` |
| `SLACK_GET_WORKSPACE_CONNECTIONS_FOR_CHANNEL` | Get workspace connections for channel | read | `channel_id`; *cursor, limit* |  | `admin.conversations:read` |
| `SLACK_GET_WORKSPACE_SETTINGS` | Fetch workspace settings information | read | `team_id` |  | `admin.teams:read` |
| `SLACK_INVITE_USERS_TO_A_SLACK_CHANNEL` | Invite users to a Slack channel | write | *channel, force, users* |  | `channels:write`, `channels:write.invites`, `groups:write`, `groups:write.invites`, `im:write`, `mpim:write` |
| `SLACK_INVITE_USER_TO_CHANNEL` | Invite users to channel | write | `channel_id`, `user_ids` |  | `admin.conversations:write` |
| `SLACK_INVITE_USER_TO_WORKSPACE` | Invite user to workspace | write | `channel_ids`, `email`, `team_id`; *custom_message, email_password_policy_enabled, guest_expiration_ts, is_restricted, is_ultra_restricted, real_name, resend* |  | `admin.users:write` |
| `SLACK_JOIN_AN_EXISTING_CONVERSATION` | Join conversation by channel id | write | `channel` |  | `channels:write` |
| `SLACK_LEAVE_CONVERSATION` | Leave conversation channel | write | `channel` |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_LIST_ADMIN_APPS_APPROVED` | List approved apps | read | *certified, cursor, enterprise_id, limit, team_id* |  | `admin.apps:read` |
| `SLACK_LIST_ADMIN_APPS_REQUESTS` | List app requests | read | *certified, cursor, enterprise_id, limit, team_id* |  | `admin.apps:read` |
| `SLACK_LIST_ADMIN_EMOJI` | List admin emoji | read | *cursor, limit* |  | `admin.teams:read` |
| `SLACK_LIST_ALL_CHANNELS` | List all channels | read | *cursor, exclude_archived, limit, team_id, types* |  | `channels:read`, `groups:read`, `mpim:read`, `im:read` |
| `SLACK_LIST_ALL_USERS` | List all users | read | *cursor, include_locale, limit, team_id* |  | `users:read`, `users:read.email` |
| `SLACK_LIST_APPROVED_WORKSPACE_INVITE_REQUESTS` | List approved workspace invite requests | read | *cursor, limit, team_id* |  | `admin.invites:read` |
| `SLACK_LIST_AUTH_TEAMS` | List authorized teams | read | *cursor, include_icon, limit* |  |  |
| `SLACK_LIST_AVAILABLE_WORKSPACES` | List available workspaces | read | `user`; *cursor* |  | `users:read` |
| `SLACK_LIST_CHANNEL_BOOKMARKS` | List channel bookmarks | read | `channel_id` |  | `bookmarks:read` |
| `SLACK_LIST_CONVERSATIONS` | List conversations | read | *cursor, exclude_archived, limit, team_id, types, user* |  | `channels:read`, `groups:read`, `im:read`, `mpim:read` |
| `SLACK_LIST_CUSTOM_EMOJIS` | List team custom emojis | write | *include_categories* |  | `emoji:read` |
| `SLACK_LIST_DENIED_WORKSPACE_INVITE_REQUESTS` | List denied workspace invite requests | read | *cursor, limit, team_id* |  | `admin.invites:read` |
| `SLACK_LIST_ENTERPRISE_TEAMS` | List Enterprise teams | read | *cursor, limit* |  | `admin.teams:read` |
| `SLACK_LIST_FILES_WITH_FILTERS_IN_SLACK` | List Slack files | read | *channel, count, page, show_files_hidden_by_limit, team_id, ts_from, ts_to, types and 1 more* |  | `files:read` |
| `SLACK_LIST_IDP_GROUPS_LINKED_TO_CHANNEL` | List IDP groups linked to channel | read | `channel_id`; *team_id* |  | `admin.conversations:read` |
| `SLACK_LIST_PENDING_WORKSPACE_INVITE_REQUESTS` | List pending workspace invite requests | read | *cursor, limit, team_id* |  | `admin.invites:read` |
| `SLACK_LIST_PINNED_ITEMS` | List pinned items in a channel | read | `channel` |  | `pins:read` |
| `SLACK_LIST_REMINDERS` | List reminders | read | *team_id* |  | `reminders:read` |
| `SLACK_LIST_REMOTE_FILES` | List remote files | read | *channel, cursor, limit, ts_from, ts_to* |  | `remote_files:read` |
| `SLACK_LIST_RESTRICTED_APPS` | List Restricted Apps | read | *certified, cursor, enterprise_id, limit, team_id* |  | `admin.apps:read` |
| `SLACK_LIST_SCHEDULED_MESSAGES` | List scheduled messages | read | *channel, cursor, latest, limit, oldest, team_id* |  |  |
| `SLACK_LIST_SLACK_LIST_ITEMS` | List Slack List items | read | `list_id`; *archived, cursor, limit* |  | `lists:read` |
| `SLACK_LIST_STARRED_ITEMS` | List starred items | read | *count, cursor, limit, page, team_id* |  | `stars:read` |
| `SLACK_LIST_UNREAD_CHANNEL_MESSAGES` | List unread channel messages | read | `channel`; *cursor, include_all_metadata, limit* |  | `channels:read`, `groups:read`, `im:read`, `mpim:read`, `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACK_LIST_USER_GROUPS` | List user groups | read | *include_count, include_disabled, include_users, team_id* |  | `usergroups:read` |
| `SLACK_LIST_USER_GROUP_MEMBERS` | List all users in a user group | read | `usergroup`; *include_disabled, team_id* |  | `usergroups:read` |
| `SLACK_LIST_USER_REACTIONS` | List user reactions | read | *count, cursor, full, limit, page, team_id, user* |  | `reactions:read` |
| `SLACK_LIST_WORKSPACE_ADMINS` | List workspace admins | read | `team_id`; *cursor, limit* |  | `admin.teams:read` |
| `SLACK_LIST_WORKSPACE_OWNERS` | List workspace owners | read | `team_id`; *cursor, limit* |  | `admin.teams:read` |
| `SLACK_LIST_WORKSPACE_USERS` | List admin users | read | *cursor, include_deactivated_user_workspaces, is_active, limit, only_guests, team_id* |  | `admin.users:read` |
| `SLACK_LOOKUP_CANVAS_SECTIONS` | Lookup Canvas Sections | read | `canvas_id`, `criteria` |  | `canvases:read` |
| `SLACK_OPEN_DM` | Open DM | write | *channel, prevent_creation, return_im, users* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_PIN_ITEM` | Pin an item to a channel | write | `channel`, `timestamp` |  | `pins:write` |
| `SLACK_READ_AUDIT_LOGS` | Read Audit Logs | read | *action, actor, cursor, entity, latest, limit, oldest* |  | `auditlogs:read` |
| `SLACK_REMOVE_CALL_PARTICIPANTS` | Remove call participants | write | `id`, `users` |  | `calls:write` |
| `SLACK_REMOVE_EMOJI` | Remove emoji | write | `name` |  | `admin.teams:write` |
| `SLACK_REMOVE_REACTION_FROM_ITEM` | Remove reaction from item | write | `name`; *channel, file, file_comment, timestamp* |  | `reactions:write` |
| `SLACK_REMOVE_REMOTE_FILE` | Remove remote file | write | *external_id, file, token* |  |  |
| `SLACK_REMOVE_STAR` | Remove a star from an item | write | *channel, file, file_comment, timestamp* |  | `stars:write` |
| `SLACK_REMOVE_USER_FROM_CONVERSATION` | Remove user from conversation | write | *channel, user* |  | `channels:write`, `groups:write` |
| `SLACK_REMOVE_USER_FROM_WORKSPACE` | Remove user from workspace | write | `team_id`, `user_id` |  | `admin.users:write` |
| `SLACK_RENAME_CONVERSATION` | Rename a conversation | write | *channel, name* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_RENAME_EMOJI` | Rename an emoji | write | `name`, `new_name` |  | `admin.teams:write` |
| `SLACK_RESERVE_FILE_UPLOAD` | Reserve Slack file upload | write | `filename`, `length`; *alt_txt, snippet_type* |  | `files:write` |
| `SLACK_RESET_USER_SESSIONS` | Reset user sessions | write | `user_id`; *mobile_only, web_only* |  | `admin.users:write` |
| `SLACK_RESTRICT_APP_INSTALLATION` | Restrict app installation | write | *app_id, enterprise_id, request_id, team_id* |  | `admin.apps:write` |
| `SLACK_RETRIEVE_A_USER_S_IDENTITY_DETAILS` | Retrieve a user's identity details | read |  |  |  |
| `SLACK_RETRIEVE_CONVERSATION_INFORMATION` | Retrieve conversation information | read | *channel, include_locale, include_num_members* |  | `channels:read`, `groups:read`, `im:read`, `mpim:read` |
| `SLACK_RETRIEVE_CONVERSATION_MEMBERS_LIST` | Get conversation members | read | *channel, cursor, limit* |  | `channels:read`, `groups:read`, `mpim:read`, `im:read` |
| `SLACK_RETRIEVE_CURRENT_USER_DND_STATUS` | Retrieve user DND status | read | *team_id, user* |  | `dnd:read` |
| `SLACK_RETRIEVE_DETAILED_INFORMATION_ABOUT_A_FILE` | Retrieve detailed file information | read | `file`; *count, cursor, limit, page* |  | `files:read` |
| `SLACK_RETRIEVE_DETAILED_USER_INFORMATION` | Retrieve detailed user information | read | *include_locale, user* |  | `users:read`, `users:read.email` |
| `SLACK_RETRIEVE_MESSAGE_PERMALINK_URL` | Retrieve message permalink | read | `channel`, `message_ts` |  |  |
| `SLACK_RETRIEVE_USER_PROFILE_INFORMATION` | Retrieve user profile information | read | *include_labels, user* |  | `users.profile:read` |
| `SLACK_REVOKE_FILE_PUBLIC_SHARING` | Revoke a file's public url | write | `file` |  | `files:write` |
| `SLACK_RTM_CONNECT` | Start RTM session | write | *batch_presence_aware, presence_sub* |  |  |
| `SLACK_SCHEDULE_MESSAGE` | Schedule message | write | *attachments, blocks, channel, link_names, markdown_text, parse, post_at, reply_broadcast and 4 more* |  | `chat:write` |
| `SLACK_SCIM_GET_CONFIG` | Get SCIM service provider configuration | read |  |  | `admin` |
| `SLACK_SEARCH_ALL` | Search all content | read | `query`; *count, highlight, page, sort, sort_dir, team_id* |  | `search:read` |
| `SLACK_SEARCH_FILES` | Search Slack files | read | `query`; *count, highlight, page, sort, sort_dir, team_id* |  | `search:read` |
| `SLACK_SEARCH_MESSAGES` | Search messages | read | `query`; *auto_paginate, count, cursor, highlight, page, sort, sort_dir, team_id* |  | `search:read` |
| `SLACK_SEND_EPHEMERAL_MESSAGE` | Send ephemeral message | write | `channel`, `user`; *as_user, attachments, blocks, icon_emoji, icon_url, link_names, markdown_text, parse and 3 more* |  | `chat:write` |
| `SLACK_SEND_MESSAGE` | Send message | write | `channel`; *blocks, fallback_text, markdown_text, reply_broadcast, thread_ts, unfurl_links, unfurl_media* |  | `chat:write` |
| `SLACK_SEND_ME_MESSAGE` | Share a me message in a channel | write | *channel, text* |  | `chat:write` |
| `SLACK_SET_ADMIN_USER` | Set admin user | write | `team_id`, `user_id` |  | `admin.users:write` |
| `SLACK_SET_CANVAS_ACCESS` | Set canvas access | write | `access_level`, `canvas_id`; *channel_ids, user_ids* |  | `canvases:write` |
| `SLACK_SET_CONVERSATION_PREFS` | Set conversation preferences | write | `channel_id`, `prefs` |  | `admin.conversations:write` |
| `SLACK_SET_CONVERSATION_PURPOSE` | Set a conversation's purpose | write | *channel, purpose* |  | `channels:write`, `channels:write.topic`, `groups:write`, `groups:write.topic`, `im:write`, `im:write.topic`, `mpim:write`, `mpim:write.topic` |
| `SLACK_SET_DEFAULT_CHANNELS` | Set default channels | write | `channel_ids`, `team_id` |  | `admin.teams:write` |
| `SLACK_SET_DND_DURATION` | Set DND duration | write | `num_minutes` |  | `dnd:write` |
| `SLACK_SET_PROFILE_PHOTO` | Set profile photo | write | `image`; *crop_w, crop_x, crop_y* |  | `users.profile:write` |
| `SLACK_SET_READ_CURSOR_IN_A_CONVERSATION` | Set conversation read cursor | write | *channel, ts* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_SET_SLACK_LIST_ACCESS` | Set Slack List access | write | `access_level`, `list_id`; *channel_ids, user_ids* |  | `lists:write` |
| `SLACK_SET_THE_TOPIC_OF_A_CONVERSATION` | Set conversation topic | write | *channel, topic* |  | `channels:write`, `channels:write.topic`, `groups:write`, `groups:write.topic`, `im:write`, `im:write.topic`, `mpim:write`, `mpim:write.topic` |
| `SLACK_SET_USER_PRESENCE` | Set user presence | write | `presence` |  | `users:write` |
| `SLACK_SET_USER_PROFILE` | Set Slack user profile information | write | *name, profile, user, value* |  | `users.profile:write` |
| `SLACK_SET_WORKSPACES_FOR_CHANNEL` | Set workspaces for channel | write | `channel_id`; *org_channel, target_team_ids, team_id* |  | `admin.conversations:write` |
| `SLACK_SET_WORKSPACE_DESCRIPTION` | Set workspace description | write | `description`, `team_id` |  | `admin.teams:write` |
| `SLACK_SET_WORKSPACE_ICON` | Set workspace icon | write | `image_url`, `team_id` |  | `admin.teams:write` |
| `SLACK_SET_WORKSPACE_NAME` | Set workspace name | write | `name`, `team_id` |  | `admin.teams:write` |
| `SLACK_SET_WORKSPACE_OWNER` | Set workspace owner | write | `team_id`, `user_id` |  | `admin.users:write` |
| `SLACK_SHARE_REMOTE_FILE` | Share a remote file in channels | write | `channels`; *external_id, file* |  | `remote_files:share` |
| `SLACK_START_CALL` | Start call | write | `external_unique_id`, `join_url`; *created_by, date_start, desktop_app_join_url, external_display_id, title, users* |  | `calls:write` |
| `SLACK_TEST_AUTH` | Test authentication | read |  |  |  |
| `SLACK_UNARCHIVE_CHANNEL` | Unarchive channel | write | `channel` |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACK_UNPIN_ITEM` | Unpin message from channel | write | `channel`, `timestamp` |  | `pins:write` |
| `SLACK_UPDATES_A_SLACK_MESSAGE` | Update a Slack message | write | `channel`, `ts`; *as_user, blocks, fallback_text, file_ids, markdown_text, metadata, reply_broadcast* |  | `chat:write` |
| `SLACK_UPDATE_CALL_INFO` | Update call information | write | `id`; *desktop_app_join_url, join_url, title* |  | `calls:write` |
| `SLACK_UPDATE_REMOTE_FILE` | Update an existing remote file | write | *external_id, external_url, file, filetype, indexable_file_contents, preview_image, title, token* |  |  |
| `SLACK_UPDATE_SLACK_LIST` | Update Slack List | write | `id`; *description_blocks, name, todo_mode* |  | `lists:write` |
| `SLACK_UPDATE_SLACK_LIST_ITEM` | Update Slack List item | write | `cells`, `list_id` |  | `lists:write` |
| `SLACK_UPDATE_USER_GROUP` | Update Slack user group | write | `usergroup`; *additional_channels, channels, description, enable_section, handle, include_count, name, team_id* |  | `usergroups:write` |
| `SLACK_UPDATE_USER_GROUP_MEMBERS` | Update user group members | write | `usergroup`, `users`; *include_count, team_id* |  | `usergroups:write` |
| `SLACK_UPLOAD_OR_CREATE_A_FILE_IN_SLACK` | Upload or create a file in Slack | write | *blocks, channels, content, file, filename, filetype, initial_comment, thread_ts and 2 more* |  | `files:write` |
| `SLACK_WHO_AM_I` | Who Am I | read |  |  |  |
