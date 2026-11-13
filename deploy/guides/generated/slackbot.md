# Slackbot: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

105 tools: 42 read, 63 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SLACKBOT_ADD_CALL_PARTICIPANTS` | Add call participants | write | `id`, `users` |  | `calls:write` |
| `SLACKBOT_ADD_CHANNEL_BOOKMARK` | Add channel bookmark | write | `channel_id`, `link`, `title`; *emoji, parent_id* |  | `bookmarks:write` |
| `SLACKBOT_ADD_REACTION_TO_AN_ITEM` | Add reaction to message | write | `channel`, `name`, `timestamp` |  | `reactions:write` |
| `SLACKBOT_ADD_REMOTE_FILE` | Add a remote file | write | `external_id`, `external_url`, `title`; *filetype, indexable_file_contents, preview_image* |  |  |
| `SLACKBOT_ARCHIVE_CONVERSATION` | Archive a Slack conversation | write | *channel* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_CLOSE_DM` | Close conversation channel | write | `channel` |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_COMPLETE_EXTERNAL_FILE_UPLOAD` | Complete external file upload | write | `files`; *blocks, channel_id, channels, initial_comment, thread_ts* |  | `files:write` |
| `SLACKBOT_CREATE_A_REMINDER` | Create a reminder | write | `text`, `time`; *team_id, user* |  | `reminders:write` |
| `SLACKBOT_CREATE_CANVAS` | Create Slack Canvas | write | *channel_id, document_content, title* |  | `canvases:write` |
| `SLACKBOT_CREATE_CHANNEL` | Create channel | write | `name`; *is_private, team_id* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_CREATE_LIST` | Create Slack List | write | `name`; *copy_from_list_id, description_blocks, include_copied_list_records, schema, todo_mode* |  | `lists:write` |
| `SLACKBOT_CREATE_LIST_ITEM` | Create Slack List item | write | `list_id`; *duplicated_item_id, initial_fields, parent_item_id* |  | `lists:write` |
| `SLACKBOT_CREATE_USER_GROUP` | Create a Slack user group | write | `name`; *additional_channels, channels, description, enable_section, handle, include_count, team_id* |  | `usergroups:write` |
| `SLACKBOT_CUSTOMIZE_URL_UNFURL` | Customize URL unfurl | write | *channel, metadata, source, ts, unfurl_id, unfurls, user_auth_blocks, user_auth_message and 2 more* |  | `links:write` |
| `SLACKBOT_DELETES_A_MESSAGE_FROM_A_CHAT` | Delete a message from a chat | write | *as_user, channel, ts* |  | `chat:write` |
| `SLACKBOT_DELETE_CANVAS` | Delete Slack Canvas | write | `canvas_id` |  | `canvases:write` |
| `SLACKBOT_DELETE_FILE` | Delete a file by ID | write | `file` |  | `files:write` |
| `SLACKBOT_DELETE_FILE_COMMENT` | Delete file comment | write | `file`, `id` |  | `files:write` |
| `SLACKBOT_DELETE_LIST_ACCESS` | Delete Slack List access | write | `list_id`; *channel_ids, user_ids* |  | `lists:write` |
| `SLACKBOT_DELETE_LIST_ITEM` | Delete Slack List item | write | `id`, `list_id` |  | `lists:write` |
| `SLACKBOT_DELETE_MULTIPLE_LIST_ITEMS` | Delete multiple Slack List items | write | `ids`, `list_id` |  | `lists:write` |
| `SLACKBOT_DELETE_REMINDER` | Delete a Slack reminder | write | `reminder`; *team_id* |  | `reminders:write` |
| `SLACKBOT_DELETE_SCHEDULED_MESSAGE` | Delete scheduled chat message | write | `channel`, `scheduled_message_id`; *as_user* |  | `chat:write` |
| `SLACKBOT_DISABLE_USER_GROUP` | Disable a Slack user group | write | `usergroup`; *include_count, team_id* |  | `usergroups:write` |
| `SLACKBOT_DOWNLOAD_FILE` | Download Slack file | read | `file`; *count, cursor, limit, page* |  | `files:read` |
| `SLACKBOT_EDIT_CANVAS` | Edit Slack Canvas | write | `canvas_id`; *document_content, operation, section_id, title_content* |  | `canvases:write` |
| `SLACKBOT_ENABLE_PUBLIC_SHARING_OF_A_FILE` | Share file public url | write | `file` |  | `files:write` |
| `SLACKBOT_ENABLE_USER_GROUP` | Enable a user group | write | `usergroup`; *include_count, team_id* |  | `usergroups:write` |
| `SLACKBOT_END_CALL` | End a call | write | `id`; *duration* |  | `calls:write` |
| `SLACKBOT_FETCH_CONVERSATION_HISTORY` | Fetch conversation history | read | `channel`; *cursor, include_all_metadata, inclusive, latest, limit, oldest* |  | `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACKBOT_FETCH_ITEM_REACTIONS` | Fetch item reactions | read | *channel, file, file_comment, full, timestamp* |  | `reactions:read` |
| `SLACKBOT_FETCH_MESSAGE_THREAD_FROM_A_CONVERSATION` | Retrieve conversation replies | read | `channel`, `ts`; *cursor, include_all_metadata, inclusive, latest, limit, oldest* |  | `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACKBOT_FETCH_TEAM_INFO` | Fetch team info | read | *domain, team* |  | `team:read` |
| `SLACKBOT_FIND_CHANNELS` | Find channels | read | `query`; *cursor, exact_match, exclude_archived, limit, member_only, team_id, types* |  | `channels:read`, `groups:read`, `mpim:read`, `im:read` |
| `SLACKBOT_FIND_USERS` | Find users | read | *cursor, email, exact_match, include_bots, include_deleted, include_locale, include_restricted, limit and 2 more* |  | `users:read`, `users:read.email` |
| `SLACKBOT_FIND_USER_BY_EMAIL_ADDRESS` | Lookup users by email | read | `email` |  | `users:read`, `users:read.email` |
| `SLACKBOT_GET_BOT_USER` | Fetch bot user information | read | *bot, team_id* |  | `users:read` |
| `SLACKBOT_GET_CALL_INFO` | Retrieve call information | read | `id` |  | `calls:read` |
| `SLACKBOT_GET_LIST_ITEM` | Get Slack List item | read | `id`, `list_id`; *include_is_subscribed* |  | `lists:read` |
| `SLACKBOT_GET_REMINDER` | Get reminder information | read | `reminder`; *team_id* |  | `reminders:read` |
| `SLACKBOT_GET_REMOTE_FILE` | Get remote file | read | *external_id, file* |  | `remote_files:read` |
| `SLACKBOT_GET_TEAM_PROFILE` | Retrieve team profile details | read | *team_id, visibility* |  | `users.profile:read` |
| `SLACKBOT_GET_UNREAD_MESSAGES_FROM_USER` | Get unread messages from user | read | `user_id`; *limit* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write`, `channels:read`, `groups:read`, `im:read`, `mpim:read`, `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACKBOT_GET_USER_DND_STATUS` | Get team DND status | read | `users`; *team_id* |  | `dnd:read` |
| `SLACKBOT_GET_USER_PRESENCE` | Retrieve user presence | read | *user* |  | `users:read` |
| `SLACKBOT_INVITE_USERS_TO_A_CHANNEL` | Invite users to a Slack channel | write | *channel, force, users* |  | `channels:write`, `channels:write.invites`, `groups:write`, `groups:write.invites`, `im:write`, `mpim:write` |
| `SLACKBOT_JOIN_AN_EXISTING_CONVERSATION` | Join conversation by channel id | write | `channel` |  | `channels:write` |
| `SLACKBOT_LEAVE_CONVERSATION` | Leave conversation channel | write | `channel` |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_LIST_ALL_CHANNELS` | List all channels | read | *cursor, exclude_archived, limit, team_id, types* |  | `channels:read`, `groups:read`, `mpim:read`, `im:read` |
| `SLACKBOT_LIST_ALL_USERS` | List all users | read | *cursor, include_locale, limit, team_id* |  | `users:read`, `users:read.email` |
| `SLACKBOT_LIST_AVAILABLE_WORKSPACES` | List available workspaces | read | `user`; *cursor* |  | `users:read` |
| `SLACKBOT_LIST_CHANNEL_BOOKMARKS` | List channel bookmarks | read | `channel_id` |  | `bookmarks:read` |
| `SLACKBOT_LIST_CONVERSATIONS` | List conversations | read | *cursor, exclude_archived, limit, team_id, types, user* |  | `channels:read`, `groups:read`, `im:read`, `mpim:read` |
| `SLACKBOT_LIST_CUSTOM_EMOJIS` | List team custom emojis | write | *include_categories* |  | `emoji:read` |
| `SLACKBOT_LIST_FILES_WITH_FILTERS_IN_SLACK` | List Slack files | read | *channel, count, page, show_files_hidden_by_limit, team_id, ts_from, ts_to, types and 1 more* |  | `files:read` |
| `SLACKBOT_LIST_LIST_ITEMS` | List Slack List items | read | `list_id`; *archived, cursor, limit* |  | `lists:read` |
| `SLACKBOT_LIST_PINNED_ITEMS` | List pinned items in a channel | read | `channel` |  | `pins:read` |
| `SLACKBOT_LIST_REMINDERS` | List reminders | read | *team_id* |  | `reminders:read` |
| `SLACKBOT_LIST_REMOTE_FILES` | List remote files | read | *channel, cursor, limit, ts_from, ts_to* |  | `remote_files:read` |
| `SLACKBOT_LIST_UNREAD_CHANNEL_MESSAGES` | List unread channel messages | read | `channel`; *cursor, include_all_metadata, limit* |  | `channels:read`, `groups:read`, `im:read`, `mpim:read`, `channels:history`, `groups:history`, `im:history`, `mpim:history` |
| `SLACKBOT_LIST_USER_GROUPS` | List user groups | read | *include_count, include_disabled, include_users, team_id* |  | `usergroups:read` |
| `SLACKBOT_LIST_USER_GROUP_MEMBERS` | List all users in a user group | read | `usergroup`; *include_disabled, team_id* |  | `usergroups:read` |
| `SLACKBOT_LIST_USER_REACTIONS` | List user reactions | read | *count, cursor, full, limit, page, team_id, user* |  | `reactions:read` |
| `SLACKBOT_LOOKUP_CANVAS_SECTIONS` | Lookup Canvas Sections | read | `canvas_id`, `criteria` |  | `canvases:read` |
| `SLACKBOT_OPEN_DM` | Open DM | write | *channel, prevent_creation, return_im, users* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_PIN_ITEM` | Pin an item to a channel | write | `channel`, `timestamp` |  | `pins:write` |
| `SLACKBOT_REMOVE_CALL_PARTICIPANTS` | Remove call participants | write | `id`, `users` |  | `calls:write` |
| `SLACKBOT_REMOVE_REACTION_FROM_ITEM` | Remove reaction from item | write | `name`; *channel, file, file_comment, timestamp* |  | `reactions:write` |
| `SLACKBOT_REMOVE_REMOTE_FILE` | Remove remote file | write | *external_id, file, token* |  |  |
| `SLACKBOT_REMOVE_USER_FROM_CONVERSATION` | Remove user from conversation | write | *channel, user* |  | `channels:write`, `groups:write` |
| `SLACKBOT_RENAME_CONVERSATION` | Rename a conversation | write | *channel, name* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_RESERVE_FILE_UPLOAD` | Reserve Slack file upload | write | `filename`, `length`; *alt_txt, snippet_type* |  | `files:write` |
| `SLACKBOT_RETRIEVE_CONVERSATION_INFORMATION` | Retrieve conversation information | read | *channel, include_locale, include_num_members* |  | `channels:read`, `groups:read`, `im:read`, `mpim:read` |
| `SLACKBOT_RETRIEVE_CONVERSATION_MEMBERS_LIST` | Get conversation members | read | *channel, cursor, limit* |  | `channels:read`, `groups:read`, `mpim:read`, `im:read` |
| `SLACKBOT_RETRIEVE_CURRENT_USER_DND_STATUS` | Retrieve user DND status | read | *team_id, user* |  | `dnd:read` |
| `SLACKBOT_RETRIEVE_DETAILED_INFORMATION_ABOUT_A_FILE` | Retrieve detailed file information | read | `file`; *count, cursor, limit, page* |  | `files:read` |
| `SLACKBOT_RETRIEVE_DETAILED_USER_INFORMATION` | Retrieve detailed user information | read | *include_locale, user* |  | `users:read`, `users:read.email` |
| `SLACKBOT_RETRIEVE_USER_PROFILE_INFORMATION` | Retrieve user profile information | read | *include_labels, user* |  | `users.profile:read` |
| `SLACKBOT_REVOKE_FILE_PUBLIC_SHARING` | Revoke a file's public url | write | `file` |  | `files:write` |
| `SLACKBOT_SCHEDULE_MESSAGE` | Schedule message | write | *attachments, blocks, channel, link_names, markdown_text, parse, post_at, reply_broadcast and 4 more* |  | `chat:write` |
| `SLACKBOT_SEARCH_ALL` | Search all content | read | `query`; *count, highlight, page, sort, sort_dir, team_id* |  | `search:read`, `search:read.public` |
| `SLACKBOT_SEARCH_FILES` | Search Slack files | read | `query`; *count, highlight, page, sort, sort_dir, team_id* |  | `search:read`, `search:read.public` |
| `SLACKBOT_SEARCH_MESSAGES` | Search messages | read | `query`; *auto_paginate, count, cursor, highlight, page, sort, sort_dir, team_id* |  | `search:read`, `search:read.public` |
| `SLACKBOT_SEND_EPHEMERAL_MESSAGE` | Send ephemeral message | write | `channel`, `user`; *as_user, attachments, blocks, icon_emoji, icon_url, link_names, markdown_text, parse and 3 more* |  | `chat:write` |
| `SLACKBOT_SEND_MESSAGE` | Send message | write | `channel`; *blocks, fallback_text, markdown_text, reply_broadcast, thread_ts, unfurl_links, unfurl_media* |  | `chat:write` |
| `SLACKBOT_SEND_ME_MESSAGE` | Share a me message in a channel | write | *channel, text* |  | `chat:write` |
| `SLACKBOT_SET_CANVAS_ACCESS` | Set canvas access | write | `access_level`, `canvas_id`; *channel_ids, user_ids* |  | `canvases:write` |
| `SLACKBOT_SET_CONVERSATION_PURPOSE` | Set a conversation's purpose | write | *channel, purpose* |  | `channels:write`, `channels:write.topic`, `groups:write`, `groups:write.topic`, `im:write`, `im:write.topic`, `mpim:write`, `mpim:write.topic` |
| `SLACKBOT_SET_LIST_ACCESS` | Set Slack List access | write | `access_level`, `list_id`; *channel_ids, user_ids* |  | `lists:write` |
| `SLACKBOT_SET_READ_CURSOR_IN_A_CONVERSATION` | Set conversation read cursor | write | *channel, ts* |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_SET_THE_TOPIC_OF_A_CONVERSATION` | Set conversation topic | write | *channel, topic* |  | `channels:write`, `channels:write.topic`, `groups:write`, `groups:write.topic`, `im:write`, `im:write.topic`, `mpim:write`, `mpim:write.topic` |
| `SLACKBOT_SET_USER_PRESENCE` | Set user presence | write | `presence` |  | `users:write` |
| `SLACKBOT_SHARE_REMOTE_FILE` | Share a remote file in channels | write | `channels`; *external_id, file* |  | `remote_files:share` |
| `SLACKBOT_START_CALL` | Start call | write | `external_unique_id`, `join_url`; *created_by, date_start, desktop_app_join_url, external_display_id, title, users* |  | `calls:write` |
| `SLACKBOT_TEST_AUTH` | Test authentication | read |  |  |  |
| `SLACKBOT_UNARCHIVE_CHANNEL` | Unarchive channel | write | `channel` |  | `channels:write`, `groups:write`, `im:write`, `mpim:write` |
| `SLACKBOT_UNPIN_ITEM` | Unpin message from channel | write | `channel`, `timestamp` |  | `pins:write` |
| `SLACKBOT_UPDATES_A_MESSAGE` | Update a Slack message | write | `channel`, `ts`; *as_user, blocks, fallback_text, file_ids, markdown_text, metadata, reply_broadcast* |  | `chat:write` |
| `SLACKBOT_UPDATE_CALL_INFO` | Update call information | write | `id`; *desktop_app_join_url, join_url, title* |  | `calls:write` |
| `SLACKBOT_UPDATE_LIST` | Update Slack List | write | `id`; *description_blocks, name, todo_mode* |  | `lists:write` |
| `SLACKBOT_UPDATE_LIST_ITEM` | Update Slack List item | write | `cells`, `list_id` |  | `lists:write` |
| `SLACKBOT_UPDATE_REMOTE_FILE` | Update an existing remote file | write | *external_id, external_url, file, filetype, indexable_file_contents, preview_image, title, token* |  |  |
| `SLACKBOT_UPDATE_USER_GROUP` | Update Slack user group | write | `usergroup`; *additional_channels, channels, description, enable_section, handle, include_count, name, team_id* |  | `usergroups:write` |
| `SLACKBOT_UPDATE_USER_GROUP_MEMBERS` | Update user group members | write | `usergroup`, `users`; *include_count, team_id* |  | `usergroups:write` |
| `SLACKBOT_UPLOAD_OR_CREATE_A_FILE_IN_SLACK` | Upload or create a file in Slack | write | *blocks, channels, content, file, filename, filetype, initial_comment, thread_ts and 2 more* |  | `files:write` |
